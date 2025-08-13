# OpenTelemetry Collector (Kubernetes)

Use o Helm para instalar o OpenTelemetry Collector:

```sh
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
helm install otel-collector open-telemetry/opentelemetry-collector --namespace observability --create-namespace
```

Personalize o values.yaml para configurar receivers, exporters e pipelines.
