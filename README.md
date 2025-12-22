# Auto Namespace OpenTelemetry Instrumentation

A lightweight Kubernetes operator that automates the injection of OpenTelemetry Instrumentation resources. It watches for namespaces labeled with `apm-observe=true` and applies the necessary configuration to enable monitoring without manual intervention.

## 🚀 Quick Overview

- **What it does:** Automatically creates `Instrumentation` resources in specific namespaces.
- **Trigger:** Namespaces labeled with `apm-observe=true`.
- **Stack:** Designed to work with the OpenTelemetry Operator and SigNoz (or any OTLP backend).

---

## 🛠 Prerequisites & Dependencies

Before installing the operator, ensure your cluster has the necessary monitoring stack dependencies. Follow this order:

### 1. Install Cert-Manager
Required for the OpenTelemetry Operator..

```bash
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.7.1/cert-manager.yaml
```

### 2. Install OpenTelemetry Operator

```bash
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts

helm install opentelemetry-operator open-telemetry/opentelemetry-operator \
  --set "manager.collectorImage.repository=otel/opentelemetry-collector-k8s" \
  --namespace opentelemetry-operator-system \
  --create-namespace
```

### 3. Install SigNoz (Backend)

*Note: This repository includes a custom values file at `signoz/signoz-values.yaml`. Please edit this file to configure your storage or credentials before running the commands below.*

```bash
kubectl create namespace platform
helm repo add signoz https://charts.signoz.io || true
helm repo update

# Install using the custom values file included in this repo
helm install dmsignoz signoz/signoz -n platform -f signoz/signoz-values.yaml
```

---

## 📦 Installation

You can install the auto-instrumentation operator via the Helm repository (recommended) or from the local source.

### Option A: From Helm Repository (Recommended)

```bash
kubectl create namespace auto-observe
helm repo add sevgingalibov https://SevginGalibov.github.io/auto-ns-opentelemetry-instrumentation/
helm repo update

helm install auto-ns-opentelemetry-instrumentation sevgingalibov/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe \
  --create-namespace \
  --set image.tag=0.0.4
```

### Option B: Local Development Chart

```bash
kubectl create namespace auto-observe
helm install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe \
  --set image.tag=0.0.4
```

---

## 💡 Usage

Once the operator is running, you only need to label a namespace to enable auto-instrumentation.

### 1. Label a Namespace
Create your application namespace and add the `apm-observe=true` label:

```bash
kubectl create namespace my-app
kubectl label namespace my-app apm-observe=true
```

### 2. Verify Injection
The operator will automatically detect the label and create the resource:

```bash
kubectl -n my-app get instrumentation
```

### 3. Annotate Your Workloads

Finally, simply add the language-specific annotation to your Pods or Deployments. The OpenTelemetry Operator will detect this and inject the agent automatically.

```bash
Language	Annotation
.NET	instrumentation.opentelemetry.io/inject-dotnet: "true"
Java	instrumentation.opentelemetry.io/inject-java: "true"
Node.js	instrumentation.opentelemetry.io/inject-nodejs: "true"
Python	instrumentation.opentelemetry.io/inject-python: "true"
```
---

## ⚙️ Configuration

You can customize the instrumentation settings (e.g., OTLP endpoints or ignored namespaces) by editing the `values.yaml` file in the Helm chart.

**Example `values.yaml` snippet:**

```yaml
instrumentation:
  exporter:
    endpoint: "http://dmsignoz-otel-collector.platform.svc:4318"
  env:
    - name: OTEL_LOGS_EXPORTER
      value: "otlp"
  java:
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: "http://dmsignoz-otel-collector.platform.svc:4318"

# Namespaces to explicitly ignore (even if labeled)
ignoreNamespaces:
  - kube-system
  - kube-public
  - default
```

To update configuration on a running cluster:
```bash
helm upgrade auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe
```

---

## 🏗 Development & Build

Instructions for building the operator image manually.

### Build Multi-Arch Image (AMD64 & ARM64)
Use `docker buildx` to create a manifest for multiple architectures.

```bash
# Ensure buildx is ready
docker buildx create --use || true

# Build and push
docker buildx build --platform linux/amd64,linux/arm64 \
  -f Dockerfile.amd64 \
  -t ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.4 \
  --push .
```

### Simple AMD64 Build

```bash
GOOS=linux GOARCH=amd64 go build -o operator-amd64 ./
docker build -f Dockerfile.amd64 -t ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.4 .
docker push ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.4
```

---

## 🔍 Troubleshooting

**Check Operator Logs:**
```bash
kubectl -n auto-observe logs -l app=auto-ns-opentelemetry-instrumentation -f
```

**Check Platform Status:**
```bash
kubectl -n platform get pods
kubectl -n opentelemetry-operator-system get pods
```
