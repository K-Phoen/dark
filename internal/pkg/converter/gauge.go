package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (converter *JSON) convertGauge(panel sdk.Panel) grabana.DashboardPanel {
	gauge := &grabana.DashboardGauge{
		Title:         panel.Title,
		Span:          panelSpan(panel),
		Unit:          panel.GaugePanel.FieldConfig.Defaults.Unit,
		Decimals:      panel.GaugePanel.FieldConfig.Defaults.Decimals,
		Transparent:   panel.Transparent,
		Orientation:   converter.convertGaugeOrientation(panel),
		ValueType:     converter.convertGaugeValueType(panel),
		ThresholdMode: converter.convertGaugeThresholdMode(panel),
		Thresholds:    converter.convertGaugeThresholds(panel),
	}

	if panel.GaugePanel.Options.Text != nil {
		gauge.ValueFontSize = panel.GaugePanel.Options.Text.ValueSize
		gauge.TitleFontSize = panel.GaugePanel.Options.Text.TitleSize
	}

	if panel.Description != nil {
		gauge.Description = *panel.Description
	}
	if panel.Repeat != nil {
		gauge.Repeat = *panel.Repeat
	}
	if panel.RepeatDirection != nil {
		gauge.RepeatDirection = sdkRepeatDirectionToYAML(*panel.RepeatDirection)
	}
	if panel.Height != nil {
		gauge.Height = *(panel.Height).(*string)
	}
	if panel.Datasource != nil {
		gauge.Datasource = panel.Datasource.LegacyName
	}
	if len(panel.Links) != 0 {
		gauge.Links = converter.convertPanelLinks(panel.Links)
	}

	for _, target := range panel.GaugePanel.Targets {
		graphTarget := converter.convertTarget(target)
		if graphTarget == nil {
			continue
		}

		gauge.Targets = append(gauge.Targets, *graphTarget)
	}

	return grabana.DashboardPanel{Gauge: gauge}
}

func (converter *JSON) convertGaugeValueType(panel sdk.Panel) string {
	if len(panel.GaugePanel.Options.ReduceOptions.Calcs) != 1 {
		return "last_non_null"
	}

	valueType := panel.GaugePanel.Options.ReduceOptions.Calcs[0]

	switch valueType {
	case "first":
		return "first"
	case "firstNotNull":
		return "first_non_null"
	case "last":
		return "last"
	case "lastNotNull":
		return "last_non_null"

	case "min":
		return "min"
	case "max":
		return "max"
	case "mean":
		return "avg"

	case "count":
		return "count"
	case "sum":
		return "total"
	case "range":
		return "range"

	default:
		converter.logger.Warn("unknown value type", zap.String("value type", valueType))
		return "last_non_null"
	}
}

func (converter *JSON) convertGaugeOrientation(panel sdk.Panel) string {
	switch panel.GaugePanel.Options.Orientation {
	case "", "auto":
		return "auto"
	case "horizontal":
		return "horizontal"
	case "vertical":
		return "vertical"
	default:
		converter.logger.Warn("unknown orientation", zap.String("orientation", panel.GaugePanel.Options.Orientation))
		return "auto"
	}
}

func (converter *JSON) convertGaugeThresholdMode(panel sdk.Panel) string {
	switch panel.GaugePanel.FieldConfig.Defaults.Thresholds.Mode {
	case "":
		return "absolute"
	case "absolute":
		return "absolute"
	case "percentage":
		return "relative"
	default:
		converter.logger.Warn("unknown threshold mode", zap.String("mode", panel.GaugePanel.FieldConfig.Defaults.Thresholds.Mode))
		return "absolute"
	}
}

func (converter *JSON) convertGaugeThresholds(panel sdk.Panel) []grabana.GaugeThresholdStep {
	steps := make([]grabana.GaugeThresholdStep, 0, len(panel.GaugePanel.FieldConfig.Defaults.Thresholds.Steps))

	for _, step := range panel.GaugePanel.FieldConfig.Defaults.Thresholds.Steps {
		steps = append(steps, grabana.GaugeThresholdStep{
			Color: step.Color,
			Value: step.Value,
		})
	}

	return steps
}
