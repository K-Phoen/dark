package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertLogsPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "200px"
	datasource := "prometheus"

	logsPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Logs panel",
			Description: strPtr("panel desc"),
			Type:        "logs",
			Transparent: true,
			Height:      &height,
			Datasource:  &sdk.DatasourceRef{LegacyName: datasource},
		},
		LogsPanel: &sdk.LogsPanel{
			Targets: []sdk.Target{
				{
					Expr:         "loki_query",
					LegendFormat: "Legend",
					RefID:        "A",
				},
			},
			Options: sdk.LogsOptions{
				DedupStrategy:      "none",
				WrapLogMessage:     true,
				ShowTime:           true,
				ShowLabels:         true,
				ShowCommonLabels:   true,
				PrettifyLogMessage: true,
				SortOrder:          "Descending",
				EnableLogDetails:   true,
			},
		},
	}

	converted, ok := converter.convertDataPanel(logsPanel)

	req.True(ok)
	req.True(converted.Logs.Transparent)
	req.Equal("Logs panel", converted.Logs.Title)
	req.Equal("panel desc", converted.Logs.Description)
	req.Equal(height, converted.Logs.Height)
	req.Equal(datasource, converted.Logs.Datasource)
	req.False(converted.Logs.Visualization.HideLogDetails)
	req.True(converted.Logs.Visualization.PrettifyJSON)
	req.True(converted.Logs.Visualization.UniqueLabels)
	req.True(converted.Logs.Visualization.WrapLines)
	req.True(converted.Logs.Visualization.Time)
	req.True(converted.Logs.Visualization.CommonLabels)
	req.Equal("none", converted.Logs.Visualization.Deduplication)
	req.Equal("desc", converted.Logs.Visualization.Order)
}

func TestConvertLogsLinks(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	sdkPanel := sdk.NewLogs("")
	sdkPanel.Links = []sdk.Link{
		{Title: "logs title", URL: strPtr("logs url")},
	}

	converted, ok := converter.convertDataPanel(*sdkPanel)

	req.True(ok)
	req.NotNil(converted.Logs)

	panel := converted.Logs
	req.Len(panel.Links, 1)
	req.Equal("logs title", panel.Links[0].Title)
	req.Equal("logs url", panel.Links[0].URL)
}
