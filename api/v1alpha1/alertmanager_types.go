package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "make" to regenerate code after modifying this file

// AlertManagerSpec defines the desired state of AlertManager
type AlertManagerSpec struct {
	ContactPoints       []ContactPoint  `json:"contact_points,omitempty"`
	DefaultContactPoint string          `json:"default_contact_point,omitempty"`
	Routing             []RoutingPolicy `json:"routing,omitempty"`
}

type ContactPoint struct {
	Name     string             `json:"name"`
	Contacts []ContactPointType `json:"contacts"`
}

type ContactPointType struct {
	Email    *EmailContactType    `json:"email,omitempty"`
	Slack    *SlackContactType    `json:"slack,omitempty"`
	Opsgenie *OpsgenieContactType `json:"opsgenie,omitempty"`
}

type EmailContactType struct {
	To      []string `json:"to"`
	Single  bool     `json:"single,omitempty"`
	Message string   `json:"message,omitempty"`
}

type SlackContactType struct {
	Webhook ValueOrRef `json:"webhook,omitempty"`
	Title   string     `json:"title,omitempty"`
	Body    string     `json:"body,omitempty"`
}

type OpsgenieContactType struct {
	APIURL           string     `json:"api_url,omitempty"`
	APIKey           ValueOrRef `json:"api_key,omitempty"`
	AutoClose        bool       `json:"auto_close,omitempty"`
	OverridePriority bool       `json:"override_priority,omitempty"`
}

type RoutingPolicy struct {
	ContactPoint string               `json:"to"`
	Rules        []LabelsMatchingRule `json:"if_labels,omitempty"`
}

type LabelsMatchingRule struct {
	Eq         map[string]string `json:"eq,omitempty"`
	Neq        map[string]string `json:"neq,omitempty"`
	Matches    map[string]string `json:"matches,omitempty"`
	NotMatches map[string]string `json:"not_matches,omitempty"`
}

// AlertManagerStatus defines the observed state of AlertManager
type AlertManagerStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
//+kubebuilder:printcolumn:name="Message",type=string,JSONPath=`.status.message`

// AlertManager is the Schema for the alertmanagers API
type AlertManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertManagerSpec   `json:"spec,omitempty"`
	Status AlertManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AlertManagerList contains a list of AlertManager
type AlertManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertManager{}, &AlertManagerList{})
}
