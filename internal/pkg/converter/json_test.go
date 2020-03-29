package converter

import (
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
	variable.Current = sdk.Current{Text: "30s"}
	variable.Options = []sdk.Option{
		{Text: "10s", Value: "10s"},
		{Text: "30s", Value: "30s"},
		{Text: "1m", Value: "1m"},
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
	variable.Current = sdk.Current{Text: "85th", Value: "85"}
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
	req.Equal("85th", custom.Default)
	req.True(reflect.DeepEqual(custom.ValuesMap, map[string]string{
		"50th": "50",
		"85th": "85",
		"99th": "99",
	}))
}

func TestConvertConstVar(t *testing.T) {
	req := require.New(t)

	variable := defaultVar("const")
	variable.Name = "var_const"
	variable.Label = "Label const"
	variable.Current = sdk.Current{Text: "85th", Value: "85"}
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
