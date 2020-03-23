package dashboards

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/decoder"
	"gopkg.in/yaml.v2"
)

type Creator struct {
	grabanaClient *grabana.Client
}

func NewCreator(grabanaClient *grabana.Client) *Creator {
	return &Creator{grabanaClient: grabanaClient}
}

func (creator *Creator) FromRawSpec(folderName string, uid string, rawJSON []byte) error {
	spec := make(map[string]interface{})
	if err := json.Unmarshal(rawJSON, &spec); err != nil {
		return fmt.Errorf("could not unmarshall  dashboard json spec: %w", err)
	}

	dashboardYaml, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("could not convert dashboard spec to yaml: %w", err)
	}

	dashboardBuilder, err := decoder.UnmarshalYAML(bytes.NewBuffer(dashboardYaml))
	if err != nil {
		return fmt.Errorf("could not unmarshall dashboard YAML spec: %w", err)
	}

	dashboard.UID(uid)(&dashboardBuilder)

	return creator.upsertDashboard(folderName, dashboardBuilder)
}

func (creator *Creator) Delete(uid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return creator.grabanaClient.DeleteDashboard(ctx, uid)
}

func (creator *Creator) upsertDashboard(folderName string, dashboardBuilder dashboard.Builder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	folder, err := creator.grabanaClient.FindOrCreateFolder(ctx, folderName)
	if err != nil {
		return err
	}

	if _, err := creator.grabanaClient.UpsertDashboard(ctx, folder, dashboardBuilder); err != nil {
		return fmt.Errorf("could not create dashboard: %w", err)
	}

	return nil
}
