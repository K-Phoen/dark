package converter

import (
	"strings"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/grabana/singlestat"
	"github.com/K-Phoen/sdk"
)

func (converter *JSON) convertSingleStat(panel sdk.Panel) grabana.DashboardPanel {
	singleStat := &grabana.DashboardSingleStat{
		Title:         panel.Title,
		Span:          panelSpan(panel),
		Unit:          panel.SinglestatPanel.Format,
		Decimals:      &panel.SinglestatPanel.Decimals,
		ValueType:     panel.SinglestatPanel.ValueName,
		Transparent:   panel.Transparent,
		ValueFontSize: panel.SinglestatPanel.ValueFontSize,
	}

	if panel.Description != nil {
		singleStat.Description = *panel.Description
	}
	if panel.Repeat != nil {
		singleStat.Repeat = *panel.Repeat
	}
	if panel.RepeatDirection != nil {
		singleStat.RepeatDirection = sdkRepeatDirectionToYAML(*panel.RepeatDirection)
	}
	if panel.Height != nil {
		singleStat.Height = *(panel.Height).(*string)
	}
	if panel.Datasource != nil {
		singleStat.Datasource = panel.Datasource.LegacyName
	}
	if len(panel.Links) != 0 {
		singleStat.Links = converter.convertPanelLinks(panel.Links)
	}

	thresholds := strings.Split(panel.SinglestatPanel.Thresholds, ",")
	if len(thresholds) == 2 {
		singleStat.Thresholds = [2]string{thresholds[0], thresholds[1]}
	}

	if len(panel.SinglestatPanel.Colors) == 3 {
		singleStat.Colors = [3]string{
			panel.SinglestatPanel.Colors[0],
			panel.SinglestatPanel.Colors[1],
			panel.SinglestatPanel.Colors[2],
		}
	}

	var colorOpts []string
	if panel.SinglestatPanel.ColorBackground {
		colorOpts = append(colorOpts, "background")
	}
	if panel.SinglestatPanel.ColorValue {
		colorOpts = append(colorOpts, "value")
	}
	if len(colorOpts) != 0 {
		singleStat.Color = colorOpts
	}

	if panel.SinglestatPanel.SparkLine.Show && panel.SinglestatPanel.SparkLine.Full {
		singleStat.SparkLine = "full"
	}
	if panel.SinglestatPanel.SparkLine.Show && !panel.SinglestatPanel.SparkLine.Full {
		singleStat.SparkLine = "bottom"
	}

	// Font sizes
	if panel.SinglestatPanel.PrefixFontSize != nil && *panel.SinglestatPanel.PrefixFontSize != "" {
		singleStat.PrefixFontSize = *panel.SinglestatPanel.PrefixFontSize
	}
	if panel.SinglestatPanel.PostfixFontSize != nil && *panel.SinglestatPanel.PostfixFontSize != "" {
		singleStat.PostfixFontSize = *panel.SinglestatPanel.PostfixFontSize
	}

	// ranges to text mapping
	singleStat.RangesToText = converter.convertSingleStatRangesToText(panel)

	for _, target := range panel.SinglestatPanel.Targets {
		graphTarget := converter.convertTarget(target)
		if graphTarget == nil {
			continue
		}

		singleStat.Targets = append(singleStat.Targets, *graphTarget)
	}

	return grabana.DashboardPanel{SingleStat: singleStat}
}

func (converter *JSON) convertSingleStatRangesToText(panel sdk.Panel) []singlestat.RangeMap {
	if panel.SinglestatPanel.MappingType == nil || *panel.SinglestatPanel.MappingType != 2 {
		return nil
	}

	mappings := make([]singlestat.RangeMap, 0, len(panel.SinglestatPanel.RangeMaps))
	for _, mapping := range panel.SinglestatPanel.RangeMaps {
		converted := singlestat.RangeMap{
			From: "",
			To:   "",
			Text: "",
		}

		if mapping.From != nil {
			converted.From = *mapping.From
		}
		if mapping.To != nil {
			converted.To = *mapping.To
		}
		if mapping.Text != nil {
			converted.Text = *mapping.Text
		}

		mappings = append(mappings, converted)
	}

	return mappings
}
