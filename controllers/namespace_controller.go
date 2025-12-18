package controllers

import (
    "context"
    "io"
    "os"
    "strings"
    "time"

    "k8s.io/apimachinery/pkg/api/errors"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/apimachinery/pkg/util/yaml"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"
)

// NamespaceReconciler watches namespaces and applies instrumentation manifests
type NamespaceReconciler struct {
    client.Client
    Scheme    *runtime.Scheme
    InstrPath string
    // map for quick lookup of ignored namespaces
    IgnoreNamespaces map[string]bool
}

// SetupWithManager registers the controller with the manager
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&corev1.Namespace{}).
        Complete(r)
}

// Reconcile logic: if namespace has label apm-observe=true, apply instrumentation.yaml into that namespace
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    var ns corev1.Namespace
    if err := r.Get(ctx, types.NamespacedName{Name: req.Name}, &ns); err != nil {
        if errors.IsNotFound(err) {
            return ctrl.Result{}, nil
        }
        return ctrl.Result{}, err
    }

    // ignore namespaces configured
    if r.IgnoreNamespaces != nil && r.IgnoreNamespaces[ns.Name] {
        logger.V(1).Info("namespace is in ignore list, skipping", "namespace", ns.Name)
        return ctrl.Result{}, nil
    }

    labels := ns.Labels
    if labels == nil || strings.ToLower(labels["apm-observe"]) != "true" {
        // nothing to do
        return ctrl.Result{}, nil
    }

    logger.Info("apm-observe label found, applying instrumentation", "namespace", req.Name)

    // Read instrumentation file
    b, err := os.ReadFile(r.InstrPath)
    if err != nil {
        logger.Error(err, "failed to read instrumentation file", "path", r.InstrPath)
        return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
    }

    // Decode multi-document YAML
    dec := yaml.NewYAMLOrJSONDecoder(strings.NewReader(string(b)), 4096)
    for {
        var obj unstructured.Unstructured
        if err := dec.Decode(&obj); err != nil {
            if err == io.EOF {
                break
            }
            logger.Error(err, "failed to decode YAML document")
            return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
        }

        if obj.GetName() == "" {
            // try to fallback to metadata.name inside map
            // if still empty, skip
            continue
        }

        // Set namespace on the object so it gets applied into target namespace
        obj.SetNamespace(req.Name)

        // Try to get existing
        var existing unstructured.Unstructured
        existing.SetGroupVersionKind(obj.GroupVersionKind())
        key := types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.GetName()}
        if err := r.Get(ctx, key, &existing); err != nil {
            if errors.IsNotFound(err) {
                if err := r.Create(ctx, &obj); err != nil {
                    logger.Error(err, "failed to create object", "gvk", obj.GroupVersionKind(), "name", obj.GetName())
                    continue
                }
                logger.Info("created object", "kind", obj.GetKind(), "name", obj.GetName())
                continue
            }
            logger.Error(err, "failed to get existing object")
            continue
        }

        // Update: preserve resourceVersion
        obj.SetResourceVersion(existing.GetResourceVersion())
        if err := r.Update(ctx, &obj); err != nil {
            logger.Error(err, "failed to update object", "name", obj.GetName())
            continue
        }
        logger.Info("updated object", "kind", obj.GetKind(), "name", obj.GetName())
    }

    return ctrl.Result{}, nil
}
