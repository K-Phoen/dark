# DARK

[![codecov](https://codecov.io/gh/K-Phoen/dark/branch/master/graph/badge.svg)](https://codecov.io/gh/K-Phoen/dark)

**D**ashboards **A**s **R**esources in **K**ubernetes.

DARK provides a way to define and deploy Grafana dashboards via Kubernetes, next to the services they monitor.

If you are looking for a way to version your dashboards and deploy them across all environments, like you would do
with your services, then this project is meant for you.

In addition to dashboards, DARK can also manage:

* data sources
* [API keys](./docs/usage/creating-api-keys.md)
* alert manager configuration (alert routing, notification policies, ...)

## Design goals

* full description of dashboards via YAML
* integration with `kubectl`
* seamless integration with Grafana
* delegate YAML decoding and dashboard generation to [Grabana](https://github.com/K-Phoen/grabana)

## Non-goals

* creation and operation of Grafana instances themselves. For that, check out [`grafana-operator/grafana-operator`](https://github.com/grafana-operator/grafana-operator)

## Example dashboard

```yaml
apiVersion: k8s.kevingomez.fr/v1
kind: GrafanaDashboard
metadata:
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
        - timeseries:
            title: HTTP Rate
            datasource: prometheus-default
            targets:
              - prometheus:
                  query: "rate(promhttp_metric_handler_requests_total[$interval])"
                  legend: "{{ handler }} - {{ code }}"
        
        - graph:
            title: Heap allocations
            datasource: prometheus-default
            targets:
              - prometheus:
                  query: "go_memstats_heap_alloc_bytes"
                  legend: "{{ job }}"
```

More examples can be found in the [`examples/`](./examples) folder

## Installation & usage

Check out [the documentation](docs/index.md) to dig deeper into how to set up and use DARK.

## Adopters

[Companies using DARK](ADOPTERS.md).

## License

This library is under the [MIT](LICENSE) license.