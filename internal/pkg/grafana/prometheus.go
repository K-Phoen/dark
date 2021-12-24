package grafana

import (
	"context"
	"fmt"
	"time"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/prometheus"
	"k8s.io/apimachinery/pkg/types"
)

func (datasources *Datasources) prometheusSpecToModel(ctx context.Context, objectRef types.NamespacedName, ds *v1alpha1.PrometheusDatasource) (datasource.Datasource, error) {
	opts, err := datasources.prometheusSpecToOptions(ctx, objectRef, ds)
	if err != nil {
		return nil, err
	}

	return prometheus.New(objectRef.Name, ds.URL, opts...), nil
}

func (datasources *Datasources) prometheusSpecToOptions(ctx context.Context, objectRef types.NamespacedName, spec *v1alpha1.PrometheusDatasource) ([]prometheus.Option, error) {
	opts := []prometheus.Option{}

	if spec.Default != nil && *spec.Default {
		opts = append(opts, prometheus.Default())
	}
	if spec.ForwardOauth != nil && *spec.ForwardOauth {
		opts = append(opts, prometheus.ForwardOauthIdentity())
	}
	if spec.ForwardCredentials != nil && *spec.ForwardCredentials {
		opts = append(opts, prometheus.WithCredentials())
	}
	if spec.SkipTLSVerify != nil && *spec.SkipTLSVerify {
		opts = append(opts, prometheus.SkipTLSVerify())
	}
	if len(spec.ForwardCookies) != 0 {
		opts = append(opts, prometheus.ForwardCookies(spec.ForwardCookies...))
	}
	if spec.ScrapeInterval != "" {
		interval, err := time.ParseDuration(spec.ScrapeInterval)
		if err != nil {
			return nil, fmt.Errorf("could not parse scrape interval: %w", err)
		}

		opts = append(opts, prometheus.ScrapeInterval(interval))
	}
	if spec.QueryTimeout != "" {
		timeout, err := time.ParseDuration(spec.QueryTimeout)
		if err != nil {
			return nil, fmt.Errorf("could not parse query timout: %w", err)
		}

		opts = append(opts, prometheus.QueryTimeout(timeout))
	}
	if spec.AccessMode != "" {
		if spec.AccessMode != "proxy" && spec.AccessMode != "direct" {
			return nil, ErrInvalidAccessMode
		}

		opts = append(opts, prometheus.AccessMode(prometheus.Access(spec.AccessMode)))
	}
	if spec.HTTPMethod != "" {
		opts = append(opts, prometheus.HTTPMethod(spec.HTTPMethod))
	}
	if spec.BasicAuth != nil {
		username, password, err := datasources.basicAuthCredentials(ctx, objectRef.Namespace, spec.BasicAuth)
		if err != nil {
			return nil, err
		}

		opts = append(opts, prometheus.BasicAuth(username, password))
	}
	if spec.CACertificate != nil {
		caCertificate, err := datasources.refReader.RefToValue(ctx, objectRef.Namespace, *spec.CACertificate)
		if err != nil {
			return nil, fmt.Errorf("could not extract CA certificate: %w", err)
		}

		opts = append(opts, prometheus.WithCertificate(caCertificate))
	}

	return opts, nil
}
