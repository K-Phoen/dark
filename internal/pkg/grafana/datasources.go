package grafana

import (
	"context"
	"fmt"
	"time"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/jaeger"
	"github.com/K-Phoen/grabana/datasource/prometheus"
	"github.com/K-Phoen/grabana/datasource/stackdriver"
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

func (datasources *Datasources) jaegerSpecToModel(ctx context.Context, objectRef types.NamespacedName, ds *v1alpha1.JaegerDatasource) (datasource.Datasource, error) {
	opts, err := datasources.jaegerSpecToOptions(ctx, objectRef, ds)
	if err != nil {
		return nil, err
	}

	return jaeger.New(objectRef.Name, ds.URL, opts...), nil
}
func (datasources *Datasources) jaegerSpecToOptions(ctx context.Context, objectRef types.NamespacedName, spec *v1alpha1.JaegerDatasource) ([]jaeger.Option, error) {
	opts := []jaeger.Option{}

	if spec.Default != nil && *spec.Default {
		opts = append(opts, jaeger.Default())
	}
	if spec.ForwardOauth != nil && *spec.ForwardOauth {
		opts = append(opts, jaeger.ForwardOauthIdentity())
	}
	if spec.ForwardCredentials != nil && *spec.ForwardCredentials {
		opts = append(opts, jaeger.WithCredentials())
	}
	if spec.SkipTLSVerify != nil && *spec.SkipTLSVerify {
		opts = append(opts, jaeger.SkipTLSVerify())
	}
	if len(spec.ForwardCookies) != 0 {
		opts = append(opts, jaeger.ForwardCookies(spec.ForwardCookies...))
	}
	if spec.Timeout != "" {
		timeout, err := time.ParseDuration(spec.Timeout)
		if err != nil {
			return nil, err
		}

		opts = append(opts, jaeger.Timeout(timeout))
	}
	if spec.BasicAuth != nil {
		username, password, err := datasources.basicAuthCredentials(ctx, objectRef.Namespace, spec.BasicAuth)
		if err != nil {
			return nil, err
		}

		opts = append(opts, jaeger.BasicAuth(username, password))
	}
	if spec.CACertificate != nil {
		caCertificate, err := datasources.refReader.RefToValue(ctx, objectRef.Namespace, *spec.CACertificate)
		if err != nil {
			return nil, err
		}

		opts = append(opts, jaeger.WithCertificate(caCertificate))
	}
	if spec.NodeGraph != nil && *spec.NodeGraph {
		opts = append(opts, jaeger.WithNodeGraph())
	}

	return opts, nil
}

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
			return nil, err
		}

		opts = append(opts, prometheus.ScrapeInterval(interval))
	}
	if spec.QueryTimeout != "" {
		timeout, err := time.ParseDuration(spec.QueryTimeout)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		opts = append(opts, prometheus.WithCertificate(caCertificate))
	}

	return opts, nil
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

func (datasources *Datasources) stackdriverSpecToModel(ctx context.Context, objectRef types.NamespacedName, ds *v1alpha1.StackdriverDatasource) (datasource.Datasource, error) {
	opts := []stackdriver.Option{}

	if ds.Default != nil && *ds.Default {
		opts = append(opts, stackdriver.Default())
	}
	if ds.JWTAuthentication != nil {
		jwtKey, err := datasources.refReader.RefToValue(ctx, objectRef.Namespace, *ds.JWTAuthentication)
		if err != nil {
			return nil, err
		}

		opts = append(opts, stackdriver.JWTAuthentication(jwtKey))
	}

	return stackdriver.New(objectRef.Name, opts...), nil
}
