# Prometheus (Kubernetes)

Use o Helm para instalar o Prometheus:

```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus --namespace observability --create-namespace
```

Personalize o scrape_config conforme sua stack.
