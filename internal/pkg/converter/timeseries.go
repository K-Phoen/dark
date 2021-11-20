package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
)

func (converter *JSON) convertTimeSeries(panel sdk.Panel) grabana.DashboardPanel {
	tsPanel := &grabana.DashboardTimeSeries{
		Title:         panel.Title,
		Span:          panelSpan(panel),
		Transparent:   panel.Transparent,
		Alert:         converter.convertAlert(panel),
		Legend:        converter.convertTimeSeriesLegend(panel.TimeseriesPanel.Options.Legend),
		Visualization: converter.convertTimeSeriesVisualization(panel),
		Axis:          nil,
	}

	if panel.Description != nil {
		tsPanel.Description = *panel.Description
	}
	if panel.Repeat != nil {
		tsPanel.Repeat = *panel.Repeat
	}
	if panel.Height != nil {
		tsPanel.Height = panel.Height.(string)
	}
	if panel.Datasource != nil {
		tsPanel.Datasource = *panel.Datasource
	}

	for _, target := range panel.TimeseriesPanel.Targets {
		tsTarget := converter.convertTarget(target)
		if tsTarget == nil {
			continue
		}

		tsPanel.Targets = append(tsPanel.Targets, *tsTarget)
	}

	return grabana.DashboardPanel{TimeSeries: tsPanel}
}

func (converter *JSON) convertTimeSeriesVisualization(panel sdk.Panel) *grabana.TimeSeriesVisualization {
	tsViz := &grabana.TimeSeriesVisualization{
		FillOpacity: &panel.TimeseriesPanel.FieldConfig.Defaults.Custom.FillOpacity,
		PointSize:   &panel.TimeseriesPanel.FieldConfig.Defaults.Custom.PointSize,
	}

	// Tooltip mode
	if panel.TimeseriesPanel.Options.Tooltip.Mode == "none" {
		tsViz.Tooltip = "none"
	} else if panel.TimeseriesPanel.Options.Tooltip.Mode == "multi" {
		tsViz.Tooltip = "all_series"
	} else {
		tsViz.Tooltip = "single_series"
	}

	// Gradient mode
	if panel.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode == "none" {
		tsViz.GradientMode = "none"
	} else if panel.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode == "hue" {
		tsViz.GradientMode = "hue"
	} else if panel.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode == "scheme" {
		tsViz.GradientMode = "scheme"
	} else {
		tsViz.GradientMode = "opacity"
	}

	return tsViz
}

func (converter *JSON) convertTimeSeriesLegend(legend sdk.TimeseriesLegendOptions) []string {
	options := []string{}

	// Display mode
	if legend.DisplayMode == "list" {
		options = append(options, "as_list")
	} else if legend.DisplayMode == "hidden" {
		options = append(options, "hide")
	} else {
		options = append(options, "as_table")
	}

	// Placement
	if legend.Placement == "right" {
		options = append(options, "to_the_right")
	} else {
		options = append(options, "to_bottom")
	}

	// Automatic calculations
	calcs := map[string]string{
		"first":        "first",
		"firstNotNull": "first_non_null",
		"last":         "last",
		"lastNotNull":  "last_non_null",

		"min":  "min",
		"max":  "max",
		"mean": "avg",

		"count": "count",
		"sum":   "total",
		"range": "range",
	}

	for sdkCalc, grabanaCalc := range calcs {
		if !stringInSlice(sdkCalc, legend.Calcs) {
			continue
		}

		options = append(options, grabanaCalc)
	}

	return options
}
