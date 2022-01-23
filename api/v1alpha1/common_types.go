package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

type ValueOrRef struct {
	// Only one of the following may be specified.
	Value    string    `json:"value,omitempty"`
	ValueRef *ValueRef `json:"valueFrom,omitempty"`
}

type ValueRef struct {
	SecretKeyRef *v1.SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}
