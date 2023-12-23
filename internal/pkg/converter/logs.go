package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/grabana/logs"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (converter *JSON) convertLogs(panel sdk.Panel) grabana.DashboardPanel {
	convertedLogs := &grabana.DashboardLogs{
		Title:         panel.Title,
		Span:          panelSpan(panel),
		Transparent:   panel.Transparent,
		Visualization: converter.convertLogsVizualization(panel),
	}

	if panel.Description != nil {
		convertedLogs.Description = *panel.Description
	}
	if panel.Repeat != nil {
		convertedLogs.Repeat = *panel.Repeat
	}
	if panel.RepeatDirection != nil {
		convertedLogs.RepeatDirection = sdkRepeatDirectionToYAML(*panel.RepeatDirection)
	}
	if panel.Height != nil {
		convertedLogs.Height = *(panel.Height).(*string)
	}
	if panel.Datasource != nil {
		convertedLogs.Datasource = panel.Datasource.LegacyName
	}
	if len(panel.Links) != 0 {
		convertedLogs.Links = converter.convertPanelLinks(panel.Links)
	}

	for _, target := range panel.LogsPanel.Targets {
		logsTarget := converter.convertLogsTarget(target)

		convertedLogs.Targets = append(convertedLogs.Targets, logsTarget)
	}

	return grabana.DashboardPanel{Logs: convertedLogs}
}

func (converter *JSON) convertLogsVizualization(panel sdk.Panel) *grabana.LogsVisualization {
	return &grabana.LogsVisualization{
		Time:           panel.LogsPanel.Options.ShowTime,
		UniqueLabels:   panel.LogsPanel.Options.ShowLabels,
		CommonLabels:   panel.LogsPanel.Options.ShowCommonLabels,
		WrapLines:      panel.LogsPanel.Options.WrapLogMessage,
		PrettifyJSON:   panel.LogsPanel.Options.PrettifyLogMessage,
		HideLogDetails: !panel.LogsPanel.Options.EnableLogDetails,
		Order:          converter.convertLogsSortOrder(panel.LogsPanel.Options.SortOrder),
		Deduplication:  converter.convertLogsDedupStrategy(panel.LogsPanel.Options.DedupStrategy),
	}
}

func (converter *JSON) convertLogsDedupStrategy(strategy string) string {
	switch strategy {
	case string(logs.None):
		return "none"
	case string(logs.Exact):
		return "exact"
	case string(logs.Numbers):
		return "numbers"
	case string(logs.Signature):
		return "signature"
	default:
		converter.logger.Warn("unhandled logs dedup strategy: skipped", zap.String("strategy", strategy))
		return ""
	}
}

func (converter *JSON) convertLogsSortOrder(order string) string {
	switch order {
	case string(logs.Asc):
		return "asc"
	case string(logs.Desc):
		return "desc"
	default:
		converter.logger.Warn("unhandled sort order: skipped", zap.String("order", order))
		return ""
	}
}

func (converter *JSON) convertLogsTarget(target sdk.Target) grabana.LogsTarget {
	return grabana.LogsTarget{
		Loki: &grabana.LokiTarget{
			Query:  target.Expr,
			Legend: target.LegendFormat,
			Ref:    target.RefID,
			Hidden: target.Hide,
		},
	}
}
