package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertHeatmapPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "400px"
	datasource := "prometheus"

	heatmapPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "heatmap panel",
			Type:        "heatmap",
			Description: strPtr("heatmap description"),
			Transparent: true,
			Height:      &height,
			Datasource:  &sdk.DatasourceRef{LegacyName: datasource},
		},
		HeatmapPanel: &sdk.HeatmapPanel{
			HideZeroBuckets: true,
			HighlightCards:  true,
			ReverseYBuckets: true,
			DataFormat:      "tsbuckets",
		},
	}

	converted, ok := converter.convertDataPanel(heatmapPanel)

	req.True(ok)
	req.True(converted.Heatmap.Transparent)
	req.Equal("heatmap panel", converted.Heatmap.Title)
	req.Equal("heatmap description", converted.Heatmap.Description)
	req.Equal(height, converted.Heatmap.Height)
	req.Equal(datasource, converted.Heatmap.Datasource)
	req.True(converted.Heatmap.ReverseYBuckets)
	req.True(converted.Heatmap.HideZeroBuckets)
	req.True(converted.Heatmap.HighlightCards)
	req.Equal("time_series_buckets", converted.Heatmap.DataFormat)
}
