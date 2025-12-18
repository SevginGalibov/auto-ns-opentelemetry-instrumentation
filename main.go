package main

import (
    "context"
    "flag"
    "os"
    "strings"
    "time"

    "k8s.io/apimachinery/pkg/runtime"
    utilruntime "k8s.io/apimachinery/pkg/util/runtime"
    clientgoscheme "k8s.io/client-go/kubernetes/scheme"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"

    "github.com/SevginGalibov/auto-ns-opentelemetry-instrumentation/controllers"
)

var (
    scheme = runtime.NewScheme()
)

func init() {
    utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
    var metricsAddr string
    var enableLeaderElection bool
    flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
    flag.BoolVar(&enableLeaderElection, "leader-elect", false,
        "Enable leader election for controller manager.")
    flag.Parse()

    ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:             scheme,
        MetricsBindAddress: metricsAddr,
        LeaderElection:     enableLeaderElection,
        LeaderElectionID:   "auto-ns-opentelemetry-instrumentation",
    })
    if err != nil {
        ctrl.Log.Error(err, "unable to start manager")
        os.Exit(1)
    }

    instrPath := os.Getenv("INSTR_PATH")
    if instrPath == "" {
        instrPath = "/etc/instrumentation/instrumentation.yaml"
    }

    // IGNORE_NAMESPACES env: comma separated list
    ignoreEnv := os.Getenv("IGNORE_NAMESPACES")
    ignoreMap := make(map[string]bool)
    if ignoreEnv != "" {
        for _, n := range strings.Split(ignoreEnv, ",") {
            n = strings.TrimSpace(n)
            if n != "" {
                ignoreMap[n] = true
            }
        }
    }

    if err = (&controllers.NamespaceReconciler{
        Client:           mgr.GetClient(),
        Scheme:           mgr.GetScheme(),
        InstrPath:        instrPath,
        IgnoreNamespaces: ignoreMap,
    }).SetupWithManager(mgr); err != nil {
        ctrl.Log.Error(err, "unable to create controller", "controller", "Namespace")
        os.Exit(1)
    }

    ctx := ctrl.SetupSignalHandler()
    ctrl.Log.Info("starting manager")
    if err := mgr.Start(ctx); err != nil {
        ctrl.Log.Error(err, "problem running manager")
        os.Exit(1)
    }

    // Ensure main doesn't exit
    <-context.Background().Done()
    time.Sleep(time.Second)
}
