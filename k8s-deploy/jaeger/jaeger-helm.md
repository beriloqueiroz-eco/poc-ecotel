# Jaeger Operator via Helm

O Jaeger Operator facilita o deploy e gestão do Jaeger no Kubernetes.

## Instalação com operador

```sh
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm repo update
helm install jaeger-operator jaegertracing/jaeger-operator --namespace monitoring
kubectl apply -f jaeger.yaml
```

## Instalação sem operador

```sh
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm repo update
helm install jaeger jaegertracing/jaeger --namespace monitoring -f values.yaml
```

## acesso

- <http://jaeger-query.monitoring.svc.cluster.local:16686/>
