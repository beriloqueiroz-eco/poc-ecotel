# todo

- istio enviando dados para o otel
- grafana com relacionamentos by yaml

conteúdo de config do jaeger no grafana

```json
{
    "pdcInjected": false,
    "tracesToLogsV2": {
        "customQuery": true,
        "datasourceUid": "aev06pdbgtp1cc",
        "filterBySpanID": false,
        "filterByTraceID": false,
        "query": "{service_name=~\".+\"} | traceId=\"${__trace.traceId}\"",
        "spanEndTimeShift": "2m",
        "spanStartTimeShift": "-2m",
        "tags": [
            {
                "key": "",
                "value": ""
            }
        ]
    },
    "tracesToMetrics": {
        "datasourceUid": "fev06o3d82ha8c",
        "queries": [
            {
                "name": "Metrics By Trace",
                "query": "{trace_id=\"${__trace.traceId}\"}"
            }
        ],
        "spanEndTimeShift": "2m",
        "spanStartTimeShift": "-2m"
    }
}
```
