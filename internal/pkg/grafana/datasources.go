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
)

var ErrDatasourceNotConfigured = fmt.Errorf("datasource not configured")
var ErrInvalidAccessMode = fmt.Errorf("invalid access mode")

type Datasources struct {
	logger        logr.Logger
	grabanaClient *grabana.Client
}

func NewDatasources(logger logr.Logger, grabanaClient *grabana.Client) *Datasources {
	return &Datasources{
		logger:        logger,
		grabanaClient: grabanaClient,
	}
}

func (datasources *Datasources) SpecToModel(name string, spec v1alpha1.DatasourceSpec) (datasource.Datasource, error) {
	if spec.Prometheus != nil {
		return prometheusSpecToModel(name, spec.Prometheus)
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

func prometheusSpecToModel(name string, ds *v1alpha1.PrometheusDatasource) (datasource.Datasource, error) {
	opts, err := prometheusSpecToOptions(ds)
	if err != nil {
		return nil, err
	}

	return prometheus.New(name, ds.URL, opts...), nil
}

func prometheusSpecToOptions(promSpec *v1alpha1.PrometheusDatasource) ([]prometheus.Option, error) {
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

	return opts, nil
}
