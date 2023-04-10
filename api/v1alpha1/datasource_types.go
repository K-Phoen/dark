package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "make" to regenerate code after modifying this file

func init() {
	SchemeBuilder.Register(&Datasource{}, &DatasourceList{})
}

// DatasourceStatus defines the observed state of Datasource
type DatasourceStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=datasources;datasource;grafana-datasources
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
//+kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`

// Datasource is the Schema for the datasources API
type Datasource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatasourceSpec   `json:"spec"`
	Status DatasourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DatasourceList contains a list of Datasource
type DatasourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Datasource `json:"items"`
}

type DatasourceSpec struct {
	Prometheus  *PrometheusDatasource  `json:"prometheus,omitempty"`
	Stackdriver *StackdriverDatasource `json:"stackdriver,omitempty"`
	Jaeger      *JaegerDatasource      `json:"jaeger,omitempty"`
	Loki        *LokiDatasource        `json:"loki,omitempty"`
	Tempo       *TempoDatasource       `json:"tempo,omitempty"`
	CloudWatch  *CloudWatchDatasource  `json:"cloudwatch,omitempty"`
}

type PrometheusDatasource struct {
	// +kubebuilder:validation:Required
	URL                string   `json:"url"`
	Default            *bool    `json:"default,omitempty"`
	ForwardOauth       *bool    `json:"forward_oauth,omitempty"`
	ForwardCredentials *bool    `json:"forward_credentials,omitempty"`
	SkipTLSVerify      *bool    `json:"skip_tls_verify,omitempty"`
	ForwardCookies     []string `json:"forward_cookies,omitempty"`
	ScrapeInterval     string   `json:"scrape_interval,omitempty"`
	QueryTimeout       string   `json:"query_timeout,omitempty"`
	// +kubebuilder:validation:Enum=POST;GET
	HTTPMethod string `json:"http_method,omitempty"`
	// +kubebuilder:validation:Enum=proxy;direct
	AccessMode    string               `json:"access_mode,omitempty"`
	BasicAuth     *BasicAuth           `json:"basic_auth,omitempty"`
	CACertificate *ValueOrRef          `json:"ca_certificate,omitempty"`
	Exemplars     []PrometheusExemplar `json:"exemplars,omitempty"`
}

type PrometheusExemplar struct {
	LabelName string `json:"label_name"`

	// The data source the exemplar is going to navigate to.
	// Set this value for internal exemplar links.
	Datasource *ValueOrDatasourceRef `json:"datasource,omitempty"`

	// The URL of the trace backend the user would go to see its trace.
	// Set this value for external exemplar links.
	URL string `json:"url,omitempty"`
}

type LokiDatasource struct {
	// +kubebuilder:validation:Required
	URL     string `json:"url"`
	Default *bool  `json:"default,omitempty"`

	ForwardOauth       *bool       `json:"forward_oauth,omitempty"`
	ForwardCredentials *bool       `json:"forward_credentials,omitempty"`
	SkipTLSVerify      *bool       `json:"skip_tls_verify,omitempty"`
	ForwardCookies     []string    `json:"forward_cookies,omitempty"`
	Timeout            string      `json:"timeout,omitempty"`
	BasicAuth          *BasicAuth  `json:"basic_auth,omitempty"`
	CACertificate      *ValueOrRef `json:"ca_certificate,omitempty"`

	MaximumLines  *int               `json:"maximum_lines,omitempty"`
	DerivedFields []LokiDerivedField `json:"derived_fields,omitempty"`
}

type LokiDerivedField struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	URL string `json:"url"`
	// Used to parse and capture some part of the log message. You can use the captured groups in the template.
	// +kubebuilder:validation:Required
	Regex string `json:"regex"`
	// Used to override the button label when this derived field is found in a log.
	// Optional.
	URLDisplayLabel string `json:"url_label,omitempty"`
	// For internal links
	// Optional.
	Datasource *ValueOrDatasourceRef `json:"datasource,omitempty"`
}

type TempoDatasource struct {
	// +kubebuilder:validation:Required
	URL     string `json:"url"`
	Default *bool  `json:"default,omitempty"`

	ForwardOauth       *bool       `json:"forward_oauth,omitempty"`
	ForwardCredentials *bool       `json:"forward_credentials,omitempty"`
	SkipTLSVerify      *bool       `json:"skip_tls_verify,omitempty"`
	ForwardCookies     []string    `json:"forward_cookies,omitempty"`
	Timeout            string      `json:"timeout,omitempty"`
	BasicAuth          *BasicAuth  `json:"basic_auth,omitempty"`
	CACertificate      *ValueOrRef `json:"ca_certificate,omitempty"`

	NodeGraph *bool `json:"node_graph,omitempty"`

	TraceToLogs *TraceToLogs `json:"trace_to_logs,omitempty"`
}

type TraceToLogs struct {
	Datasource     ValueOrDatasourceRef `json:"datasource"`
	Tags           []string             `json:"tags,omitempty"`
	SpanStartShift string               `json:"span_start_shift,omitempty"`
	SpanEndShift   string               `json:"span_end_shift,omitempty"`
	FilterByTrace  *bool                `json:"filter_by_trace,omitempty"`
	FilterBySpan   *bool                `json:"filter_by_span,omitempty"`
}

type StackdriverDatasource struct {
	Default           *bool       `json:"default,omitempty"`
	JWTAuthentication *ValueOrRef `json:"jwt_authentication,omitempty"`
}

type JaegerDatasource struct {
	// +kubebuilder:validation:Required
	URL     string `json:"url"`
	Default *bool  `json:"default,omitempty"`

	ForwardOauth       *bool       `json:"forward_oauth,omitempty"`
	ForwardCredentials *bool       `json:"forward_credentials,omitempty"`
	SkipTLSVerify      *bool       `json:"skip_tls_verify,omitempty"`
	ForwardCookies     []string    `json:"forward_cookies,omitempty"`
	Timeout            string      `json:"timeout,omitempty"`
	BasicAuth          *BasicAuth  `json:"basic_auth,omitempty"`
	CACertificate      *ValueOrRef `json:"ca_certificate,omitempty"`

	NodeGraph *bool `json:"node_graph,omitempty"`

	TraceToLogs *TraceToLogs `json:"trace_to_logs,omitempty"`
}

type CloudWatchDatasource struct {
	Default *bool `json:"default,omitempty"`

	Auth *CloudWatchAuth `json:"auth,omitempty"`

	// Endpoint specifies a custom endpoint for the CloudWatch service.
	Endpoint string `json:"endpoint,omitempty"`

	// DefaultRegion sets the default region to use.
	DefaultRegion string `json:"default_region,omitempty"`

	// AssumeRoleARN specifies the ARN of a role to assume.
	// Format: arn:aws:iam:*
	AssumeRoleARN string `json:"assume_role_arn,omitempty"`

	// ExternalID specifies the external identifier of a role to assume in another account.
	ExternalID string `json:"external_id,omitempty"`

	// CustomMetricsNamespaces specifies a list of namespaces for custom metrics.
	CustomMetricsNamespaces []string `json:"custom_metrics_namespaces,omitempty"`
}

type CloudWatchAuth struct {
	Keys *CloudWatchAuthKeys `json:"keys,omitempty"`
}

type CloudWatchAuthKeys struct {
	Access string      `json:"access"`
	Secret *ValueOrRef `json:"secret"`
}

type BasicAuth struct {
	Username ValueOrRef `json:"username"`
	Password ValueOrRef `json:"password"`
}

type ValueOrDatasourceRef struct {
	// Only one of the following may be specified.
	UID  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
}
