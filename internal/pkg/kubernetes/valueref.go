package kubernetes

import (
	"context"
	"fmt"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ErrInvalidValueRef = fmt.Errorf("invalid value ref")

type ValueRefReader struct {
	logger logr.Logger
	client client.Reader
}

func NewValueRefReader(logger logr.Logger, client client.Reader) *ValueRefReader {
	return &ValueRefReader{
		logger: logger,
		client: client,
	}
}

func (reader *ValueRefReader) RefToValue(ctx context.Context, namespace string, ref v1alpha1.ValueOrRef) (string, error) {
	if ref.Value == "" && (ref.ValueRef == nil || ref.ValueRef.SecretKeyRef == nil) {
		return "", ErrInvalidValueRef
	}

	if ref.Value != "" {
		return ref.Value, nil
	}

	return reader.readSecret(ctx, namespace, ref.ValueRef.SecretKeyRef)
}

func (reader *ValueRefReader) readSecret(ctx context.Context, namespace string, ref *v1.SecretKeySelector) (string, error) {
	reader.logger.Info("fetching secret", "namespace", namespace, "name", ref.Name)

	secret := &v1.Secret{}
	if err := reader.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: ref.Name}, secret); err != nil {
		reader.logger.Error(err, "unable to fetch secret")
		return "", err
	}

	if _, ok := secret.Data[ref.Key]; !ok {
		// key doesn't exist in secret, but the ref was marked as optional
		if ref.Optional != nil && *ref.Optional {
			return "", nil
		}

		return "", fmt.Errorf("key '%s' does not exist in secret '%s'", ref.Key, ref.Name)
	}

	return string(secret.Data[ref.Key]), nil
}
