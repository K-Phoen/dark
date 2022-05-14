package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (converter *JSON) convertStat(panel sdk.Panel) grabana.DashboardPanel {
	stat := &grabana.DashboardStat{
		Title:         panel.Title,
		Span:          panelSpan(panel),
		Unit:          panel.StatPanel.FieldConfig.Defaults.Unit,
		Decimals:      panel.StatPanel.FieldConfig.Defaults.Decimals,
		Transparent:   panel.Transparent,
		ValueFontSize: panel.StatPanel.Options.Text.ValueSize,
		TitleFontSize: panel.StatPanel.Options.Text.TitleSize,
		Orientation:   converter.convertStatOrientation(panel),
		Text:          converter.convertStatTextMode(panel),
		ValueType:     converter.convertStatValueType(panel),
		ColorMode:     converter.convertStatColorMode(panel),
		ThresholdMode: converter.convertStatThresholdMode(panel),
		Thresholds:    converter.convertStatThresholds(panel),
	}

	if panel.Description != nil {
		stat.Description = *panel.Description
	}
	if panel.Repeat != nil {
		stat.Repeat = *panel.Repeat
	}
	if panel.Height != nil {
		stat.Height = *(panel.Height).(*string)
	}
	if panel.Datasource != nil {
		stat.Datasource = *panel.Datasource
	}
	if panel.StatPanel.Options.GraphMode == "area" {
		stat.SparkLine = true
	}

	for _, target := range panel.StatPanel.Targets {
		graphTarget := converter.convertTarget(target)
		if graphTarget == nil {
			continue
		}

		stat.Targets = append(stat.Targets, *graphTarget)
	}

	return grabana.DashboardPanel{Stat: stat}
}

func (converter *JSON) convertStatOrientation(panel sdk.Panel) string {
	switch panel.StatPanel.Options.Orientation {
	case "":
		return "auto"
	case "horizontal":
		return "horizontal"
	case "vertical":
		return "vertical"
	default:
		converter.logger.Warn("unknown orientation", zap.String("orientation", panel.StatPanel.Options.Orientation))
		return "auto"
	}
}

func (converter *JSON) convertStatTextMode(panel sdk.Panel) string {
	switch panel.StatPanel.Options.TextMode {
	case "auto":
		return "auto"
	case "value":
		return "value"
	case "name":
		return "name"
	case "value_and_name":
		return "value_and_name"
	case "none":
		return "none"
	default:
		converter.logger.Warn("unknown text mode", zap.String("mode", panel.StatPanel.Options.TextMode))
		return "auto"
	}
}

func (converter *JSON) convertStatValueType(panel sdk.Panel) string {
	if len(panel.StatPanel.Options.ReduceOptions.Calcs) != 1 {
		return "last_non_null"
	}

	valueType := panel.StatPanel.Options.ReduceOptions.Calcs[0]

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

func (converter *JSON) convertStatColorMode(panel sdk.Panel) string {
	switch panel.StatPanel.Options.ColorMode {
	case "background":
		return "background"
	case "value":
		return "value"
	default:
		converter.logger.Warn("unknown color mode", zap.String("color_mode", panel.StatPanel.Options.ColorMode))
		return "value"
	}
}

func (converter *JSON) convertStatThresholdMode(panel sdk.Panel) string {
	switch panel.StatPanel.FieldConfig.Defaults.Thresholds.Mode {
	case "":
		return "absolute"
	case "absolute":
		return "absolute"
	case "percentage":
		return "relative"
	default:
		converter.logger.Warn("unknown threshold mode", zap.String("mode", panel.StatPanel.FieldConfig.Defaults.Thresholds.Mode))
		return "absolute"
	}
}

func (converter *JSON) convertStatThresholds(panel sdk.Panel) []grabana.StatThresholdStep {
	steps := make([]grabana.StatThresholdStep, 0, len(panel.StatPanel.FieldConfig.Defaults.Thresholds.Steps))

	for _, step := range panel.StatPanel.FieldConfig.Defaults.Thresholds.Steps {
		steps = append(steps, grabana.StatThresholdStep{
			Color: step.Color,
			Value: step.Value,
		})
	}

	return steps
}
