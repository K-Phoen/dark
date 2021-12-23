package grafana

import (
	"context"
	"time"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/loki"
	"k8s.io/apimachinery/pkg/types"
)

func (datasources *Datasources) lokiSpecToModel(ctx context.Context, objectRef types.NamespacedName, ds *v1alpha1.LokiDatasource) (datasource.Datasource, error) {
	opts, err := datasources.lokiSpecToOptions(ctx, objectRef, ds)
	if err != nil {
		return nil, err
	}

	return loki.New(objectRef.Name, ds.URL, opts...), nil
}

func (datasources *Datasources) lokiSpecToOptions(ctx context.Context, objectRef types.NamespacedName, spec *v1alpha1.LokiDatasource) ([]loki.Option, error) {
	opts := []loki.Option{}

	if spec.Default != nil && *spec.Default {
		opts = append(opts, loki.Default())
	}
	if spec.ForwardOauth != nil && *spec.ForwardOauth {
		opts = append(opts, loki.ForwardOauthIdentity())
	}
	if spec.ForwardCredentials != nil && *spec.ForwardCredentials {
		opts = append(opts, loki.WithCredentials())
	}
	if spec.SkipTLSVerify != nil && *spec.SkipTLSVerify {
		opts = append(opts, loki.SkipTLSVerify())
	}
	if len(spec.ForwardCookies) != 0 {
		opts = append(opts, loki.ForwardCookies(spec.ForwardCookies...))
	}
	if spec.Timeout != "" {
		timeout, err := time.ParseDuration(spec.Timeout)
		if err != nil {
			return nil, err
		}

		opts = append(opts, loki.Timeout(timeout))
	}
	if spec.BasicAuth != nil {
		username, password, err := datasources.basicAuthCredentials(ctx, objectRef.Namespace, spec.BasicAuth)
		if err != nil {
			return nil, err
		}

		opts = append(opts, loki.BasicAuth(username, password))
	}
	if spec.CACertificate != nil {
		caCertificate, err := datasources.refReader.RefToValue(ctx, objectRef.Namespace, *spec.CACertificate)
		if err != nil {
			return nil, err
		}

		opts = append(opts, loki.WithCertificate(caCertificate))
	}
	if spec.MaximumLines != nil && *spec.MaximumLines != 0 {
		opts = append(opts, loki.MaximumLines(*spec.MaximumLines))
	}

	return opts, nil
}
