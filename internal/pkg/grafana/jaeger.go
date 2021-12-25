package grafana

import (
	"context"
	"fmt"
	"time"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/jaeger"
	"k8s.io/apimachinery/pkg/types"
)

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
			return nil, fmt.Errorf("could not parse timout: %w", err)
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
			return nil, fmt.Errorf("could not extract CA certificate: %w", err)
		}

		opts = append(opts, jaeger.WithCertificate(caCertificate))
	}
	if spec.NodeGraph != nil && *spec.NodeGraph {
		opts = append(opts, jaeger.WithNodeGraph())
	}
	if spec.TraceToLogs != nil {
		opt, err := datasources.jaegerTraceToLogs(ctx, spec.TraceToLogs)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

//nolint:dupl
func (datasources *Datasources) jaegerTraceToLogs(ctx context.Context, spec *v1alpha1.TraceToLogs) (jaeger.Option, error) {
	opts := []jaeger.TraceToLogsOption{}

	datasourceUID, err := datasources.datasourceUIDFromRef(ctx, &spec.Datasource)
	if err != nil {
		return nil, fmt.Errorf("could not infer datasource UID from reference: %w", err)
	}

	if len(spec.Tags) != 0 {
		opts = append(opts, jaeger.Tags(spec.Tags...))
	}
	if spec.FilterByTrace != nil && *spec.FilterByTrace {
		opts = append(opts, jaeger.FilterByTrace())
	}
	if spec.FilterBySpan != nil && *spec.FilterBySpan {
		opts = append(opts, jaeger.FilterBySpan())
	}
	if spec.SpanStartShift != "" {
		shift, err := time.ParseDuration(spec.SpanStartShift)
		if err != nil {
			return nil, err
		}

		opts = append(opts, jaeger.SpanStartShift(shift))
	}
	if spec.SpanEndShift != "" {
		shift, err := time.ParseDuration(spec.SpanEndShift)
		if err != nil {
			return nil, err
		}

		opts = append(opts, jaeger.SpanEndShift(shift))
	}

	return jaeger.TraceToLogs(datasourceUID, opts...), nil
}
