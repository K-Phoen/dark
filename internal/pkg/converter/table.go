package converter

import (
	grabana "github.com/K-Phoen/grabana/decoder"
	grabanaTable "github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (converter *JSON) convertTable(panel sdk.Panel) grabana.DashboardPanel {
	table := &grabana.DashboardTable{
		Title:       panel.Title,
		Span:        panelSpan(panel),
		Transparent: panel.Transparent,
	}

	if panel.Description != nil {
		table.Description = *panel.Description
	}
	if panel.Height != nil {
		table.Height = *(panel.Height).(*string)
	}
	if panel.Datasource != nil {
		table.Datasource = panel.Datasource.LegacyName
	}
	if len(panel.Links) != 0 {
		table.Links = converter.convertPanelLinks(panel.Links)
	}

	for _, target := range panel.TablePanel.Targets {
		graphTarget := converter.convertTarget(target)
		if graphTarget == nil {
			continue
		}

		table.Targets = append(table.Targets, *graphTarget)
	}

	// hidden columns
	for _, columnStyle := range panel.TablePanel.Styles {
		if columnStyle.Type != "hidden" {
			continue
		}

		table.HiddenColumns = append(table.HiddenColumns, columnStyle.Pattern)
	}

	// time series aggregations
	if panel.TablePanel.Transform == "timeseries_aggregations" {
		for _, column := range panel.TablePanel.Columns {
			table.TimeSeriesAggregations = append(table.TimeSeriesAggregations, grabanaTable.Aggregation{
				Label: column.TextType,
				Type:  grabanaTable.AggregationType(column.Value),
			})
		}
	} else {
		converter.logger.Warn("unhandled transform type: skipped", zap.String("transform", panel.TablePanel.Transform), zap.String("panel", panel.Title))
	}

	return grabana.DashboardPanel{Table: table}
}
