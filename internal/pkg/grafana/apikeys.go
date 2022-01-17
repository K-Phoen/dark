package grafana

import (
	"context"
	"fmt"

	"github.com/K-Phoen/dark/internal/pkg/kubernetes"
	"github.com/K-Phoen/grabana"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
)

type APIKey struct {
	Name string
	Role string

	SecretName      string
	SecretNamespace string
	TokenKey        string
}

func (key APIKey) GrabanaRole() (grabana.APIKeyRole, error) {
	switch key.Role {
	case "admin":
		return grabana.AdminRole, nil
	case "editor":
		return grabana.EditorRole, nil
	case "viewer":
		return grabana.ViewerRole, nil
	}

	return grabana.ViewerRole, fmt.Errorf("invalid role")
}

type secretsWriter interface {
	Read(ctx context.Context, namespace string, ref v1.SecretKeySelector) (string, error)
	Upsert(ctx context.Context, request kubernetes.SecretUpsertRequest) error
}

type APIKeys struct {
	logger        logr.Logger
	grabanaClient *grabana.Client
	secrets       secretsWriter
}

func NewAPIKeys(logger logr.Logger, grabanaClient *grabana.Client, secrets secretsWriter) *APIKeys {
	return &APIKeys{
		logger:        logger,
		grabanaClient: grabanaClient,
		secrets:       secrets,
	}
}

func (keys *APIKeys) Reconcile(ctx context.Context, key APIKey) error {
	logger := keys.logger.WithValues("key", key.Name)

	existingGrafanaKeys, err := keys.grabanaClient.APIKeys(ctx)
	if err != nil {
		logger.Error(err, "could not check existing keys in Grafana")
		return err
	}

	// the API key does not exist in Grafana, we need to create it and its Kubernetes secret
	if _, ok := existingGrafanaKeys[key.Name]; !ok {
		return keys.createKey(ctx, key)
	}

	// the API key exists, but the secret does not. We need to re-create both
	nope := false
	_, err = keys.secrets.Read(ctx, key.SecretNamespace, v1.SecretKeySelector{
		LocalObjectReference: v1.LocalObjectReference{
			Name: key.SecretName,
		},
		Key:      key.TokenKey,
		Optional: &nope,
	})

	if err == kubernetes.ErrSecretNotFound || err == kubernetes.ErrKeyNotFoundInSecret {
		if err := keys.Delete(ctx, key.Name); err != nil {
			return err
		}

		return keys.createKey(ctx, key)
	}

	// api key exists, secret exist = nothing to do (we assume secret is up-to-date)
	return nil
}

func (keys *APIKeys) Delete(ctx context.Context, name string) error {
	if err := keys.grabanaClient.DeleteAPIKeyByName(ctx, name); err != nil && err != grabana.ErrAPIKeyNotFound {
		keys.logger.Error(err, "could not delete key in Grafana", "key", name)
		return err
	}

	return nil
}

func (keys *APIKeys) createKey(ctx context.Context, key APIKey) error {
	logger := keys.logger.WithValues("key", key.Name)

	fmt.Printf("creating new api key\n")

	role, err := key.GrabanaRole()
	if err != nil {
		return err
	}

	// create the token in Grafana
	token, err := keys.grabanaClient.CreateAPIKey(ctx, grabana.CreateAPIKeyRequest{
		Name: key.Name,
		Role: role,
	})
	if err != nil {
		logger.Error(err, "could not create Grafana API key")
		return err
	}

	secretPayload := make(map[string][]byte)
	secretPayload[key.TokenKey] = []byte(token)

	// map it to a kubernetes secret
	err = keys.secrets.Upsert(ctx, kubernetes.SecretUpsertRequest{
		Name:      key.SecretName,
		Namespace: key.SecretNamespace,
		Data:      secretPayload,
	})
	if err != nil {
		logger.Error(err, "could not create Kubernetes secret for API key")
		return err
	}

	return nil
}
