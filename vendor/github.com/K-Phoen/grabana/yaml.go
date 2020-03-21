package grabana

import (
	"fmt"
	"io"

	"github.com/K-Phoen/grabana/graph"
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/singlestat"
	"github.com/K-Phoen/grabana/table"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/text"
	"github.com/K-Phoen/grabana/variable/constant"
	"github.com/K-Phoen/grabana/variable/custom"
	"github.com/K-Phoen/grabana/variable/interval"
	"github.com/K-Phoen/grabana/variable/query"
	"gopkg.in/yaml.v2"
)

func UnmarshalYAML(input io.Reader) (DashboardBuilder, error) {
	decoder := yaml.NewDecoder(input)
	decoder.SetStrict(true)

	parsed := &DashboardYaml{}
	if err := decoder.Decode(parsed); err != nil {
		return DashboardBuilder{}, err
	}

	return parsed.ToDashboardBuilder()
}

type DashboardYaml struct {
	Title           string
	Editable        bool
	SharedCrosshair bool `yaml:"shared_crosshair"`
	Tags            []string
	AutoRefresh     string `yaml:"auto_refresh"`

	TagsAnnotation []TagAnnotation `yaml:"tags_annotations"`
	Variables      []DashboardVariable

	Rows []DashboardRow
}

type DashboardVariable struct {
	Type  string
	Name  string
	Label string

	// used for "interval", "const" and "custom"
	Default string

	// used for "interval"
	Values []string

	// used for "const" and "custom"
	ValuesMap map[string]string `yaml:"values_map"`

	// used for "query"
	Datasource string
	Request    string
}

func (dashboard *DashboardYaml) ToDashboardBuilder() (DashboardBuilder, error) {
	emptyDashboard := DashboardBuilder{}
	opts := []DashboardBuilderOption{
		dashboard.editable(),
		dashboard.sharedCrossHair(),
	}

	if len(dashboard.Tags) != 0 {
		opts = append(opts, Tags(dashboard.Tags))
	}

	if dashboard.AutoRefresh != "" {
		opts = append(opts, AutoRefresh(dashboard.AutoRefresh))
	}

	for _, tagAnnotation := range dashboard.TagsAnnotation {
		opts = append(opts, TagsAnnotation(tagAnnotation))
	}

	for _, variable := range dashboard.Variables {
		opt, err := variable.toOption()
		if err != nil {
			return emptyDashboard, err
		}

		opts = append(opts, opt)
	}

	for _, row := range dashboard.Rows {
		opt, err := row.toOption()
		if err != nil {
			return emptyDashboard, err
		}

		opts = append(opts, opt)
	}

	return NewDashboardBuilder(dashboard.Title, opts...), nil
}

func (dashboard *DashboardYaml) sharedCrossHair() DashboardBuilderOption {
	if dashboard.SharedCrosshair {
		return SharedCrossHair()
	}

	return DefaultTooltip()
}

func (dashboard *DashboardYaml) editable() DashboardBuilderOption {
	if dashboard.Editable {
		return Editable()
	}

	return ReadOnly()
}

func (variable *DashboardVariable) toOption() (DashboardBuilderOption, error) {
	switch variable.Type {
	case "interval":
		return variable.asInterval(), nil
	case "query":
		return variable.asQuery(), nil
	case "const":
		return variable.asConst(), nil
	case "custom":
		return variable.asCustom(), nil
	}

	return nil, fmt.Errorf("unknown dashboard variable type '%s'", variable.Type)
}

func (variable *DashboardVariable) asInterval() DashboardBuilderOption {
	opts := []interval.Option{
		interval.Values(variable.Values),
	}

	if variable.Label != "" {
		opts = append(opts, interval.Label(variable.Label))
	}
	if variable.Default != "" {
		opts = append(opts, interval.Default(variable.Default))
	}

	return VariableAsInterval(variable.Name, opts...)
}

func (variable *DashboardVariable) asQuery() DashboardBuilderOption {
	opts := []query.Option{
		query.Request(variable.Request),
	}

	if variable.Datasource != "" {
		opts = append(opts, query.DataSource(variable.Datasource))
	}
	if variable.Label != "" {
		opts = append(opts, query.Label(variable.Label))
	}

	return VariableAsQuery(variable.Name, opts...)
}

func (variable *DashboardVariable) asConst() DashboardBuilderOption {
	opts := []constant.Option{
		constant.Values(variable.ValuesMap),
	}

	if variable.Default != "" {
		opts = append(opts, constant.Default(variable.Default))
	}
	if variable.Label != "" {
		opts = append(opts, constant.Label(variable.Label))
	}

	return VariableAsConst(variable.Name, opts...)
}

func (variable *DashboardVariable) asCustom() DashboardBuilderOption {
	opts := []custom.Option{
		custom.Values(variable.ValuesMap),
	}

	if variable.Default != "" {
		opts = append(opts, custom.Default(variable.Default))
	}
	if variable.Label != "" {
		opts = append(opts, custom.Label(variable.Label))
	}

	return VariableAsCustom(variable.Name, opts...)
}

type DashboardRow struct {
	Name   string
	Panels []DashboardPanel
}

func (r DashboardRow) toOption() (DashboardBuilderOption, error) {
	opts := []row.Option{}

	for _, panel := range r.Panels {
		opt, err := panel.toOption()
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return Row(r.Name, opts...), nil
}

type DashboardPanel struct {
	Graph      *DashboardGraph
	Table      *DashboardTable
	SingleStat *DashboardSingleStat `yaml:"single_stat"`
	Text       *DashboardText
}

func (panel DashboardPanel) toOption() (row.Option, error) {
	if panel.Graph != nil {
		return panel.Graph.toOption()
	}
	if panel.Table != nil {
		return panel.Table.toOption()
	}
	if panel.SingleStat != nil {
		return panel.SingleStat.toOption()
	}
	if panel.Text != nil {
		return panel.Text.toOption()
	}

	return nil, fmt.Errorf("panel not configured")
}

type DashboardGraph struct {
	Title      string
	Span       float32
	Height     string
	Datasource string
	Targets    []Target
}

func (graphPanel DashboardGraph) toOption() (row.Option, error) {
	opts := []graph.Option{}

	if graphPanel.Span != 0 {
		opts = append(opts, graph.Span(graphPanel.Span))
	}
	if graphPanel.Height != "" {
		opts = append(opts, graph.Height(graphPanel.Height))
	}
	if graphPanel.Datasource != "" {
		opts = append(opts, graph.DataSource(graphPanel.Datasource))
	}

	for _, t := range graphPanel.Targets {
		opt, err := graphPanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithGraph(graphPanel.Title, opts...), nil
}

func (graphPanel *DashboardGraph) target(t Target) (graph.Option, error) {
	if t.Prometheus != nil {
		return graph.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}

	return nil, fmt.Errorf("target not configured")
}

type Target struct {
	Prometheus *prometheusTarget
}

type prometheusTarget struct {
	Query  string
	Legend string
	Ref    string
}

func (t prometheusTarget) toOptions() []prometheus.Option {
	var opts []prometheus.Option

	if t.Legend != "" {
		opts = append(opts, prometheus.Legend(t.Legend))
	}
	if t.Ref != "" {
		opts = append(opts, prometheus.Ref(t.Ref))
	}

	return opts
}

type DashboardTable struct {
	Title                  string
	Span                   float32
	Height                 string
	Datasource             string
	Targets                []Target
	HiddenColumns          []string            `yaml:"hidden_columns"`
	TimeSeriesAggregations []table.Aggregation `yaml:"time_series_aggregations"`
}

func (tablePanel DashboardTable) toOption() (row.Option, error) {
	opts := []table.Option{}

	if tablePanel.Span != 0 {
		opts = append(opts, table.Span(tablePanel.Span))
	}
	if tablePanel.Height != "" {
		opts = append(opts, table.Height(tablePanel.Height))
	}
	if tablePanel.Datasource != "" {
		opts = append(opts, table.DataSource(tablePanel.Datasource))
	}

	for _, t := range tablePanel.Targets {
		opt, err := tablePanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	for _, column := range tablePanel.HiddenColumns {
		opts = append(opts, table.HideColumn(column))
	}

	if len(tablePanel.TimeSeriesAggregations) != 0 {
		opts = append(opts, table.AsTimeSeriesAggregations(tablePanel.TimeSeriesAggregations))
	}

	return row.WithTable(tablePanel.Title, opts...), nil
}

func (tablePanel *DashboardTable) target(t Target) (table.Option, error) {
	if t.Prometheus != nil {
		return table.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}

	return nil, fmt.Errorf("target not configured")
}

type DashboardSingleStat struct {
	Title      string
	Span       float32
	Height     string
	Datasource string
	Unit       string
	Targets    []Target
	Thresholds [2]string
	Colors     [3]string
	Color      []string
}

func (singleStatPanel DashboardSingleStat) toOption() (row.Option, error) {
	opts := []singlestat.Option{}

	if singleStatPanel.Span != 0 {
		opts = append(opts, singlestat.Span(singleStatPanel.Span))
	}
	if singleStatPanel.Height != "" {
		opts = append(opts, singlestat.Height(singleStatPanel.Height))
	}
	if singleStatPanel.Datasource != "" {
		opts = append(opts, singlestat.DataSource(singleStatPanel.Datasource))
	}
	if singleStatPanel.Unit != "" {
		opts = append(opts, singlestat.Unit(singleStatPanel.Unit))
	}
	if singleStatPanel.Thresholds[0] != "" {
		opts = append(opts, singlestat.Thresholds(singleStatPanel.Thresholds))
	}
	if singleStatPanel.Colors[0] != "" {
		opts = append(opts, singlestat.Colors(singleStatPanel.Colors))
	}

	for _, colorTarget := range singleStatPanel.Color {
		switch colorTarget {
		case "value":
			opts = append(opts, singlestat.ColorValue())
		case "background":
			opts = append(opts, singlestat.ColorBackground())
		default:
			return nil, fmt.Errorf("invalid coloring target '%s'", colorTarget)
		}
	}

	for _, t := range singleStatPanel.Targets {
		opt, err := singleStatPanel.target(t)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return row.WithSingleStat(singleStatPanel.Title, opts...), nil
}

func (singleStatPanel DashboardSingleStat) target(t Target) (singlestat.Option, error) {
	if t.Prometheus != nil {
		return singlestat.WithPrometheusTarget(t.Prometheus.Query, t.Prometheus.toOptions()...), nil
	}

	return nil, fmt.Errorf("target not configured")
}

type DashboardText struct {
	Title    string
	Span     float32
	Height   string
	HTML     string
	Markdown string
}

func (textPanel DashboardText) toOption() (row.Option, error) {
	opts := []text.Option{}

	if textPanel.Span != 0 {
		opts = append(opts, text.Span(textPanel.Span))
	}
	if textPanel.Height != "" {
		opts = append(opts, text.Height(textPanel.Height))
	}
	if textPanel.Markdown != "" {
		opts = append(opts, text.Markdown(textPanel.Markdown))
	}
	if textPanel.HTML != "" {
		opts = append(opts, text.HTML(textPanel.HTML))
	}

	return row.WithText(textPanel.Title, opts...), nil
}
