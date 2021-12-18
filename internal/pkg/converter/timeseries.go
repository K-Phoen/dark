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
		Axis:          converter.convertTimeSeriesAxis(panel),
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

func (converter *JSON) convertTimeSeriesAxis(panel sdk.Panel) *grabana.TimeSeriesAxis {
	fieldConfig := panel.TimeseriesPanel.FieldConfig

	tsAxis := &grabana.TimeSeriesAxis{
		Unit:  fieldConfig.Defaults.Unit,
		Label: fieldConfig.Defaults.Custom.AxisLabel,
	}

	// decimals
	if fieldConfig.Defaults.Decimals != nil {
		tsAxis.Decimals = fieldConfig.Defaults.Decimals
	}

	// boundaries
	if fieldConfig.Defaults.Min != nil {
		tsAxis.Min = fieldConfig.Defaults.Min
	}
	if fieldConfig.Defaults.Max != nil {
		tsAxis.Max = fieldConfig.Defaults.Max
	}
	if fieldConfig.Defaults.Custom.AxisSoftMin != nil {
		tsAxis.SoftMin = fieldConfig.Defaults.Custom.AxisSoftMin
	}
	if fieldConfig.Defaults.Custom.AxisSoftMax != nil {
		tsAxis.SoftMax = fieldConfig.Defaults.Custom.AxisSoftMax
	}

	// placement
	switch fieldConfig.Defaults.Custom.AxisPlacement {
	case "hidden":
		tsAxis.Display = "hidden"
	case "left":
		tsAxis.Display = "left"
	case "right":
		tsAxis.Display = "right"
	case "auto":
		tsAxis.Display = "auto"
	}

	// scale
	switch fieldConfig.Defaults.Custom.ScaleDistribution.Type {
	case "linear":
		tsAxis.Scale = "linear"
	case "log":
		if fieldConfig.Defaults.Custom.ScaleDistribution.Log == 2 {
			tsAxis.Scale = "log2"
		} else {
			tsAxis.Scale = "log10"
		}
	}

	return tsAxis
}

func (converter *JSON) convertTimeSeriesVisualization(panel sdk.Panel) *grabana.TimeSeriesVisualization {
	tsViz := &grabana.TimeSeriesVisualization{
		FillOpacity: &panel.TimeseriesPanel.FieldConfig.Defaults.Custom.FillOpacity,
		PointSize:   &panel.TimeseriesPanel.FieldConfig.Defaults.Custom.PointSize,
	}

	// Tooltip mode
	switch panel.TimeseriesPanel.Options.Tooltip.Mode {
	case "none":
		tsViz.Tooltip = "none"
	case "multi":
		tsViz.Tooltip = "all_series"
	default:
		tsViz.Tooltip = "single_series"
	}

	// Gradient mode
	switch panel.TimeseriesPanel.FieldConfig.Defaults.Custom.GradientMode {
	case "none":
		tsViz.GradientMode = "none"
	case "hue":
		tsViz.GradientMode = "hue"
	case "scheme":
		tsViz.GradientMode = "scheme"
	default:
		tsViz.GradientMode = "opacity"
	}

	return tsViz
}

func (converter *JSON) convertTimeSeriesLegend(legend sdk.TimeseriesLegendOptions) []string {
	options := []string{}

	// Display mode
	switch legend.DisplayMode {
	case "list":
		options = append(options, "as_list")
	case "hidden":
		options = append(options, "hide")
	default:
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
