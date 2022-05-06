package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
	"github.com/K-Phoen/grabana/decoder"
	"gopkg.in/yaml.v3"
)

type Creator struct {
	grabanaClient *grabana.Client
}

func NewCreator(grabanaClient *grabana.Client) *Creator {
	return &Creator{grabanaClient: grabanaClient}
}

func (creator *Creator) FromRawSpec(ctx context.Context, folderName string, uid string, rawJSON []byte) error {
	spec := make(map[string]interface{})
	if err := json.Unmarshal(rawJSON, &spec); err != nil {
		return fmt.Errorf("could not unmarshall dashboard json spec: %w", err)
	}

	dashboardYaml, err := yaml.Marshal(spec)
	if err != nil {
		return fmt.Errorf("could not convert dashboard spec to yaml: %w", err)
	}

	dashboardBuilder, err := decoder.UnmarshalYAML(bytes.NewBuffer(dashboardYaml))
	if err != nil {
		return fmt.Errorf("could not unmarshall dashboard YAML spec: %w", err)
	}

	if err := dashboard.UID(uid)(&dashboardBuilder); err != nil {
		return fmt.Errorf("could not set dashboard UID: %w", err)
	}

	return creator.upsertDashboard(ctx, folderName, dashboardBuilder)
}

func (creator *Creator) Delete(ctx context.Context, uid string) error {
	err := creator.grabanaClient.DeleteDashboard(ctx, uid)
	if err != nil && err != grabana.ErrDashboardNotFound {
		return err
	}

	return nil
}

func (creator *Creator) upsertDashboard(ctx context.Context, folderName string, dashboardBuilder dashboard.Builder) error {
	folder, err := creator.grabanaClient.FindOrCreateFolder(ctx, folderName)
	if err != nil {
		return err
	}

	if _, err := creator.grabanaClient.UpsertDashboard(ctx, folder, dashboardBuilder); err != nil {
		return fmt.Errorf("could not create dashboard: %w", err)
	}

	return nil
}
