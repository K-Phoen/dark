package grafana

import (
	"context"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/stackdriver"
	"k8s.io/apimachinery/pkg/types"
)

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
