package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/decoder"
	"gopkg.in/yaml.v3"
)

type Datasources struct {
	grabanaClient *grabana.Client
}

func NewDatasources(grabanaClient *grabana.Client) *Datasources {
	return &Datasources{grabanaClient: grabanaClient}
}

func (datasources *Datasources) RawSpecToModel(rawJSON []byte) (datasource.Datasource, error) {
	spec := make(map[string]interface{})
	if err := json.Unmarshal(rawJSON, &spec); err != nil {
		return nil, fmt.Errorf("could not unmarshall datasource json spec: %w", err)
	}

	datasourceYaml, err := yaml.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("could not convert datasource spec to yaml: %w", err)
	}

	model, err := decoder.UnmarshalYAMLDatasource(bytes.NewBuffer(datasourceYaml))
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall datasource YAML spec: %w", err)
	}

	return model, nil
}

func (datasources *Datasources) Upsert(ctx context.Context, model datasource.Datasource) error {
	return datasources.grabanaClient.UpsertDatasource(ctx, model)
}

func (datasources *Datasources) Delete(_ context.Context, _ datasource.Datasource) error {
	// TODO
	return nil
}
