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
helm install auto-ns charts/auto-ns-opentelemetry-instrumentation --namespace auto-observe
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
docker build -f Dockerfile.amd64 -t ghcr.io/<user>/auto-ns:0.0.1 .
docker push ghcr.io/<user>/auto-ns:0.0.1
```

#### Helm Upgrade

```bash
helm upgrade --install auto-ns charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe \
  --set image.tag=0.0.1
```

### Debug

Operator logları:

```bash
kubectl -n auto-observe logs -l app=auto-ns-opentelemetry-instrumentation -f
```

---

## English

### What is it?

A Kubernetes operator that automatically watches namespaces and applies OpenTelemetry Instrumentation resources to namespaces labeled with `apm-observe=true`.

### Usage

#### 1. Install

```bash
kubectl create namespace auto-observe
helm install auto-ns charts/auto-ns-opentelemetry-instrumentation --namespace auto-observe
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
docker build -f Dockerfile.amd64 -t ghcr.io/<user>/auto-ns:0.0.1 .
docker push ghcr.io/<user>/auto-ns:0.0.1
```

#### Helm Upgrade

```bash
helm upgrade --install auto-ns charts/auto-ns-opentelemetry-instrumentation \
  --namespace auto-observe \
  --set image.tag=0.0.1
```

### Debug

Operator logs:

```bash
kubectl -n auto-observe logs -l app=auto-ns-opentelemetry-instrumentation -f
```
