# Jaeger (Kubernetes)

Use o Helm para instalar o Jaeger Operator:

```sh
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm repo update
helm install jaeger-operator jaegertracing/jaeger-operator --namespace observability --create-namespace
```

Depois, crie o recurso Jaeger (all-in-one) via manifest YAML ou Helm values.
