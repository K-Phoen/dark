package grafana

import (
	"context"
	"fmt"
	"time"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/prometheus"
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

func (datasources *Datasources) prometheusSpecToModel(ctx context.Context, objectRef types.NamespacedName, ds *v1alpha1.PrometheusDatasource) (datasource.Datasource, error) {
	opts, err := datasources.prometheusSpecToOptions(ctx, objectRef, ds)
	if err != nil {
		return nil, err
	}

	return prometheus.New(objectRef.Name, ds.URL, opts...), nil
}

func (datasources *Datasources) prometheusSpecToOptions(ctx context.Context, objectRef types.NamespacedName, promSpec *v1alpha1.PrometheusDatasource) ([]prometheus.Option, error) {
	opts := []prometheus.Option{}

	if promSpec.Default != nil && *promSpec.Default {
		opts = append(opts, prometheus.Default())
	}
	if promSpec.ForwardOauth != nil && *promSpec.ForwardOauth {
		opts = append(opts, prometheus.ForwardOauthIdentity())
	}
	if promSpec.ForwardCredentials != nil && *promSpec.ForwardCredentials {
		opts = append(opts, prometheus.WithCredentials())
	}
	if promSpec.SkipTLSVerify != nil && *promSpec.SkipTLSVerify {
		opts = append(opts, prometheus.SkipTLSVerify())
	}
	if len(promSpec.ForwardCookies) != 0 {
		opts = append(opts, prometheus.ForwardCookies(promSpec.ForwardCookies...))
	}
	if promSpec.ScrapeInterval != "" {
		interval, err := time.ParseDuration(promSpec.ScrapeInterval)
		if err != nil {
			return nil, err
		}

		opts = append(opts, prometheus.ScrapeInterval(interval))
	}
	if promSpec.QueryTimeout != "" {
		timeout, err := time.ParseDuration(promSpec.QueryTimeout)
		if err != nil {
			return nil, err
		}

		opts = append(opts, prometheus.QueryTimeout(timeout))
	}
	if promSpec.AccessMode != "" {
		if promSpec.AccessMode != "proxy" && promSpec.AccessMode != "direct" {
			return nil, ErrInvalidAccessMode
		}

		opts = append(opts, prometheus.AccessMode(prometheus.Access(promSpec.AccessMode)))
	}
	if promSpec.HTTPMethod != "" {
		opts = append(opts, prometheus.HTTPMethod(promSpec.HTTPMethod))
	}
	if promSpec.BasicAuth != nil {
		basicOpts, err := datasources.basicAuthOptions(ctx, objectRef.Namespace, promSpec.BasicAuth)
		if err != nil {
			return nil, err
		}

		opts = append(opts, basicOpts)
	}

	return opts, nil
}

func (datasources *Datasources) basicAuthOptions(ctx context.Context, namespace string, auth *v1alpha1.BasicAuth) (prometheus.Option, error) {
	username, err := datasources.refReader.RefToValue(ctx, namespace, auth.Username)
	if err != nil {
		return nil, err
	}
	password, err := datasources.refReader.RefToValue(ctx, namespace, auth.Password)
	if err != nil {
		return nil, err
	}

	return prometheus.BasicAuth(username, password), nil
}
