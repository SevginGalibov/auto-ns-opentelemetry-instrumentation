# auto-ns-opentelemetry-instrumentation

**[English](#english) | [Türkçe](#türkçe)**

---

## Quick overview

This repository contains a small Kubernetes operator that applies OpenTelemetry Instrumentation resources into namespaces labeled with `apm-observe=true`. It also includes a Helm chart and example values to install the operator and the instrumentation config.

This README provides a concise, production-oriented Quick Start for installing the monitoring stack and then deploying the operator. Follow the order below.

Order of installation:
- Cert-Manager
- OpenTelemetry Operator
- SigNoz stack
- auto-ns-opentelemetry-instrumentation

---

## Türkçe

Kısa: Kurulum adımları gerektiği sırada ve doğrulama komutlarıyla aşağıdadır.

Önkoşullar
- `kubectl`, `helm` yüklü ve cluster erişiminiz var.

1) Cert-Manager

```bash
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.7.1/cert-manager.yaml
kubectl -n cert-manager get pods
```

2) OpenTelemetry Operator

```bash
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
helm install opentelemetry-operator open-telemetry/opentelemetry-operator \
  --set "manager.collectorImage.repository=otel/opentelemetry-collector-k8s" \
  --namespace opentelemetry-operator-system --create-namespace
kubectl -n opentelemetry-operator-system get pods
```

3) SigNoz stack (platform namespace)

Edit: repository içinde `signoz/signoz-values.yaml` dosyasını ihtiyaçlarınıza göre düzenleyin (endpoints, storage, secrets).

```bash
kubectl create namespace platform
helm repo add signoz https://charts.signoz.io || true
helm repo update
helm install dmsignoz signoz/signoz -n platform -f signoz/signoz-values.yaml
kubectl -n platform get pods
```

4) auto-ns-opentelemetry-instrumentation (operator)

Chart ve config repo içindedir. Operator image parametresini `values.yaml` içinde kontrol edin.

Not: iki kurulum seçeneği vardır — (A) yerel chart dizininden doğrudan kurulum (geliştirme/test) ve (B) yayınlanmış Helm repo (gh-pages) üzerinden kurulum (kullanıcılar için kolay, production).

A) Yerel chart (mevcut repo içindeki `charts/` dizininden):

```bash
helm install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe --create-namespace --set image.tag=0.0.1
```

B) Yayınlanmış Helm repo üzerinden (ör. GitHub Pages `gh-pages` kullanıyorsanız):

```bash
helm repo add sevgingalibov https://SevginGalibov.github.io/auto-ns-opentelemetry-instrumentation/
helm repo update
helm install auto-ns-opentelemetry-instrumentation sevgingalibov/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe --create-namespace --set image.tag=0.0.1
```

# Label a test namespace to trigger instrumentation
kubectl create namespace my-app
kubectl label namespace my-app apm-observe=true
kubectl -n my-app get instrumentation

Doğrulama

```bash
kubectl -n auto-observe get pods
kubectl -n platform get pods
kubectl -n opentelemetry-operator-system get pods
```

Notlar
- `signoz/signoz-values.yaml` dosyasını root içindeki `signoz/` klasöründe bulabilirsiniz.
- Eğer özel registry, TLS, veya storage ayarları gerekiyorsa values dosyasını düzenleyin.

---

## English

Concise Quick Start for production-like installation. Follow the steps in order.

Prerequisites
- `kubectl` and `helm` configured for your cluster.

1) Install Cert-Manager

```bash
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.7.1/cert-manager.yaml
kubectl -n cert-manager get pods
```

2) Install OpenTelemetry Operator

```bash
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
helm install opentelemetry-operator open-telemetry/opentelemetry-operator \
  --set "manager.collectorImage.repository=otel/opentelemetry-collector-k8s" \
  --namespace opentelemetry-operator-system --create-namespace
kubectl -n opentelemetry-operator-system get pods
```

3) Install SigNoz (platform namespace)

Edit the values file located at `signoz/signoz-values.yaml` in this repository to set endpoints, storage and credentials.

```bash
kubectl create namespace platform
helm repo add signoz https://charts.signoz.io || true
helm repo update
helm install dmsignoz signoz/signoz -n platform -f signoz/signoz-values.yaml
kubectl -n platform get pods
```

4) Install auto-ns-opentelemetry-instrumentation operator

The chart and configuration are included in this repo. Verify operator image settings in the chart `values.yaml`.

Note: there are two installation options — (A) local chart (development/testing) and (B) installed from a published Helm repository (recommended for consumers).

A) Local chart (install from the `charts/` directory in this repository):

```bash
helm install auto-ns-opentelemetry-instrumentation charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe --create-namespace --set image.tag=0.0.1
```

B) From published Helm repo (e.g. GitHub Pages `gh-pages`):

```bash
helm repo add sevgingalibov https://SevginGalibov.github.io/auto-ns-opentelemetry-instrumentation/
helm repo update
helm install auto-ns-opentelemetry-instrumentation sevgingalibov/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe --create-namespace --set image.tag=0.0.1
```

# Label a namespace to apply instrumentation
kubectl create namespace my-app
kubectl label namespace my-app apm-observe=true
kubectl -n my-app get instrumentation

Verification

```bash
kubectl -n auto-observe get pods
kubectl -n platform get pods
kubectl -n opentelemetry-operator-system get pods
```

Notes
- The SigNoz values file is `signoz/signoz-values.yaml`. Update it before install if you use custom endpoints, storage backends or secrets.
- Keep chart and image versions in sync (bump `Chart.yaml` and `values.yaml` when releasing new operator versions).
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
