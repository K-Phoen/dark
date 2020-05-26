# DARK

[![codecov](https://codecov.io/gh/K-Phoen/dark/branch/master/graph/badge.svg)](https://codecov.io/gh/K-Phoen/dark)

**D**ashboards **A**s **R**esources in **K**ubernetes.

DARK provides a way to define and deploy Grafana dashboards via Kubernetes, next to the services they monitor.

If you are looking for a way to version your dashboards and deploy them across all environments, like you would do
with your services, then this project is meant for you.

## Design goals

* full description of dashboards via YAML
* compatibility with `kubectl`
* seamless integration with Grafana
* delegate YAML decoding and dashboard generation to [Grabana](https://github.com/K-Phoen/grabana)

## Setup

**Note:** review these manifests to ensure that they fit your cluster's configuration.

Setup the [CRD](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/):

```sh
kubectl apply -f k8s/crd.yaml
```

Create a secret to store Grafana's API token (with `editor` access level):

```sh
kubectl create secret generic dark-tokens --from-literal=grafana=TOKEN_HERE
```

Deploy DARK's controller:

```sh
kubectl apply -f k8s/deployment.yaml
```

## Dashboard definition

Define a dashboard:

```yaml
# k8s/example-dashboard.yml
apiVersion: k8s.kevingomez.fr/v1
kind: GrafanaDashboard

metadata:
  # must be unique across dashboards
  name: example-dashboard

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

For more information on the YAML schema used to describe dashboards, see [Grabana](https://github.com/K-Phoen/grabana).

Apply the configuration:

```sh
kubectl apply -f k8s/example-dashboard.yml
```

And verify that the dashboard was created:

```sh
kubectl get dashboards
kubectl get events | grep dark
```

## Converting Grafana JSON dashboard to YAML

To ease the transition from existing, raw Grafana dashboards to DARK, a converter is provided.
It takes the path to a JSON dashboard and a path for the destination YAML file.

```sh
docker run --rm -it -v $(pwd):/workspace kphoen/dark-converter:latest dashboard.json converted-dashboard.yaml
```

## License

This library is under the [MIT](LICENSE) license.
