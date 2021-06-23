package converter

import (
	"bytes"
	"reflect"
	"testing"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/grafana-tools/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func defaultVar(varType string) sdk.TemplateVar {
	return sdk.TemplateVar{
		Type:  varType,
		Name:  "var",
		Label: "Label",
	}
}

func TestConvertInvalidJSONToYAML(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToYAML(bytes.NewBufferString(""), bytes.NewBufferString(""))

	req.Error(err)
}

func TestConvertValidJSONToYaml(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToYAML(bytes.NewBufferString("{}"), bytes.NewBufferString(""))

	req.NoError(err)
}

func TestConvertInvalidJSONToK8SManifest(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToK8SManifest(bytes.NewBufferString(""), bytes.NewBufferString(""), "Folder", "test-dashboard")

	req.Error(err)
}

func TestConvertValidJSONK8SManifest(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	err := converter.ToK8SManifest(bytes.NewBufferString("{}"), bytes.NewBufferString(""), "Folder", "test-dashboard")

	req.NoError(err)
}

func TestConvertGeneralSettings(t *testing.T) {
	req := require.New(t)

	board := &sdk.Board{}
	board.Title = "title"
	board.SharedCrosshair = true
	board.Editable = true
	board.Tags = []string{"tag", "other"}
	board.Refresh = &sdk.BoolString{
		Flag:  true,
		Value: "5s",
	}

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertGeneralSettings(board, dashboard)

	req.Equal("title", dashboard.Title)
	req.Equal("5s", dashboard.AutoRefresh)
	req.Equal([]string{"tag", "other"}, dashboard.Tags)
	req.True(dashboard.Editable)
	req.True(dashboard.SharedCrosshair)
}

func TestConvertUnknownVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("unknown")

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Variables, 0)
}

func TestConvertIntervalVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("interval")
	variable.Name = "var_interval"
	variable.Label = "Label interval"
	variable.Current = sdk.Current{Text: &sdk.StringSliceString{Value: []string{"30sec"}, Valid: true}, Value: "30s"}
	variable.Options = []sdk.Option{
		{Text: "10sec", Value: "10s"},
		{Text: "30sec", Value: "30s"},
		{Text: "1min", Value: "1m"},
	}

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Variables, 1)
	req.NotNil(dashboard.Variables[0].Interval)

	interval := dashboard.Variables[0].Interval

	req.Equal("var_interval", interval.Name)
	req.Equal("Label interval", interval.Label)
	req.Equal("30s", interval.Default)
	req.ElementsMatch([]string{"10s", "30s", "1m"}, interval.Values)
}

func TestConvertCustomVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("custom")
	variable.Name = "var_custom"
	variable.Label = "Label custom"
	variable.Current = sdk.Current{Text: &sdk.StringSliceString{Value: []string{"85th"}, Valid: true}, Value: "85"}
	variable.Options = []sdk.Option{
		{Text: "50th", Value: "50"},
		{Text: "85th", Value: "85"},
		{Text: "99th", Value: "99"},
	}

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Variables, 1)
	req.NotNil(dashboard.Variables[0].Custom)

	custom := dashboard.Variables[0].Custom

	req.Equal("var_custom", custom.Name)
	req.Equal("Label custom", custom.Label)
	req.Equal("85", custom.Default)
	req.True(reflect.DeepEqual(custom.ValuesMap, map[string]string{
		"50th": "50",
		"85th": "85",
		"99th": "99",
	}))
}

func TestConvertDatasourceVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("datasource")
	variable.Name = "var_datasource"
	variable.Label = "Label datasource"

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Variables, 1)
	req.NotNil(dashboard.Variables[0].Datasource)

	dsVar := dashboard.Variables[0].Datasource

	req.Equal("var_datasource", dsVar.Name)
	req.Equal("Label datasource", dsVar.Label)
	req.False(dsVar.IncludeAll)
}

func TestConvertConstVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("const")
	variable.Name = "var_const"
	variable.Label = "Label const"
	variable.Current = sdk.Current{Text: &sdk.StringSliceString{Value: []string{"85th"}, Valid: true}, Value: "85"}
	variable.Options = []sdk.Option{
		{Text: "85th", Value: "85"},
		{Text: "99th", Value: "99"},
	}

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Variables, 1)
	req.NotNil(dashboard.Variables[0].Const)

	constant := dashboard.Variables[0].Const

	req.Equal("var_const", constant.Name)
	req.Equal("Label const", constant.Label)
	req.Equal("85th", constant.Default)
	req.True(reflect.DeepEqual(constant.ValuesMap, map[string]string{
		"85th": "85",
		"99th": "99",
	}))
}

func TestConvertQueryVar(t *testing.T) {
	req := require.New(t)
	datasource := "prometheus-default"

	variable := defaultVar("query")
	variable.Name = "var_query"
	variable.Label = "Query"
	variable.IncludeAll = true
	variable.Current = sdk.Current{Value: "$__all"}
	variable.Datasource = &datasource
	variable.Query = "prom_query"

	converter := NewJSON(zap.NewNop())

	dashboard := &grabana.DashboardModel{}
	converter.convertVariables([]sdk.TemplateVar{variable}, dashboard)

	req.Len(dashboard.Variables, 1)
	req.NotNil(dashboard.Variables[0].Query)

	query := dashboard.Variables[0].Query

	req.Equal("var_query", query.Name)
	req.Equal("Query", query.Label)
	req.Equal(datasource, query.Datasource)
	req.Equal("prom_query", query.Request)
	req.True(query.IncludeAll)
	req.True(query.DefaultAll)
}

func TestConvertRow(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	row := converter.convertRow(sdk.Panel{CommonPanel: sdk.CommonPanel{Title: "Row title"}})

	req.Equal("Row title", row.Name)
}

func TestConvertCollaspedRow(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	row := converter.convertRow(sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "Row title",
		},
		RowPanel: &sdk.RowPanel{
			Collapsed: true,
		},
	})

	req.Equal("Row title", row.Name)
	req.True(row.Collapse)
}

func TestConvertTargetFailsIfNoValidTargetIsGiven(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	convertedTarget := converter.convertTarget(sdk.Target{})
	req.Nil(convertedTarget)
}

func TestConvertTargetWithPrometheusTarget(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		Expr:         "prometheus_query",
		LegendFormat: "{{ field }}",
		RefID:        "A",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.Nil(convertedTarget.Stackdriver)
	req.Equal("prometheus_query", convertedTarget.Prometheus.Query)
	req.Equal("{{ field }}", convertedTarget.Prometheus.Legend)
	req.Equal("A", convertedTarget.Prometheus.Ref)
}

func TestConvertTargetWithGraphiteTarget(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		Target: "graphite_query",
		RefID:  "A",
		Hide:   true,
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.NotNil(convertedTarget.Graphite)
	req.Equal("graphite_query", convertedTarget.Graphite.Query)
	req.Equal("A", convertedTarget.Graphite.Ref)
	req.True(convertedTarget.Graphite.Hidden)
}

func TestConvertTargetWithInfluxDBTarget(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		Measurement: "influxdb_query",
		RefID:       "A",
		Hide:        true,
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.NotNil(convertedTarget.InfluxDB)
	req.Equal("influxdb_query", convertedTarget.InfluxDB.Query)
	req.Equal("A", convertedTarget.InfluxDB.Ref)
	req.True(convertedTarget.InfluxDB.Hidden)
}

func TestConvertTargetWithStackdriverTargetFailsIfNoMetricKind(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricType: "pubsub.googleapis.com/subscription/ack_message_count",
	}

	convertedTarget := converter.convertTarget(target)

	req.Nil(convertedTarget)
}

func TestConvertTargetWithStackdriverTargetIgnoresUnknownCrossSeriesReducer(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind:         "DELTA",
		MetricType:         "pubsub.googleapis.com/subscription/ack_message_count",
		CrossSeriesReducer: "unknown",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.NotNil(convertedTarget.Stackdriver)
	req.Empty(convertedTarget.Stackdriver.Aggregation)
}

func TestConvertTargetWithStackdriverTargetIgnoresUnknownAligner(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind:       "DELTA",
		MetricType:       "pubsub.googleapis.com/subscription/ack_message_count",
		PerSeriesAligner: "unknown",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.NotNil(convertedTarget.Stackdriver)
	req.Empty(convertedTarget.Stackdriver.Alignment)
}

func TestConvertTargetWithStackdriverTarget(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind:         "DELTA",
		MetricType:         "pubsub.googleapis.com/subscription/ack_message_count",
		CrossSeriesReducer: "REDUCE_MEAN",
		PerSeriesAligner:   "ALIGN_DELTA",
		AlignmentPeriod:    "stackdriver-auto",
		GroupBys:           []string{"field"},
		AliasBy:            "legend",
		RefID:              "A",
		Filters: []string{
			"resource.label.subscription_id",
			"=",
			"subscription_name",
			"AND",
			"other-property",
			"!=",
			"other-value",
			"AND",
			"regex-property",
			"=~",
			"regex-value",
			"AND",
			"regex-not-property",
			"!=~",
			"regex-not-value",
		},
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.Nil(convertedTarget.Prometheus)
	req.NotNil(convertedTarget.Stackdriver)
	req.Equal("delta", convertedTarget.Stackdriver.Type)
	req.Equal("pubsub.googleapis.com/subscription/ack_message_count", convertedTarget.Stackdriver.Metric)
	req.Equal("mean", convertedTarget.Stackdriver.Aggregation)
	req.Equal("stackdriver-auto", convertedTarget.Stackdriver.Alignment.Period)
	req.Equal("delta", convertedTarget.Stackdriver.Alignment.Method)
	req.Equal("legend", convertedTarget.Stackdriver.Legend)
	req.Equal("A", convertedTarget.Stackdriver.Ref)
	req.ElementsMatch([]string{"field"}, convertedTarget.Stackdriver.GroupBy)
	req.EqualValues(map[string]string{"resource.label.subscription_id": "subscription_name"}, convertedTarget.Stackdriver.Filters.Eq)
	req.EqualValues(map[string]string{"other-property": "other-value"}, convertedTarget.Stackdriver.Filters.Neq)
	req.EqualValues(map[string]string{"regex-property": "regex-value"}, convertedTarget.Stackdriver.Filters.Matches)
	req.EqualValues(map[string]string{"regex-not-property": "regex-not-value"}, convertedTarget.Stackdriver.Filters.NotMatches)
}

func TestConvertTargetWithStackdriverGauge(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind: "GAUGE",
		MetricType: "pubsub.googleapis.com/subscription/ack_message_count",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.Nil(convertedTarget.Prometheus)
	req.NotNil(convertedTarget.Stackdriver)
	req.Equal("gauge", convertedTarget.Stackdriver.Type)
	req.Equal("pubsub.googleapis.com/subscription/ack_message_count", convertedTarget.Stackdriver.Metric)
}

func TestConvertTargetWithStackdriverCumulative(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	target := sdk.Target{
		MetricKind: "CUMULATIVE",
		MetricType: "pubsub.googleapis.com/subscription/ack_message_count",
	}

	convertedTarget := converter.convertTarget(target)

	req.NotNil(convertedTarget)
	req.Nil(convertedTarget.Prometheus)
	req.NotNil(convertedTarget.Stackdriver)
	req.Equal("cumulative", convertedTarget.Stackdriver.Type)
	req.Equal("pubsub.googleapis.com/subscription/ack_message_count", convertedTarget.Stackdriver.Metric)
}

func TestConvertTagAnnotationIgnoresBuiltIn(t *testing.T) {
	req := require.New(t)

	annotation := sdk.Annotation{Name: "Annotations & Alerts"}
	dashboard := &grabana.DashboardModel{}

	NewJSON(zap.NewNop()).convertAnnotations([]sdk.Annotation{annotation}, dashboard)

	req.Len(dashboard.TagsAnnotation, 0)
}

func TestConvertTagAnnotationIgnoresUnknownTypes(t *testing.T) {
	req := require.New(t)

	annotation := sdk.Annotation{Name: "Will be ignored", Type: "dashboard"}
	dashboard := &grabana.DashboardModel{}

	NewJSON(zap.NewNop()).convertAnnotations([]sdk.Annotation{annotation}, dashboard)

	req.Len(dashboard.TagsAnnotation, 0)
}

func TestConvertTagAnnotation(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	datasource := "-- Grafana --"
	annotation := sdk.Annotation{
		Type:       "tags",
		Datasource: &datasource,
		IconColor:  "#5794F2",
		Name:       "Deployments",
		Tags:       []string{"deploy"},
	}
	dashboard := &grabana.DashboardModel{}

	converter.convertAnnotations([]sdk.Annotation{annotation}, dashboard)

	req.Len(dashboard.TagsAnnotation, 1)
	req.Equal("Deployments", dashboard.TagsAnnotation[0].Name)
	req.ElementsMatch([]string{"deploy"}, dashboard.TagsAnnotation[0].Tags)
	req.Equal("#5794F2", dashboard.TagsAnnotation[0].IconColor)
	req.Equal(datasource, dashboard.TagsAnnotation[0].Datasource)
}

func TestConvertLegend(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	rawLegend := sdk.Legend{
		AlignAsTable: true,
		Avg:          true,
		Current:      true,
		HideEmpty:    true,
		HideZero:     true,
		Max:          true,
		Min:          true,
		RightSide:    true,
		Show:         true,
		Total:        true,
	}

	legend := converter.convertLegend(rawLegend)

	req.ElementsMatch(
		[]string{"as_table", "to_the_right", "min", "max", "avg", "current", "total", "no_null_series", "no_zero_series"},
		legend,
	)
}

func TestConvertCanHideLegend(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	legend := converter.convertLegend(sdk.Legend{Show: false})
	req.ElementsMatch([]string{"hide"}, legend)
}

func TestConvertAxis(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	rawAxis := sdk.Axis{
		Format:  "bytes",
		LogBase: 2,
		Min:     &sdk.FloatString{Value: 0},
		Max:     &sdk.FloatString{Value: 42},
		Show:    true,
		Label:   "Axis",
	}

	axis := converter.convertAxis(rawAxis)

	req.Equal("bytes", *axis.Unit)
	req.Equal("Axis", axis.Label)
	req.EqualValues(0, *axis.Min)
	req.EqualValues(42, *axis.Max)
	req.False(*axis.Hidden)
}

func TestConvertTextPanelWithMarkdown(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "200px"

	textPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Text panel",
			Transparent: true,
			Height:      &height,
			Type:        "text",
		},
		TextPanel: &sdk.TextPanel{
			Mode:    "markdown",
			Content: "# hello world",
		},
	}

	converted, ok := converter.convertDataPanel(textPanel)

	req.True(ok)
	req.True(converted.Text.Transparent)
	req.Equal("Text panel", converted.Text.Title)
	req.Equal("# hello world", converted.Text.Markdown)
	req.Equal(height, converted.Text.Height)
}

func TestConvertTextPanelWithHTML(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())

	textPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Text panel html",
			Type:        "text",
			Description: strPtr("panel description"),
		},
		TextPanel: &sdk.TextPanel{
			Mode:    "html",
			Content: "<h1>hello world</h1>",
		},
	}

	converted, ok := converter.convertDataPanel(textPanel)

	req.True(ok)
	req.False(converted.Text.Transparent)
	req.Equal("Text panel html", converted.Text.Title)
	req.Equal("panel description", converted.Text.Description)
	req.Equal("<h1>hello world</h1>", converted.Text.HTML)
}

func TestConvertSingleStatPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "200px"
	datasource := "prometheus"

	singlestatPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "Singlestat panel",
			Description: strPtr("panel desc"),
			Type:        "singlestat",
			Transparent: true,
			Height:      &height,
			Datasource:  &datasource,
		},
		SinglestatPanel: &sdk.SinglestatPanel{
			Format:          "none",
			ValueName:       "current",
			Colors:          []string{"blue", "red", "green"},
			ColorBackground: true,
			ColorValue:      true,
		},
	}

	converted, ok := converter.convertDataPanel(singlestatPanel)

	req.True(ok)
	req.True(converted.SingleStat.Transparent)
	req.Equal("Singlestat panel", converted.SingleStat.Title)
	req.Equal("panel desc", converted.SingleStat.Description)
	req.Equal("none", converted.SingleStat.Unit)
	req.Equal("current", converted.SingleStat.ValueType)
	req.Equal(height, converted.SingleStat.Height)
	req.Equal(datasource, converted.SingleStat.Datasource)
	req.True(reflect.DeepEqual(converted.SingleStat.Colors, [3]string{
		"blue", "red", "green",
	}))
	req.True(reflect.DeepEqual(converted.SingleStat.Color, []string{
		"background", "value",
	}))
}

func TestConvertHeatmapPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "400px"
	datasource := "prometheus"

	heatmapPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "heatmap panel",
			Type:        "heatmap",
			Description: strPtr("heatmap description"),
			Transparent: true,
			Height:      &height,
			Datasource:  &datasource,
		},
		HeatmapPanel: &sdk.HeatmapPanel{
			HideZeroBuckets: true,
			HighlightCards:  true,
			ReverseYBuckets: true,
			DataFormat:      "tsbuckets",
		},
	}

	converted, ok := converter.convertDataPanel(heatmapPanel)

	req.True(ok)
	req.True(converted.Heatmap.Transparent)
	req.Equal("heatmap panel", converted.Heatmap.Title)
	req.Equal("heatmap description", converted.Heatmap.Description)
	req.Equal(height, converted.Heatmap.Height)
	req.Equal(datasource, converted.Heatmap.Datasource)
	req.True(converted.Heatmap.ReverseYBuckets)
	req.True(converted.Heatmap.HideZeroBuckets)
	req.True(converted.Heatmap.HightlightCards)
	req.Equal("time_series_buckets", converted.Heatmap.DataFormat)
}

func TestConvertGraphPanel(t *testing.T) {
	req := require.New(t)

	converter := NewJSON(zap.NewNop())
	height := "400px"
	datasource := "prometheus"

	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title:       "graph panel",
			Type:        "graph",
			Description: strPtr("graph description"),
			Transparent: true,
			Height:      &height,
			Datasource:  &datasource,
		},
		GraphPanel: &sdk.GraphPanel{},
	}

	converted, ok := converter.convertDataPanel(graphPanel)

	req.True(ok)
	req.NotNil(converted.Graph)

	graph := converted.Graph
	req.True(graph.Transparent)
	req.Equal("graph panel", graph.Title)
	req.Equal("graph description", graph.Description)
	req.Equal(height, graph.Height)
	req.Equal(datasource, graph.Datasource)
}

func TestConvertVisualization(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())
	enabled := true

	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "graph panel",
			Type:  "graph",
		},
		GraphPanel: &sdk.GraphPanel{
			NullPointMode: "connected",
			SteppedLine:   true,
			SeriesOverrides: []sdk.SeriesOverride{
				{
					Alias:  "alias",
					Dashes: &enabled,
				},
			},
		},
	}

	visualization := converter.convertVisualization(graphPanel)

	req.True(visualization.Staircase)
	req.Equal("connected", visualization.NullValue)
	req.Len(visualization.Overrides, 1)
}

func TestConvertGraphOverridesWithNoOverride(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())

	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "graph panel",
			Type:  "graph",
		},
		GraphPanel: &sdk.GraphPanel{},
	}

	overrides := converter.convertGraphOverrides(graphPanel)

	req.Len(overrides, 0)
}

func TestConvertGraphOverridesWithOneOverride(t *testing.T) {
	req := require.New(t)
	converter := NewJSON(zap.NewNop())
	color := "red"
	enabled := true
	number := 2

	graphPanel := sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			Title: "heatmap panel",
			Type:  "graph",
		},
		GraphPanel: &sdk.GraphPanel{
			SeriesOverrides: []sdk.SeriesOverride{
				{
					Alias:  "alias",
					Color:  &color,
					Dashes: &enabled,
					Fill:   &number,
					Lines:  &enabled,
				},
			},
		},
	}

	overrides := converter.convertGraphOverrides(graphPanel)

	req.Len(overrides, 1)

	override := overrides[0]

	req.Equal("alias", override.Alias)
	req.Equal(color, override.Color)
	req.True(*override.Dashes)
	req.True(*override.Lines)
	req.Equal(number, *override.Fill)
}

func strPtr(input string) *string {
	return &input
}
