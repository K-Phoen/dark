package converter

import (
	"fmt"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (converter *JSON) convertTimeSeries(panel sdk.Panel) grabana.DashboardPanel {
	tsPanel := &grabana.DashboardTimeSeries{
		Title:         panel.Title,
		Span:          panelSpan(panel),
		Transparent:   panel.Transparent,
		Legend:        converter.convertTimeSeriesLegend(panel.TimeseriesPanel.Options.Legend),
		Visualization: converter.convertTimeSeriesVisualization(panel),
		Axis:          converter.convertTimeSeriesAxis(panel),
		Overrides:     converter.convertTimeSeriesOverrides(panel),
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
		tsPanel.Datasource = panel.Datasource.LegacyName
	}
	if len(panel.Links) != 0 {
		tsPanel.Links = converter.convertPanelLinks(panel.Links)
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
		LineWidth:   &panel.TimeseriesPanel.FieldConfig.Defaults.Custom.LineWidth,
	}

	if panel.TimeseriesPanel.FieldConfig.Defaults.Custom.DrawStyle == "line" {
		tsViz.LineInterpolation = converter.convertTimeSeriesLineInterpolation(panel)
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

	// Stacking mode
	switch panel.TimeseriesPanel.FieldConfig.Defaults.Custom.Stacking.Mode {
	case "none":
		tsViz.Stack = "none"
	case "normal":
		tsViz.Stack = "normal"
	case "percent":
		tsViz.Stack = "percent"
	default:
		tsViz.Stack = "none"
	}

	return tsViz
}

func (converter *JSON) convertTimeSeriesLineInterpolation(panel sdk.Panel) string {
	mode := panel.TimeseriesPanel.FieldConfig.Defaults.Custom.LineInterpolation

	switch mode {
	case "smooth":
		return "smooth"
	case "linear":
		return "linear"
	case "stepBefore":
		return "step_before"
	case "stepAfter":
		return "step_after"
	default:
		converter.logger.Warn("invalid line interpolation mode, defaulting to smooth", zap.String("interpolation_mode", mode))
		return "smooth"
	}
}

func (converter *JSON) convertTimeSeriesLegend(legend sdk.TimeseriesLegendOptions) []string {
	options := []string{}

	// Hidden legend?
	if legend.Show != nil && !*legend.Show {
		options = append(options, "hide")
	} else {
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

func (converter *JSON) convertTimeSeriesOverrides(panel sdk.Panel) []grabana.TimeSeriesOverride {
	overrides := make([]grabana.TimeSeriesOverride, 0, len(panel.TimeseriesPanel.FieldConfig.Overrides))

	for _, sdkOverride := range panel.TimeseriesPanel.FieldConfig.Overrides {
		override, err := converter.convertTimeSeriesOverride(sdkOverride)
		if err != nil {
			converter.logger.Warn("could not convert field override: skipping", zap.Error(err))
			continue
		}

		overrides = append(overrides, override)
	}

	return overrides
}

func (converter *JSON) convertTimeSeriesOverride(sdkOverride sdk.FieldConfigOverride) (grabana.TimeSeriesOverride, error) {
	override := grabana.TimeSeriesOverride{}

	matcher, err := converter.convertTimeSeriesOverrideMatcher(sdkOverride.Matcher)
	if err != nil {
		return override, err
	}
	override.Matcher = matcher

	properties, err := converter.convertTimeSeriesOverrideProperties(sdkOverride.Properties)
	if err != nil {
		return override, err
	}
	override.Properties = properties

	return override, nil
}

func (converter *JSON) convertTimeSeriesOverrideMatcher(matcher struct {
	ID      string `json:"id"`
	Options string `json:"options"`
}) (grabana.TimeSeriesOverrideMatcher, error) {
	switch matcher.ID {
	case "byName":
		return grabana.TimeSeriesOverrideMatcher{FieldName: &matcher.Options}, nil
	case "byFrameRefID":
		return grabana.TimeSeriesOverrideMatcher{QueryRef: &matcher.Options}, nil
	case "byRegexp":
		return grabana.TimeSeriesOverrideMatcher{Regex: &matcher.Options}, nil
	case "byType":
		return grabana.TimeSeriesOverrideMatcher{Type: &matcher.Options}, nil
	default:
		return grabana.TimeSeriesOverrideMatcher{}, fmt.Errorf("unknown field override matcher '%s'", matcher.ID)
	}
}

func (converter *JSON) convertTimeSeriesOverrideProperties(sdkProperties []sdk.FieldConfigOverrideProperty) (grabana.TimeSeriesOverrideProperties, error) {
	properties := grabana.TimeSeriesOverrideProperties{}

	for _, sdkProperty := range sdkProperties {
		converter.convertTimeSeriesOverrideProperty(sdkProperty, &properties)
	}

	return properties, nil
}

func (converter *JSON) convertTimeSeriesOverrideProperty(sdkProperty sdk.FieldConfigOverrideProperty, properties *grabana.TimeSeriesOverrideProperties) {
	switch sdkProperty.ID {
	case "unit":
		properties.Unit = strPtr(sdkProperty.Value.(string))
	case "custom.axisPlacement":
		properties.AxisDisplay = strPtr(sdkProperty.Value.(string))
	case "custom.fillOpacity":
		properties.FillOpacity = intPtr(int(sdkProperty.Value.(float64)))
	case "custom.stacking":
		options, ok := sdkProperty.Value.(map[string]interface{})
		if !ok {
			converter.logger.Warn("could not convert custom.stacking field override: invalid options")
			break
		}
		properties.Stack = strPtr(options["mode"].(string))
	case "custom.transform":
		transformType := sdkProperty.Value.(string)
		if transformType != "negative-Y" {
			converter.logger.Warn("could not convert transform field override: invalid option")
			break
		}
		properties.NegativeY = boolPtr(true)
	case "color":
		options, ok := sdkProperty.Value.(map[string]interface{})
		if !ok {
			converter.logger.Warn("could not convert color field override: invalid options")
			break
		}
		if options["mode"] != "fixed" {
			converter.logger.Warn("could not convert color field override: unsupported mode")
			break
		}

		properties.Color = strPtr(options["fixedColor"].(string))
	default:
		converter.logger.Warn(fmt.Sprintf("unhandled override type '%s'", sdkProperty.ID))
	}
}
