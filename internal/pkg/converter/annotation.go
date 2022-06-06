package converter

import (
	grabanaDashboard "github.com/K-Phoen/grabana/dashboard"
	grabana "github.com/K-Phoen/grabana/decoder"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

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
		datasource = annotation.Datasource.LegacyName
	}

	dashboard.TagsAnnotation = append(dashboard.TagsAnnotation, grabanaDashboard.TagAnnotation{
		Name:       annotation.Name,
		Datasource: datasource,
		IconColor:  annotation.IconColor,
		Tags:       annotation.Tags,
	})
}
