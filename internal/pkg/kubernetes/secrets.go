package kubernetes

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ErrSecretNotFound = fmt.Errorf("secret not found")
var ErrKeyNotFoundInSecret = fmt.Errorf("key not found")

type SecretUpsertRequest struct {
	Name      string
	Namespace string
	Data      map[string][]byte
}

type Secrets struct {
	logger logr.Logger
	client client.Client
}

func NewSecrets(logger logr.Logger, client client.Client) *Secrets {
	return &Secrets{
		logger: logger,
		client: client,
	}
}

func (secrets *Secrets) Upsert(ctx context.Context, request SecretUpsertRequest) error {
	logger := secrets.logger.WithValues("namespace", request.Namespace, "name", request.Name)
	logger.Info("upserting secret")

	secret := &v1.Secret{}

	// if a secret for this API key already exists, we delete it
	err := secrets.client.Get(ctx, client.ObjectKey{Namespace: request.Namespace, Name: request.Name}, secret)
	if err != nil && !apierrors.IsNotFound(err) {
		logger.Error(err, "unable to check secret existence")
		return err
	}
	// the secret was found
	if err == nil {
		if err := secrets.client.Delete(ctx, secret); err != nil {
			logger.Error(err, "unable to delete existing secret")
			return err
		}
	}

	// now we can safely re-create the secret with a new value
	err = secrets.client.Create(ctx, &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.Name,
			Namespace: request.Namespace,
			Annotations: map[string]string{
				"app.kubernetes.io/managed-by": "dark",
			},
		},
		Data: request.Data,
	})
	if err != nil {
		logger.Error(err, "unable to create secret")
		return err
	}

	return nil
}

func (secrets *Secrets) Read(ctx context.Context, namespace string, ref v1.SecretKeySelector) (string, error) {
	logger := secrets.logger.WithValues("namespace", namespace, "name", ref.Name)
	logger.Info("fetching secret")

	secret := &v1.Secret{}
	if err := secrets.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: ref.Name}, secret); err != nil {
		if apierrors.IsNotFound(err) {
			return "", ErrSecretNotFound
		}

		logger.Error(err, "unable to fetch secret")
		return "", err
	}

	if _, ok := secret.Data[ref.Key]; !ok {
		// key doesn't exist in secret, but the ref was marked as optional
		if ref.Optional != nil && *ref.Optional {
			return "", nil
		}

		return "", fmt.Errorf("key '%s' does not exist: %w", ref.Key, ErrKeyNotFoundInSecret)
	}

	return string(secret.Data[ref.Key]), nil
}
