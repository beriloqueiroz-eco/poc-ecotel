# OpenTelemetry Collector Operator via Helm

O OpenTelemetry Collector Operator permite gerenciar coletores como recursos nativos do Kubernetes.

## Instalação

```sh
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
helm install opentelemetry-operator open-telemetry/opentelemetry-operator --namespace observability --create-namespace
```

## Criando um Collector

Crie um arquivo `otel-collector.yaml` e aplique:

```sh
kubectl apply -f otel-collector.yaml
```
