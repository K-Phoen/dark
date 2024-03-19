package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertStatPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "200px"
	datasource := "prometheus"

	statPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Stat panel",
			Description: strPtr("panel desc"),
			Type:        "stat",
			Transparent: true,
			Height:      &height,
			Datasource:  &sdk.DatasourceRef{LegacyName: datasource},
		},
		StatPanel: &sdk.StatPanel{
			Targets: []sdk.Target{
				{
					Expr:         "prometheus_query",
					LegendFormat: "{{ field }}",
					RefID:        "A",
				},
			},
			Options: sdk.StatOptions{
				Orientation: "auto",
				TextMode:    "value_and_name",
				ColorMode:   "value",
				GraphMode:   "area",
				JustifyMode: "auto",
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

	converted, ok := converter.convertDataPanel(statPanel)

	req.True(ok)
	req.True(converted.Stat.Transparent)
	req.Equal("Stat panel", converted.Stat.Title)
	req.Equal("panel desc", converted.Stat.Description)
	req.Equal("none", converted.Stat.Unit)
	req.Equal("last", converted.Stat.ValueType)
	req.Equal("auto", converted.Stat.Orientation)
	req.Equal("value", converted.Stat.ColorMode)
	req.Equal("value_and_name", converted.Stat.Text)
	req.Equal(10, converted.Stat.ValueFontSize)
	req.Equal(20, converted.Stat.TitleFontSize)
	req.Equal(height, converted.Stat.Height)
	req.Equal(datasource, converted.Stat.Datasource)
	req.True(converted.Stat.SparkLine)
}

func TestConvertStatLinks(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	sdkPanel := sdk.NewStat("")
	sdkPanel.Links = []sdk.Link{
		{Title: "stat title", URL: strPtr("stat url")},
	}

	converted, ok := converter.convertDataPanel(*sdkPanel)

	req.True(ok)
	req.NotNil(converted.Stat)

	panel := converted.Stat
	req.Len(panel.Links, 1)
	req.Equal("stat title", panel.Links[0].Title)
	req.Equal("stat url", panel.Links[0].URL)
}
