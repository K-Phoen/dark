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
	Prometheus *PrometheusDatasource `json:"prometheus,omitempty"`
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
	AccessMode string `json:"access_mode,omitempty"`
}
