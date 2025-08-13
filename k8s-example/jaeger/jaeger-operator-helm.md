# Jaeger Operator via Helm

O Jaeger Operator facilita o deploy e gestão do Jaeger no Kubernetes.

## Instalação

```sh
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm repo update
helm install jaeger-operator jaegertracing/jaeger-operator --namespace observability --create-namespace
```

## Criando uma instância Jaeger (All-in-One)

Crie um arquivo `jaeger.yaml`:

```yaml
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: jaeger-all-in-one
  namespace: observability
spec:
  strategy: allInOne
  allInOne:
    options:
      log-level: debug
```

E aplique:

```sh
kubectl apply -f jaeger.yaml
```
