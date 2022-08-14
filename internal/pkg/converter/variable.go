package converter

import (
	"strings"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

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
	case "datasource":
		converter.convertDatasourceVar(variable, dashboard)
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
		Hide:    converter.convertVarHide(variable),
	}

	for _, opt := range variable.Options {
		interval.Values = append(interval.Values, opt.Value)
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Interval: interval})
}

func (converter *JSON) convertCustomVar(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	custom := &grabana.VariableCustom{
		Name:       variable.Name,
		Label:      variable.Label,
		Default:    defaultOption(variable.Current),
		ValuesMap:  make(map[string]string, len(variable.Options)),
		AllValue:   variable.AllValue,
		IncludeAll: variable.IncludeAll,
		Hide:       converter.convertVarHide(variable),
	}

	for _, opt := range variable.Options {
		custom.ValuesMap[opt.Text] = opt.Value
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Custom: custom})
}

func (converter *JSON) convertQueryVar(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	datasource := ""
	if variable.Datasource != nil {
		datasource = variable.Datasource.LegacyName
	}

	query := &grabana.VariableQuery{
		Name:       variable.Name,
		Label:      variable.Label,
		Datasource: datasource,
		Regex:      variable.Regex,
		IncludeAll: variable.IncludeAll,
		DefaultAll: variable.Current.Value == "$__all",
		AllValue:   variable.AllValue,
		Hide:       converter.convertVarHide(variable),
	}

	if variable.Query != nil {
		if request, ok := variable.Query.(string); ok {
			query.Request = request
		}
		if request, ok := variable.Query.(map[string]interface{}); ok {
			query.Request = request["query"].(string)
		}
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Query: query})
}

func (converter *JSON) convertDatasourceVar(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	datasource := &grabana.VariableDatasource{
		Name:       variable.Name,
		Label:      variable.Label,
		Regex:      variable.Regex,
		IncludeAll: variable.IncludeAll,
		Hide:       converter.convertVarHide(variable),
	}

	if variable.Query != nil {
		datasource.Type = variable.Query.(string)
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Datasource: datasource})
}

func (converter *JSON) convertConstVar(variable sdk.TemplateVar, dashboard *grabana.DashboardModel) {
	constant := &grabana.VariableConst{
		Name:      variable.Name,
		Label:     variable.Label,
		Default:   strings.Join(variable.Current.Text.Value, ","),
		ValuesMap: make(map[string]string, len(variable.Options)),
		Hide:      converter.convertVarHide(variable),
	}

	for _, opt := range variable.Options {
		constant.ValuesMap[opt.Text] = opt.Value
	}

	dashboard.Variables = append(dashboard.Variables, grabana.DashboardVariable{Const: constant})
}

func (converter *JSON) convertVarHide(variable sdk.TemplateVar) string {
	switch variable.Hide {
	case 0:
		return ""
	case 1:
		return "label"
	case 2:
		return "variable"
	default:
		converter.logger.Warn("unknown hide value for variable %s", zap.String("variable", variable.Name))
		return ""
	}
}
