package grafana

import (
	"context"
	"fmt"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/cloudwatch"
	"k8s.io/apimachinery/pkg/types"
)

func (datasources *Datasources) cloudWatchSpecToModel(ctx context.Context, objectRef types.NamespacedName, ds *v1alpha1.CloudWatchDatasource) (datasource.Datasource, error) {
	opts := []cloudwatch.Option{}

	if ds.Default != nil && *ds.Default {
		opts = append(opts, cloudwatch.Default())
	}
	if ds.Endpoint != "" {
		opts = append(opts, cloudwatch.Endpoint(ds.Endpoint))
	}
	if ds.DefaultRegion != "" {
		opts = append(opts, cloudwatch.DefaultRegion(ds.DefaultRegion))
	}
	if ds.AssumeRoleARN != "" {
		opts = append(opts, cloudwatch.AssumeRoleARN(ds.AssumeRoleARN))
	}
	if ds.ExternalID != "" {
		opts = append(opts, cloudwatch.ExternalID(ds.ExternalID))
	}
	if len(ds.CustomMetricsNamespaces) != 0 {
		opts = append(opts, cloudwatch.CustomMetricsNamespaces(ds.CustomMetricsNamespaces...))
	}

	fmt.Printf("%#v\n", ds.Auth.Keys)

	if ds.Auth != nil && ds.Auth.Keys != nil {
		secret, err := datasources.refReader.RefToValue(ctx, objectRef.Namespace, *ds.Auth.Keys.Secret)
		if err != nil {
			return nil, err
		}

		opts = append(opts, cloudwatch.AccessSecretAuth(ds.Auth.Keys.Access, secret))
	}

	return cloudwatch.New(objectRef.Name, opts...)
}
