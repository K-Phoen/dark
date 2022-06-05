package converter

import (
	"strconv"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (converter *JSON) convertHeatmap(panel sdk.Panel) grabana.DashboardPanel {
	heatmap := &grabana.DashboardHeatmap{
		Title:           panel.Title,
		Span:            panelSpan(panel),
		Transparent:     panel.Transparent,
		HideZeroBuckets: panel.HeatmapPanel.HideZeroBuckets,
		HighlightCards:  panel.HeatmapPanel.HighlightCards,
		ReverseYBuckets: panel.HeatmapPanel.ReverseYBuckets,
		Tooltip: &grabana.HeatmapTooltip{
			Show:          panel.HeatmapPanel.Tooltip.Show,
			ShowHistogram: panel.HeatmapPanel.Tooltip.ShowHistogram,
			Decimals:      &panel.HeatmapPanel.TooltipDecimals,
		},
		YAxis: converter.convertHeatmapYAxis(panel),
	}

	if panel.Description != nil {
		heatmap.Description = *panel.Description
	}
	if panel.Repeat != nil {
		heatmap.Repeat = *panel.Repeat
	}
	if panel.Height != nil {
		heatmap.Height = *(panel.Height).(*string)
	}
	if panel.Datasource != nil {
		heatmap.Datasource = panel.Datasource.LegacyName
	}
	if len(panel.Links) != 0 {
		heatmap.Links = converter.convertPanelLinks(panel.Links)
	}
	if panel.HeatmapPanel.DataFormat != "" {
		switch panel.HeatmapPanel.DataFormat {
		case "tsbuckets":
			heatmap.DataFormat = "time_series_buckets"
		case "time_series":
			heatmap.DataFormat = "time_series"
		default:
			converter.logger.Warn("unknown data format: skipping heatmap", zap.String("data_format", panel.HeatmapPanel.DataFormat), zap.String("heatmap_title", panel.Title))
		}
	}

	for _, target := range panel.HeatmapPanel.Targets {
		heatmapTarget := converter.convertTarget(target)
		if heatmapTarget == nil {
			continue
		}

		heatmap.Targets = append(heatmap.Targets, *heatmapTarget)
	}

	return grabana.DashboardPanel{Heatmap: heatmap}
}

func (converter *JSON) convertHeatmapYAxis(panel sdk.Panel) *grabana.HeatmapYAxis {
	panelAxis := panel.HeatmapPanel.YAxis

	axis := &grabana.HeatmapYAxis{
		Decimals: panelAxis.Decimals,
		Unit:     panelAxis.Format,
	}

	if panelAxis.Max != nil {
		max, err := strconv.ParseFloat(*panelAxis.Max, 64)
		if err != nil {
			converter.logger.Warn("could not parse max value on heatmap Y axis %s: %s", zap.String("value", *panelAxis.Max), zap.Error(err))
		} else {
			axis.Max = &max
		}
	}

	if panelAxis.Min != nil {
		min, err := strconv.ParseFloat(*panelAxis.Min, 64)
		if err != nil {
			converter.logger.Warn("could not parse min value on heatmap Y axis %s: %s", zap.String("value", *panelAxis.Min), zap.Error(err))
		} else {
			axis.Min = &min
		}
	}

	return axis
}
