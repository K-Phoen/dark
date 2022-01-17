package kubernetes

import (
	"context"
	"fmt"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
)

var ErrInvalidValueRef = fmt.Errorf("invalid value ref")

type secretsReader interface {
	Read(ctx context.Context, namespace string, ref v1.SecretKeySelector) (string, error)
}

type ValueRefReader struct {
	logger logr.Logger

	secrets secretsReader
}

func NewValueRefReader(logger logr.Logger, secrets secretsReader) *ValueRefReader {
	return &ValueRefReader{
		logger:  logger,
		secrets: secrets,
	}
}

func (reader *ValueRefReader) RefToValue(ctx context.Context, namespace string, ref v1alpha1.ValueOrRef) (string, error) {
	if ref.Value == "" && (ref.ValueRef == nil || ref.ValueRef.SecretKeyRef == nil) {
		return "", ErrInvalidValueRef
	}

	if ref.Value != "" {
		return ref.Value, nil
	}

	return reader.secrets.Read(ctx, namespace, *ref.ValueRef.SecretKeyRef)
}
