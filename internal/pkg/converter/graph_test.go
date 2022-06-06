package converter

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestConvertGraphPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "400px"
	datasource := "prometheus"

	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "graph panel",
			Type:        "graph",
			Description: strPtr("graph description"),
			Transparent: true,
			Height:      &height,
			Datasource:  &sdk.DatasourceRef{LegacyName: datasource},
		},
		GraphPanel: &sdk.GraphPanel{},
	}

	converted, ok := converter.convertDataPanel(graphPanel)

	req.True(ok)
	req.NotNil(converted.Graph)

	graph := converted.Graph
	req.True(graph.Transparent)
	req.Equal("graph panel", graph.Title)
	req.Equal("graph description", graph.Description)
	req.Equal(height, graph.Height)
	req.Equal(datasource, graph.Datasource)
}

func TestConvertGraphLinks(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Type: "graph",
			Links: []sdk.Link{
				{Title: "title", URL: strPtr("url")},
			},
		},
		GraphPanel: &sdk.GraphPanel{},
	}

	converted, ok := converter.convertDataPanel(graphPanel)

	req.True(ok)
	req.NotNil(converted.Graph)

	graph := converted.Graph
	req.Len(graph.Links, 1)
	req.Equal("title", graph.Links[0].Title)
}

func TestConvertGraphLegend(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	rawLegend := sdk.Legend{
		AlignAsTable: true,
		Avg:          true,
		Current:      true,
		HideEmpty:    true,
		HideZero:     true,
		Max:          true,
		Min:          true,
		RightSide:    true,
		Show:         true,
		Total:        true,
	}

	legend := converter.convertGraphLegend(rawLegend)

	req.ElementsMatch(
		[]string{"as_table", "to_the_right", "min", "max", "avg", "current", "total", "no_null_series", "no_zero_series"},
		legend,
	)
}

func TestConvertGraphCanHideLegend(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	legend := converter.convertGraphLegend(sdk.Legend{Show: false})
	req.ElementsMatch([]string{"hide"}, legend)
}

func TestConvertGraphAxis(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	rawAxis := sdk.Axis{
		Format:  "bytes",
		LogBase: 2,
		Min:     &sdk.FloatString{Value: 0},
		Max:     &sdk.FloatString{Value: 42},
		Show:    true,
		Label:   "Axis",
	}

	axis := converter.convertGraphAxis(rawAxis)

	req.Equal("bytes", *axis.Unit)
	req.Equal("Axis", axis.Label)
	req.EqualValues(0, *axis.Min)
	req.EqualValues(42, *axis.Max)
	req.False(*axis.Hidden)
}

func TestConvertGraphVisualization(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())
	enabled := true

	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "graph panel",
			Type:  "graph",
		},
		GraphPanel: &sdk.GraphPanel{
			NullPointMode: "connected",
			SteppedLine:   true,
			SeriesOverrides: []sdk.SeriesOverride{
				{
					Alias:  "alias",
					Dashes: &enabled,
				},
			},
		},
	}

	visualization := converter.convertGraphVisualization(graphPanel)

	req.True(visualization.Staircase)
	req.Equal("connected", visualization.NullValue)
	req.Len(visualization.Overrides, 1)
}

func TestConvertGraphOverridesWithNoOverride(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "graph panel",
			Type:  "graph",
		},
		GraphPanel: &sdk.GraphPanel{},
	}

	overrides := converter.convertGraphOverrides(graphPanel)

	req.Len(overrides, 0)
}

func TestConvertGraphOverridesWithOneOverride(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())
	color := "red"
	enabled := true
	number := 2

	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "heatmap panel",
			Type:  "graph",
		},
		GraphPanel: &sdk.GraphPanel{
			SeriesOverrides: []sdk.SeriesOverride{
				{
					Alias:  "alias",
					Color:  &color,
					Dashes: &enabled,
					Fill:   &number,
					Lines:  &enabled,
				},
			},
		},
	}

	overrides := converter.convertGraphOverrides(graphPanel)

	req.Len(overrides, 1)

	override := overrides[0]

	req.Equal("alias", override.Alias)
	req.Equal(color, override.Color)
	req.True(*override.Dashes)
	req.True(*override.Lines)
	req.Equal(number, *override.Fill)
}
