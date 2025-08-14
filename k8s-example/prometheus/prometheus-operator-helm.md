# Prometheus Operator via Helm

O Prometheus Operator é recomendado para ambientes produtivos, pois gerencia Prometheus, Alertmanager e regras de monitoramento.

## Instalação

```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack --namespace observability --create-namespace
```

