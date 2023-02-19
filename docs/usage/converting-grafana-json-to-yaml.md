# Converting a Grafana JSON dashboard to YAML

To ease the transition from existing, raw Grafana dashboards to DARK, a converter is provided.
It takes a dashboard as modelled by Grafana and converts it to a YAML file compatible with DARK.

```sh
docker run --rm -it -u $(id -u):$(id -g) -v $(pwd):/workspace kphoen/dark-converter:latest \
    convert-k8s-manifest \
        -i dashboard.json \ # input file
        -o converted-dashboard.yaml \ # output file
        --folder Dark \ # Folder in Grafana in which the dashboard should be created
        --namespace monitoring \ # Namespace in Kubernetes in which the manifest should live
        test-dashboard # Name of the Kubernetes manifest
```

## That was it!

[Return to the index to explore what you can do with DARK](../index.md)
