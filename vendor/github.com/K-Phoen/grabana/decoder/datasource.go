package decoder

import (
	"fmt"
	"time"

	"github.com/K-Phoen/grabana/datasource"
	"github.com/K-Phoen/grabana/datasource/prometheus"
)

var ErrDatasourceNotConfigured = fmt.Errorf("datasource not configured")
var ErrInvalidAccessMode = fmt.Errorf("invalid access mode")

type Datasource struct {
	Prometheus *PrometheusDatasource `yaml:",omitempty"`
}

func (datasource Datasource) toModelDatasource() (datasource.Datasource, error) {
	if datasource.Prometheus != nil {
		return datasource.Prometheus.toModelDatasource()
	}

	return nil, ErrDatasourceNotConfigured
}

type PrometheusDatasource struct {
	Name               string
	URL                string   `yaml:"url"`
	Default            *bool    `yaml:"default,omitempty"`
	ForwardOauth       *bool    `yaml:"forward_oauth,omitempty"`
	ForwardCredentials *bool    `yaml:"forward_credentials,omitempty"`
	SkipTLSVerify      *bool    `yaml:"skip_tls_verify,omitempty"`
	ForwardCookies     []string `yaml:"forward_cookies,omitempty,flow"`
	ScrapeInterval     string   `yaml:"scrape_interval,omitempty"`
	QueryTimeout       string   `yaml:"query_timeout,omitempty"`
	HTTPMethod         string   `yaml:"http_method,omitempty"`
	AccessMode         string   `yaml:"access_mode,omitempty"`
}

func (ds PrometheusDatasource) toOptions() ([]prometheus.Option, error) {
	opts := []prometheus.Option{}

	if ds.Default != nil && *ds.Default {
		opts = append(opts, prometheus.Default())
	}
	if ds.ForwardOauth != nil && *ds.ForwardOauth {
		opts = append(opts, prometheus.ForwardOauthIdentity())
	}
	if ds.ForwardCredentials != nil && *ds.ForwardCredentials {
		opts = append(opts, prometheus.WithCredentials())
	}
	if ds.SkipTLSVerify != nil && *ds.SkipTLSVerify {
		opts = append(opts, prometheus.SkipTLSVerify())
	}
	if len(ds.ForwardCookies) != 0 {
		opts = append(opts, prometheus.ForwardCookies(ds.ForwardCookies...))
	}
	if ds.ScrapeInterval != "" {
		interval, err := time.ParseDuration(ds.ScrapeInterval)
		if err != nil {
			return nil, err
		}

		opts = append(opts, prometheus.ScrapeInterval(interval))
	}
	if ds.QueryTimeout != "" {
		timeout, err := time.ParseDuration(ds.QueryTimeout)
		if err != nil {
			return nil, err
		}

		opts = append(opts, prometheus.QueryTimeout(timeout))
	}
	if ds.AccessMode != "" {
		if ds.AccessMode != "proxy" && ds.AccessMode != "direct" {
			return nil, ErrInvalidAccessMode
		}

		opts = append(opts, prometheus.AccessMode(prometheus.Access(ds.AccessMode)))
	}
	if ds.HTTPMethod != "" {
		opts = append(opts, prometheus.HTTPMethod(ds.HTTPMethod))
	}

	return opts, nil
}

func (ds PrometheusDatasource) toModelDatasource() (datasource.Datasource, error) {
	opts, err := ds.toOptions()
	if err != nil {
		return nil, err
	}

	return prometheus.New(ds.Name, ds.URL, opts...), nil
}
