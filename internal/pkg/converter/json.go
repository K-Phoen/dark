package converter

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	grabanaDashboard "github.com/K-Phoen/grabana/dashboard"
	grabana "github.com/K-Phoen/grabana/decoder"
	grabanaTable "github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/grafana-tools/sdk"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type JSON struct {
	logger *zap.Logger
}

func NewJSON(logger *zap.Logger) *JSON {
	return &JSON{
		logger: logger,
	}
}

func (converter *JSON) Convert(input io.Reader, output io.Writer) error {
	content, err := ioutil.ReadAll(input)
	if err != nil {
		converter.logger.Error("could not read input", zap.Error(err))
		return err
	}

	board := &sdk.Board{}
	if err := json.Unmarshal(content, board); err != nil {
		converter.logger.Error("could not unmarshall dashboard", zap.Error(err))
		return err
	}

	dashboard := &grabana.DashboardModel{}

	converter.convertGeneralSettings(board, dashboard)
	converter.convertVariables(board.Templating.List, dashboard)
	converter.convertAnnotations(board.Annotations.List, dashboard)
	converter.convertPanels(board.Panels, dashboard)

	converted, err := yaml.Marshal(dashboard)
	if err != nil {
		converter.logger.Error("could marshall dashboard to yaml", zap.Error(err))
		return err
	}

	_, err = output.Write(converted)

	return err
}

func (converter *JSON) convertGeneralSettings(board *sdk.Board, dashboard *grabana.DashboardModel) {
	dashboard.Title = board.Title
	dashboard.SharedCrosshair = board.SharedCrosshair
	dashboard.Tags = board.Tags
	dashboard.Editable = board.Editable
	dashboard.Time = [2]string{board.Time.From, board.Time.To}
	dashboard.Timezone = board.Timezone

	if board.Refresh != nil {
		dashboard.AutoRefresh = board.Refresh.Value
	}
}

func (converter *JSON) convertAnnotations(annotations []sdk.Annotation, dashboard *grabana.DashboardModel) {
	for _, annotation := range annotations {
		// grafana-sdk doesn't expose the "builtIn" field, so we work around that by skipping
		// the annotation we know to be built-in by its name
		if annotation.Name == "Annotations & Alerts" {
			continue
		}

		if annotation.Type != "tags" {
			converter.logger.Warn("unhandled annotation type: skipped", zap.String("type", annotation.Type), zap.String("name", annotation.Name))
			continue
		}

		converter.convertTagAnnotation(annotation, dashboard)
	}
}

func (converter *JSON) convertTagAnnotation(annotation sdk.Annotation, dashboard *grabana.DashboardModel) {
	datasource := ""
	if annotation.Datasource != nil {
		datasource = *annotation.Datasource
	}

	dashboard.TagsAnnotation = append(dashboard.TagsAnnotation, grabanaDashboard.TagAnnotation{
		Name:       annotation.Name,
		Datasource: datasource,
		IconColor:  annotation.IconColor,
		Tags:       annotation.Tags,
	})
}

func (converter *JSON) convertVariables(variables []sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	for _, variable := range variables {
		converter.convertVariable(variable, dashboard)
	}
}

func (converter *JSON) convertVariable(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	switch variable.Type {
	case "interval":
		converter.convertIntervalVar(variable, dashboard)
	case "custom":
		converter.convertCustomVar(variable, dashboard)
	case "query":
		converter.convertQueryVar(variable, dashboard)
	case "const":
		converter.convertConstVar(variable, dashboard)
	default:
		converter.logger.Warn("unhandled variable type found: skipped", zap.String("type", variable.Type), zap.String("name", variable.Name))
	}
}

func (converter *JSON) convertIntervalVar(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	interval := &grabana.VariableInterval{
		Name:    variable.Name,
		Label:   variable.Label,
		Default: defaultOption(variable.Current),
		Values:  make([]string, 0, len(variable.Options)),
	}

	for _, opt := range variable.Options {
		interval.Values = append(interval.Values, opt.Value)
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Interval: interval})
}

func (converter *JSON) convertCustomVar(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	custom := &grabana.VariableCustom{
		Name:      variable.Name,
		Label:     variable.Label,
		Default:   defaultOption(variable.Current),
		ValuesMap: make(map[string]string, len(variable.Options)),
	}

	for _, opt := range variable.Options {
		custom.ValuesMap[opt.Text] = opt.Value
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Custom: custom})
}

func (converter *JSON) convertQueryVar(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	datasource := ""
	if variable.Datasource != nil {
		datasource = *variable.Datasource
	}

	query := &grabana.VariableQuery{
		Name:       variable.Name,
		Label:      variable.Label,
		Datasource: datasource,
		Request:    variable.Query,
		IncludeAll: variable.IncludeAll,
		DefaultAll: variable.Current.Value == "$__all",
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Query: query})
}

func (converter *JSON) convertConstVar(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	constant := &grabana.VariableConst{
		Name:      variable.Name,
		Label:     variable.Label,
		Default:   variable.Current.Text,
		ValuesMap: make(map[string]string, len(variable.Options)),
	}

	for _, opt := range variable.Options {
		constant.ValuesMap[opt.Text] = opt.Value
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Const: constant})
}

func (converter *JSON) convertPanels(panels []*sdk.Panel, dashboard *grabana.DashboardModel) {
	var currentRow *grabana.DashboardRow

	for _, panel := range panels {
		if panel.Type == "row" {
			if currentRow != nil {
				dashboard.Rows = append(dashboard.Rows, *currentRow)
			}

			currentRow = converter.convertRow(*panel)
			continue
		}

		if currentRow == nil {
			currentRow = &grabana.DashboardRow{Name: "Overview"}
		}

		switch panel.Type {
		case "graph":
			currentRow.Panels = append(currentRow.Panels, converter.convertGraph(*panel))
		case "singlestat":
			currentRow.Panels = append(currentRow.Panels, converter.convertSingleStat(*panel))
		case "table":
			currentRow.Panels = append(currentRow.Panels, converter.convertTable(*panel))
		default:
			converter.logger.Warn("unhandled panel type: skipped", zap.String("type", panel.Type), zap.String("title", panel.Title))
		}
	}

	if currentRow != nil {
		dashboard.Rows = append(dashboard.Rows, *currentRow)
	}
}

func (converter *JSON) convertRow(panel sdk.Panel) *grabana.DashboardRow {
	return &grabana.DashboardRow{
		Name:   panel.Title,
		Panels: nil,
	}
}

func (converter *JSON) convertGraph(panel sdk.Panel) grabana.DashboardPanel {
	graph := &grabana.DashboardGraph{
		Title: panel.Title,
		Span:  panelSpan(panel),
		Axes: &grabana.GraphAxes{
			Bottom: converter.convertAxis(panel.Xaxis),
		},
		Legend: converter.convertLegend(panel.Legend),
	}

	if panel.Height != nil {
		graph.Height = *panel.Height
	}
	if panel.Datasource != nil {
		graph.Datasource = *panel.Datasource
	}

	if len(panel.Yaxes) == 2 {
		graph.Axes.Left = converter.convertAxis(panel.Yaxes[0])
		graph.Axes.Right = converter.convertAxis(panel.Yaxes[1])
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

func (converter *JSON) convertLegend(sdkLegend sdk.Legend) []string {
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

func (converter *JSON) convertAxis(sdkAxis sdk.Axis) *grabana.GraphAxis {
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

func (converter *JSON) convertSingleStat(panel sdk.Panel) grabana.DashboardPanel {
	singleStat := &grabana.DashboardSingleStat{
		Title:     panel.Title,
		Span:      panelSpan(panel),
		Unit:      panel.SinglestatPanel.Format,
		ValueType: panel.SinglestatPanel.ValueName,
	}

	if panel.Height != nil {
		singleStat.Height = *panel.Height
	}
	if panel.Datasource != nil {
		singleStat.Datasource = *panel.Datasource
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

	for _, target := range panel.SinglestatPanel.Targets {
		graphTarget := converter.convertTarget(target)
		if graphTarget == nil {
			continue
		}

		singleStat.Targets = append(singleStat.Targets, *graphTarget)
	}

	return grabana.DashboardPanel{SingleStat: singleStat}
}

func (converter *JSON) convertTable(panel sdk.Panel) grabana.DashboardPanel {
	table := &grabana.DashboardTable{
		Title: panel.Title,
		Span:  panelSpan(panel),
	}

	if panel.Height != nil {
		table.Height = *panel.Height
	}
	if panel.Datasource != nil {
		table.Datasource = *panel.Datasource
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

func (converter *JSON) convertTarget(target sdk.Target) *grabana.Target {
	// looks like a prometheus target
	if target.Expr != "" {
		return converter.convertPrometheusTarget(target)
	}

	if target.MetricType != "" {
		return converter.convertStackdriverTarget(target)
	}

	converter.logger.Warn("unhandled target type: skipped", zap.Any("target", target))

	return nil
}

func (converter *JSON) convertPrometheusTarget(target sdk.Target) *grabana.Target {
	return &grabana.Target{
		Prometheus: &grabana.PrometheusTarget{
			Query:  target.Expr,
			Legend: target.LegendFormat,
			Ref:    target.RefID,
			Hidden: target.Hide,
		},
	}
}

func (converter *JSON) convertStackdriverTarget(target sdk.Target) *grabana.Target {
	switch strings.ToLower(target.MetricKind) {
	case "cumulative":
	case "gauge":
	case "delta":
	default:
		converter.logger.Warn("unhandled stackdriver metric kind: target skipped", zap.Any("metricKind", target.MetricKind))
		return nil
	}

	var aggregation string
	if target.CrossSeriesReducer != "" {
		aggregationMap := map[string]string{
			string(stackdriver.ReduceNone):              "none",
			string(stackdriver.ReduceMean):              "mean",
			string(stackdriver.ReduceMin):               "min",
			string(stackdriver.ReduceMax):               "max",
			string(stackdriver.ReduceSum):               "sum",
			string(stackdriver.ReduceStdDev):            "stddev",
			string(stackdriver.ReduceCount):             "count",
			string(stackdriver.ReduceCountTrue):         "count_true",
			string(stackdriver.ReduceCountFalse):        "count_false",
			string(stackdriver.ReduceCountFractionTrue): "fraction_true",
			string(stackdriver.ReducePercentile99):      "percentile_99",
			string(stackdriver.ReducePercentile95):      "percentile_95",
			string(stackdriver.ReducePercentile50):      "percentile_50",
			string(stackdriver.ReducePercentile05):      "percentile_05",
		}

		if agg, ok := aggregationMap[target.CrossSeriesReducer]; ok {
			aggregation = agg
		} else {
			converter.logger.Warn("unhandled stackdriver crossSeriesReducer: target skipped", zap.Any("crossSeriesReducer", target.CrossSeriesReducer))
		}
	}

	var alignment *grabana.StackdriverAlignment
	if target.PerSeriesAligner != "" {
		alignmentMethodMap := map[string]string{
			string(stackdriver.AlignNone):          "none",
			string(stackdriver.AlignDelta):         "delta",
			string(stackdriver.AlignRate):          "rate",
			string(stackdriver.AlignInterpolate):   "interpolate",
			string(stackdriver.AlignNextOlder):     "next_older",
			string(stackdriver.AlignMin):           "min",
			string(stackdriver.AlignMax):           "max",
			string(stackdriver.AlignMean):          "mean",
			string(stackdriver.AlignCount):         "count",
			string(stackdriver.AlignSum):           "sum",
			string(stackdriver.AlignStdDev):        "stddev",
			string(stackdriver.AlignCountTrue):     "count_true",
			string(stackdriver.AlignCountFalse):    "count_false",
			string(stackdriver.AlignFractionTrue):  "fraction_true",
			string(stackdriver.AlignPercentile99):  "percentile_99",
			string(stackdriver.AlignPercentile95):  "percentile_95",
			string(stackdriver.AlignPercentile50):  "percentile_50",
			string(stackdriver.AlignPercentile05):  "percentile_05",
			string(stackdriver.AlignPercentChange): "percent_change",
		}

		if method, ok := alignmentMethodMap[target.PerSeriesAligner]; ok {
			alignment = &grabana.StackdriverAlignment{
				Period: target.AlignmentPeriod,
				Method: method,
			}
		} else {
			converter.logger.Warn("unhandled stackdriver perSeriesAligner: target skipped", zap.Any("perSeriesAligner", target.PerSeriesAligner))
		}
	}

	return &grabana.Target{
		Stackdriver: &grabana.StackdriverTarget{
			Type:        strings.ToLower(target.MetricKind),
			Metric:      target.MetricType,
			Filters:     converter.convertStackdriverFilters(target),
			Aggregation: aggregation,
			Alignment:   alignment,
			GroupBy:     target.GroupBys,
			Legend:      target.AliasBy,
			Ref:         target.RefID,
			Hidden:      target.Hide,
		},
	}
}

func (converter *JSON) convertStackdriverFilters(target sdk.Target) grabana.StackdriverFilters {
	filters := grabana.StackdriverFilters{
		Eq:         map[string]string{},
		Neq:        map[string]string{},
		Matches:    map[string]string{},
		NotMatches: map[string]string{},
	}

	var leftOperand, rightOperand, operator *string
	for i := range target.Filters {
		if target.Filters[i] == "AND" {
			continue
		}

		if leftOperand == nil {
			leftOperand = &target.Filters[i]
			continue
		}
		if operator == nil {
			operator = &target.Filters[i]
			continue
		}
		if rightOperand == nil {
			rightOperand = &target.Filters[i]
		}

		if leftOperand != nil && operator != nil && rightOperand != nil {
			switch *operator {
			case "=":
				filters.Eq[*leftOperand] = *rightOperand
			case "!=":
				filters.Neq[*leftOperand] = *rightOperand
			case "=~":
				filters.Matches[*leftOperand] = *rightOperand
			case "!=~":
				filters.NotMatches[*leftOperand] = *rightOperand
			default:
				converter.logger.Warn("unhandled stackdriver filter operator: filter skipped", zap.Any("operator", *operator))
			}

			leftOperand = nil
			rightOperand = nil
			operator = nil
		}
	}

	return filters

}

func panelSpan(panel sdk.Panel) float32 {
	span := panel.Span
	if span == 0 && panel.GridPos.H != nil {
		span = float32(*panel.GridPos.W / 2) // 24 units per row to 12
	}

	return span
}

func defaultOption(opt sdk.Current) string {
	if opt.Value == nil {
		return ""
	}

	return opt.Value.(string)
}
