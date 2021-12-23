package grafana

import (
	"context"
	"fmt"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
)

var ErrDatasourceNotConfigured = fmt.Errorf("datasource not configured")
var ErrInvalidAccessMode = fmt.Errorf("invalid access mode")

type refReader interface {
	RefToValue(ctx context.Context, namespace string, ref v1alpha1.ValueOrRef) (string, error)
}

type Datasources struct {
	logger        logr.Logger
	grabanaClient *grabana.Client
	refReader     refReader
}

func NewDatasources(logger logr.Logger, grabanaClient *grabana.Client, refReader refReader) *Datasources {
	return &Datasources{
		logger:        logger,
		grabanaClient: grabanaClient,
		refReader:     refReader,
	}
}

func (datasources *Datasources) SpecToModel(ctx context.Context, objectRef types.NamespacedName, spec v1alpha1.DatasourceSpec) (datasource.Datasource, error) {
	if spec.Prometheus != nil {
		return datasources.prometheusSpecToModel(ctx, objectRef, spec.Prometheus)
	}
	if spec.Stackdriver != nil {
		return datasources.stackdriverSpecToModel(ctx, objectRef, spec.Stackdriver)
	}
	if spec.Jaeger != nil {
		return datasources.jaegerSpecToModel(ctx, objectRef, spec.Jaeger)
	}
	if spec.Loki != nil {
		return datasources.lokiSpecToModel(ctx, objectRef, spec.Loki)
	}

	return nil, ErrDatasourceNotConfigured
}

func (datasources *Datasources) Upsert(ctx context.Context, model datasource.Datasource) error {
	datasources.logger.Info("upserting datasource", "name", model.Name())
	return datasources.grabanaClient.UpsertDatasource(ctx, model)
}

func (datasources *Datasources) Delete(ctx context.Context, name string) error {
	datasources.logger.Info("deleting datasource", "name", name)

	err := datasources.grabanaClient.DeleteDatasource(ctx, name)
	if err == grabana.ErrDatasourceNotFound {
		return nil
	}

	return err
}

func (datasources *Datasources) basicAuthCredentials(ctx context.Context, namespace string, auth *v1alpha1.BasicAuth) (string, string, error) {
	username, err := datasources.refReader.RefToValue(ctx, namespace, auth.Username)
	if err != nil {
		return "", "", err
	}
	password, err := datasources.refReader.RefToValue(ctx, namespace, auth.Password)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}
