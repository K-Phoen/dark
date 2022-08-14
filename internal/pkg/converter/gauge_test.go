package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertGaugePanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "200px"
	datasource := "prometheus"

	gaugePanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Gauge panel",
			Description: strPtr("panel desc"),
			Type:        "gauge",
			Transparent: true,
			Height:      &height,
			Datasource:  &sdk.DatasourceRef{LegacyName: datasource},
		},
		GaugePanel: &sdk.GaugePanel{
			Targets: []sdk.Target{
				{
					Expr: "sum(kube_pod_info{}) / sum(kube_node_status_allocatable{resource=\"pods\"})",
				},
			},
			Options: sdk.StatOptions{
				Orientation: "horizontal",
				ReduceOptions: sdk.ReduceOptions{
					Calcs: []string{"last"},
				},
				Text: &sdk.TextOptions{
					ValueSize: 10,
					TitleSize: 20,
				},
			},
			FieldConfig: sdk.FieldConfig{
				Defaults: sdk.FieldConfigDefaults{
					Unit:       "none",
					Thresholds: sdk.Thresholds{},
				},
			},
		},
	}

	converted, ok := converter.convertDataPanel(gaugePanel)

	req.True(ok)
	req.True(converted.Gauge.Transparent)
	req.Equal("Gauge panel", converted.Gauge.Title)
	req.Equal("panel desc", converted.Gauge.Description)
	req.Equal("none", converted.Gauge.Unit)
	req.Equal("last", converted.Gauge.ValueType)
	req.Equal("horizontal", converted.Gauge.Orientation)
	req.Equal(10, converted.Gauge.ValueFontSize)
	req.Equal(20, converted.Gauge.TitleFontSize)
	req.Equal(height, converted.Gauge.Height)
	req.Equal(datasource, converted.Gauge.Datasource)
}

func TestConvertGaugeLinks(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	sdkPanel := sdk.NewGauge("")
	sdkPanel.Links = []sdk.Link{
		{Title: "stat title", URL: strPtr("stat url")},
	}

	converted, ok := converter.convertDataPanel(*sdkPanel)

	req.True(ok)
	req.NotNil(converted.Gauge)

	panel := converted.Gauge
	req.Len(panel.Links, 1)
	req.Equal("stat title", panel.Links[0].Title)
	req.Equal("stat url", panel.Links[0].URL)
}
