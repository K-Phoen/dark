package converter

import (
	"testing"

	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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
