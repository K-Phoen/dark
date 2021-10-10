package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaDashboard is a specification for a GrafanaDashboard resource
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen=true
type GrafanaDashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Folder string               `json:"folder"`
	Spec   runtime.RawExtension `json:"spec"`

	Status GrafanaDashboardStatus `json:"status"`
}

// GrafanaDashboardStatus is the status for a GrafanaDashboard resource
type GrafanaDashboardStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrafanaDashboardList is a list of GrafanaDashboard resources
type GrafanaDashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []GrafanaDashboard `json:"items"`
}
