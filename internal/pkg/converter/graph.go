package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
)

func (converter *JSON) convertGraph(panel sdk.Panel) grabana.DashboardPanel {
	graph := &grabana.DashboardGraph{
		Title:       panel.Title,
		Span:        panelSpan(panel),
		Transparent: panel.Transparent,
		Axes: &grabana.GraphAxes{
			Bottom: converter.convertGraphAxis(panel.Xaxis),
		},
		Legend:        converter.convertGraphLegend(panel.GraphPanel.Legend),
		Visualization: converter.convertGraphVisualization(panel),
	}

	if panel.Description != nil {
		graph.Description = *panel.Description
	}
	if panel.Repeat != nil {
		graph.Repeat = *panel.Repeat
	}
	if panel.RepeatDirection != nil {
		graph.RepeatDirection = sdkRepeatDirectionToYAML(*panel.RepeatDirection)
	}
	if panel.Height != nil {
		graph.Height = *(panel.Height).(*string)
	}
	if panel.Datasource != nil {
		graph.Datasource = panel.Datasource.LegacyName
	}
	if len(panel.Links) != 0 {
		graph.Links = converter.convertPanelLinks(panel.Links)
	}

	if len(panel.Yaxes) == 2 {
		graph.Axes.Left = converter.convertGraphAxis(panel.Yaxes[0])
		graph.Axes.Right = converter.convertGraphAxis(panel.Yaxes[1])
	}

	for _, target := range panel.GraphPanel.Targets {
		graphTarget := converter.convertTarget(target)
		if graphTarget == nil {
			continue
		}

		graph.Targets = append(graph.Targets, *graphTarget)
	}

	return grabana.DashboardPanel{Graph: graph}
}

func (converter *JSON) convertGraphVisualization(panel sdk.Panel) *grabana.GraphVisualization {
	graphViz := &grabana.GraphVisualization{
		NullValue: panel.GraphPanel.NullPointMode,
		Staircase: panel.GraphPanel.SteppedLine,
		Overrides: converter.convertGraphOverrides(panel),
	}

	return graphViz
}

func (converter *JSON) convertGraphOverrides(panel sdk.Panel) []grabana.GraphSeriesOverride {
	if len(panel.GraphPanel.SeriesOverrides) == 0 {
		return nil
	}

	overrides := make([]grabana.GraphSeriesOverride, 0, len(panel.GraphPanel.SeriesOverrides))

	for _, sdkOverride := range panel.GraphPanel.SeriesOverrides {
		color := ""
		if sdkOverride.Color != nil {
			color = *sdkOverride.Color
		}

		overrides = append(overrides, grabana.GraphSeriesOverride{
			Alias:     sdkOverride.Alias,
			Color:     color,
			Dashes:    sdkOverride.Dashes,
			Lines:     sdkOverride.Lines,
			Fill:      sdkOverride.Fill,
			LineWidth: sdkOverride.LineWidth,
		})
	}

	return overrides
}

func (converter *JSON) convertGraphLegend(sdkLegend sdk.Legend) []string {
	var legend []string

	if !sdkLegend.Show {
		legend = append(legend, "hide")
	}
	if sdkLegend.AlignAsTable {
		legend = append(legend, "as_table")
	}
	if sdkLegend.RightSide {
		legend = append(legend, "to_the_right")
	}
	if sdkLegend.Min {
		legend = append(legend, "min")
	}
	if sdkLegend.Max {
		legend = append(legend, "max")
	}
	if sdkLegend.Avg {
		legend = append(legend, "avg")
	}
	if sdkLegend.Current {
		legend = append(legend, "current")
	}
	if sdkLegend.Total {
		legend = append(legend, "total")
	}
	if sdkLegend.HideEmpty {
		legend = append(legend, "no_null_series")
	}
	if sdkLegend.HideZero {
		legend = append(legend, "no_zero_series")
	}

	return legend
}

func (converter *JSON) convertGraphAxis(sdkAxis sdk.Axis) *grabana.GraphAxis {
	hidden := !sdkAxis.Show
	var min *float64
	var max *float64

	if sdkAxis.Min != nil {
		min = &sdkAxis.Min.Value
	}
	if sdkAxis.Max != nil {
		max = &sdkAxis.Max.Value
	}

	return &grabana.GraphAxis{
		Hidden:  &hidden,
		Label:   sdkAxis.Label,
		Unit:    &sdkAxis.Format,
		Min:     min,
		Max:     max,
		LogBase: sdkAxis.LogBase,
	}
}
