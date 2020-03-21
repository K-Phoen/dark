package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

// GrafanaDashboard is a specification for a GrafanaDashboard resource
type GrafanaDashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrafanaDashboardSpec   `json:"spec"`
	Status GrafanaDashboardStatus `json:"status"`
}

// GrafanaDashboardSpec is the spec for a GrafanaDashboard resource
type GrafanaDashboardSpec struct {
	Message string `json:"message"`
}

// GrafanaDashboardStatus is the status for a GrafanaDashboard resource
type GrafanaDashboardStatus struct {
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaDashboardList is a list of GrafanaDashboard resources
type GrafanaDashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []GrafanaDashboard `json:"items"`
}
