# Creating dashboards

## Modelling a dashboard

```yaml
# k8s/example-dashboard.yml
apiVersion: k8s.kevingomez.fr/v1
kind: GrafanaDashboard

metadata:
  # must be unique across dashboards
  name: example-dashboard
  namespace: monitoring

folder: "Awesome folder"
spec:
  title: Awesome dashboard

  shared_crosshair: true
  tags: [generated, yaml]
  auto_refresh: 10s

  variables:
    - interval:
        name: interval
        label: interval
        default: 1m
        values: [30s, 1m, 5m, 10m, 30m, 1h, 6h, 12h]

  rows:
    - name: Prometheus
      panels:
        - graph:
            title: HTTP Rate
            height: 400px
            datasource: prometheus-default
            targets:
              - prometheus:
                  query: "rate(promhttp_metric_handler_requests_total[$interval])"
                  legend: "{{handler}} - {{ code }}"
        - graph:
            title: Heap allocations
            height: 400px
            datasource: prometheus-default
            targets:
              - prometheus:
                  query: "go_memstats_heap_alloc_bytes"
                  legend: "{{job}}"
```

For more information on the YAML schema used to describe dashboards, see [Grabana](https://github.com/K-Phoen/grabana/blob/master/doc/index.md#dashboards-as-yaml).

## Deploying a dashboard

DARK dashboards are deployed like any other Kubernetes manifest:

```sh
kubectl apply -f k8s/example-dashboard.yml
```

Verify that the dashboard was created:

```sh
kubectl get dashboards example-dashboard
kubectl get events | grep example-dashboard
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)