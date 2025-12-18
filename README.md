# auto-ns-opentelemetry-instrumentation

**[English](#english) | [Türkçe](#türkçe)**

---

## Türkçe

### Nedir?

Kubernetes namespace'lerini otomatik olarak izleyen ve `apm-observe=true` etiketine sahip namespace'lere OpenTelemetry Instrumentation kaynağını uygulayan bir operatördür.

### Kullanım

#### 1. Kurulum

```bash
kubectl create namespace auto-observe
helm install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation --namespace auto-observe
```

#### 2. Namespace'i Etikitle

```bash
kubectl create namespace my-app
kubectl label namespace my-app apm-observe=true
```

Instrumentation otomatik olarak uygulanacaktır:

```bash
kubectl get instrumentation -n my-app
```

### Konfigürasyon

Instrumentation'ın env değerlerini `charts/auto-ns-opentelemetry-instrumentation/values.yaml` dosyasında düzenle:

```yaml
instrumentation:
  exporter:
    endpoint: "http://collector.svc:4318"
  env:
    - name: OTEL_LOGS_EXPORTER
      value: "otlp"
  java:
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: "http://collector.svc:4318"
```

Yoksayılacak namespace'leri ayarla:

```yaml
ignoreNamespaces:
  - kube-system
  - kube-public
  - default
```

### Build & Deploy

#### Image Build (AMD64)

```bash
GOOS=linux GOARCH=amd64 go build -o operator-amd64 ./
docker build -f Dockerfile.amd64 -t ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.1 .
docker push ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.1
```

#### Helm Upgrade

```bash
helm upgrade --install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe \
  --set image.tag=0.0.1
# auto-ns-opentelemetry-instrumentation

**[English](#english) | [Türkçe](#türkçe)**

---

## Türkçe

### Nedir?

Kubernetes namespace'lerini otomatik olarak izleyen ve `apm-observe=true` etiketine sahip namespace'lere OpenTelemetry Instrumentation kaynağını uygulayan bir operatördür.

### Kullanım

#### 1. Kurulum

```bash
kubectl create namespace auto-observe
helm install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation --namespace auto-observe
```

#### 2. Namespace'i Etikitle

```bash
kubectl create namespace my-app
kubectl label namespace my-app apm-observe=true
```

Instrumentation otomatik olarak uygulanacaktır:

```bash
kubectl get instrumentation -n my-app
```

### Konfigürasyon

Instrumentation'ın env değerlerini `charts/auto-ns-opentelemetry-instrumentation/values.yaml` dosyasında düzenle:

```yaml
instrumentation:
  exporter:
    endpoint: "http://collector.svc:4318"
  env:
    - name: OTEL_LOGS_EXPORTER
      value: "otlp"
  java:
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: "http://collector.svc:4318"
```

Yoksayılacak namespace'leri ayarla:

```yaml
ignoreNamespaces:
  - kube-system
  - kube-public
  - default
```

### Build & Deploy

#### Image Build (AMD64)

```bash
GOOS=linux GOARCH=amd64 go build -o operator-amd64 ./
docker build -f Dockerfile.amd64 -t ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.1 .
docker push ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.1
```

#### Helm Upgrade

```bash
helm upgrade --install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe \
  --set image.tag=0.0.1
```

### Debug

Operator logları:

```bash
kubectl -n auto-observe logs -l app=auto-ns-opentelemetry-instrumentation -f
```

---

## SigNoz Entegrasyonu (Türkçe)

Aşağıdaki adımlar SigNoz platformunu cluster ortamınıza kurmak ve OpenTelemetry bileşenlerini hazırlamak içindir.

1) Cert-Manager kurun:

```bash
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.7.1/cert-manager.yaml
```

2) OpenTelemetry Operator'u kurun:

```bash
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
helm install opentelemetry-operator open-telemetry/opentelemetry-operator \
  --set "manager.collectorImage.repository=otel/opentelemetry-collector-k8s" \
  --namespace opentelemetry-operator-system --create-namespace
```

3) SigNoz platform namespace'i oluşturun ve SigNoz'u yükleyin (values dosyası repo içindeki `signoz/signoz-values.yaml`):

```bash
kubectl create namespace platform
helm repo add signoz https://charts.signoz.io || true
helm repo update
helm install dmsignoz signoz/signoz -n platform -f signoz/signoz-values.yaml
```

4) Doğrulama:

```bash
kubectl -n platform get pods
kubectl -n opentelemetry-operator-system get pods
```

Not: `signoz/signoz-values.yaml` dosyası repository kökünde `signoz/` klasöründe bulunur. Gerekli endpoint/secret ayarlarını bu dosyada düzenleyin.

---

## English

### What is it?

A Kubernetes operator that automatically watches namespaces and applies OpenTelemetry Instrumentation resources to namespaces labeled with `apm-observe=true`.

### Usage

#### 1. Install

```bash
kubectl create namespace auto-observe
helm install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation --namespace auto-observe
```

#### 2. Label Namespace

```bash
kubectl create namespace my-app
kubectl label namespace my-app apm-observe=true
```

Instrumentation will be applied automatically:

```bash
kubectl get instrumentation -n my-app
```

### Configuration

Edit instrumentation env values in `charts/auto-ns-opentelemetry-instrumentation/values.yaml`:

```yaml
instrumentation:
  exporter:
    endpoint: "http://collector.svc:4318"
  env:
    - name: OTEL_LOGS_EXPORTER
      value: "otlp"
  java:
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: "http://collector.svc:4318"
```

Set namespaces to ignore:

```yaml
ignoreNamespaces:
  - kube-system
  - kube-public
  - default
```

### Build & Deploy

#### Image Build (AMD64)

```bash
GOOS=linux GOARCH=amd64 go build -o operator-amd64 ./
docker build -f Dockerfile.amd64 -t ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.1 .
docker push ghcr.io/<user>/auto-ns-opentelemetry-instrumentation:0.0.1
```

#### Helm Upgrade

```bash
helm upgrade --install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe \
  --set image.tag=0.0.1
```

### SigNoz Integration (English)

The following steps install SigNoz and prepare OpenTelemetry components on your cluster.

1) Install cert-manager:

```bash
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.7.1/cert-manager.yaml
```

2) Install OpenTelemetry Operator:

```bash
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
helm install opentelemetry-operator open-telemetry/opentelemetry-operator \
  --set "manager.collectorImage.repository=otel/opentelemetry-collector-k8s" \
  --namespace opentelemetry-operator-system --create-namespace
```

3) Create the platform namespace and install SigNoz (values file is under `signoz/signoz-values.yaml` in this repo):

```bash
kubectl create namespace platform
helm repo add signoz https://charts.signoz.io || true
helm repo update
helm install dmsignoz signoz/signoz -n platform -f signoz/signoz-values.yaml
```

4) Verify installations:

```bash
kubectl -n platform get pods
kubectl -n opentelemetry-operator-system get pods
```

Note: `signoz/signoz-values.yaml` is located under the `signoz/` directory in this repository. Edit exporter endpoints and secrets in that values file before installing.
