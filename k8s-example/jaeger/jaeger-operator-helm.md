# Jaeger Operator via Helm

O Jaeger Operator facilita o deploy e gestão do Jaeger no Kubernetes.

## Instalação

```sh
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm repo update
helm install jaeger-operator jaegertracing/jaeger-operator --namespace observability --create-namespace
kubectl apply -f jaeger.yaml
```
