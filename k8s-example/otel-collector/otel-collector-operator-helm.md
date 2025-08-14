# OpenTelemetry Collector Operator via Helm

O OpenTelemetry Collector Operator permite gerenciar coletores como recursos nativos do Kubernetes.

## Instalação

```sh
helm repo add open-telemetry https://github.com/open-telemetry/opentelemetry-helm-charts
helm repo update
helm install opentelemetry-operator open-telemetry/opentelemetry-operator --namespace monitoring
helm install otel-collector open-telemetry/opentelemetry-collector -n monitoring -f values.yaml
```

