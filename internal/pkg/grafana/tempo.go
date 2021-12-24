package grafana

import (
	"context"
	"fmt"
	"time"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/tempo"
	"k8s.io/apimachinery/pkg/types"
)

func (datasources *Datasources) tempoSpecToModel(ctx context.Context, objectRef types.NamespacedName, ds *v1alpha1.TempoDatasource) (datasource.Datasource, error) {
	opts, err := datasources.tempoSpecToOptions(ctx, objectRef, ds)
	if err != nil {
		return nil, err
	}

	return tempo.New(objectRef.Name, ds.URL, opts...), nil
}

func (datasources *Datasources) tempoSpecToOptions(ctx context.Context, objectRef types.NamespacedName, spec *v1alpha1.TempoDatasource) ([]tempo.Option, error) {
	opts := []tempo.Option{}

	if spec.Default != nil && *spec.Default {
		opts = append(opts, tempo.Default())
	}
	if spec.ForwardOauth != nil && *spec.ForwardOauth {
		opts = append(opts, tempo.ForwardOauthIdentity())
	}
	if spec.ForwardCredentials != nil && *spec.ForwardCredentials {
		opts = append(opts, tempo.WithCredentials())
	}
	if spec.SkipTLSVerify != nil && *spec.SkipTLSVerify {
		opts = append(opts, tempo.SkipTLSVerify())
	}
	if len(spec.ForwardCookies) != 0 {
		opts = append(opts, tempo.ForwardCookies(spec.ForwardCookies...))
	}
	if spec.Timeout != "" {
		timeout, err := time.ParseDuration(spec.Timeout)
		if err != nil {
			return nil, fmt.Errorf("could not parse timout: %w", err)
		}

		opts = append(opts, tempo.Timeout(timeout))
	}
	if spec.BasicAuth != nil {
		username, password, err := datasources.basicAuthCredentials(ctx, objectRef.Namespace, spec.BasicAuth)
		if err != nil {
			return nil, err
		}

		opts = append(opts, tempo.BasicAuth(username, password))
	}
	if spec.CACertificate != nil {
		caCertificate, err := datasources.refReader.RefToValue(ctx, objectRef.Namespace, *spec.CACertificate)
		if err != nil {
			return nil, fmt.Errorf("could not extract CA certificate: %w", err)
		}

		opts = append(opts, tempo.WithCertificate(caCertificate))
	}

	if spec.TraceToLogs != nil {
		opt, err := datasources.tempoTraceToLogs(ctx, spec.TraceToLogs)
		if err != nil {
			return nil, err
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func (datasources *Datasources) tempoTraceToLogs(ctx context.Context, spec *v1alpha1.TempoTraceToLogs) (tempo.Option, error) {
	opts := []tempo.TraceToLogsOption{}

	datasourceUID, err := datasources.datasourceUIDFromRef(ctx, &spec.Datasource)
	if err != nil {
		return nil, fmt.Errorf("could not infer datasource UID from reference: %w", err)
	}

	if len(spec.Tags) != 0 {
		opts = append(opts, tempo.Tags(spec.Tags...))
	}
	if spec.FilterByTrace != nil && *spec.FilterByTrace {
		opts = append(opts, tempo.FilterByTrace())
	}
	if spec.FilterBySpan != nil && *spec.FilterBySpan {
		opts = append(opts, tempo.FilterBySpan())
	}
	if spec.SpanStartShift != "" {
		shift, err := time.ParseDuration(spec.SpanStartShift)
		if err != nil {
			return nil, err
		}

		opts = append(opts, tempo.SpanStartShift(shift))
	}
	if spec.SpanEndShift != "" {
		shift, err := time.ParseDuration(spec.SpanEndShift)
		if err != nil {
			return nil, err
		}

		opts = append(opts, tempo.SpanEndShift(shift))
	}

	return tempo.TraceToLogs(datasourceUID, opts...), nil
}
