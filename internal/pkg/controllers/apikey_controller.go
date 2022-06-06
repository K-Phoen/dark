package controllers

import (
	"context"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/dark/internal/pkg/grafana"
	"github.com/K-Phoen/dark/internal/pkg/kubernetes"
	"github.com/K-Phoen/grabana"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

//nolint:gosec
const apiKeysFinalizerName = "apikeys.k8s.kevingomez.fr/finalizer" //  (these are not hardcoded credentials -_-)

type apiKeyClient interface {
	Reconcile(ctx context.Context, key grafana.APIKey) error
	Delete(ctx context.Context, name string) error
}

// APIKeyReconciler reconciles a APIKey object
type APIKeyReconciler struct {
	client.Client

	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	apiKeyClient apiKeyClient
}

func StartAPIKeyReconciler(logger logr.Logger, ctrlManager ctrl.Manager, grabanaClient *grabana.Client) error {
	apiKeys := grafana.NewAPIKeys(logger, grabanaClient, kubernetes.NewSecrets(logger, ctrlManager.GetClient()))

	reconciler := &APIKeyReconciler{
		Client:       ctrlManager.GetClient(),
		Scheme:       ctrlManager.GetScheme(),
		Recorder:     ctrlManager.GetEventRecorderFor("api-key-controller"),
		apiKeyClient: apiKeys,
	}

	return reconciler.SetupWithManager(ctrlManager)
}

//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=apikeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=apikeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=apikeys/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the APIKey object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *APIKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("reconciling")

	apiKeyManifest := &v1alpha1.APIKey{}
	if err := r.Get(ctx, req.NamespacedName, apiKeyManifest); err != nil {
		logger.Error(err, "unable to fetch APIKey")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if apiKeyManifest.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(apiKeyManifest.GetFinalizers(), apiKeysFinalizerName) {
			controllerutil.AddFinalizer(apiKeyManifest, apiKeysFinalizerName)
			if err := r.Update(ctx, apiKeyManifest); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		logger.Info("deleting API key")

		// The object is being deleted
		if containsString(apiKeyManifest.GetFinalizers(), apiKeysFinalizerName) {
			logger.Info("finalizer found, deleting API key from grafana")

			// our finalizer is present, so lets handle any external dependency
			if err := r.apiKeyClient.Delete(ctx, apiKeyManifest.Name); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(apiKeyManifest, apiKeysFinalizerName)
			if err := r.Update(ctx, apiKeyManifest); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// handle actual reconciliation
	return r.doReconcileManifest(ctx, apiKeyManifest)

}

func (r *APIKeyReconciler) doReconcileManifest(ctx context.Context, manifest *v1alpha1.APIKey) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	key := grafana.APIKey{
		Name:            manifest.Name,
		Role:            manifest.Spec.Role,
		SecretName:      manifest.Name,
		SecretNamespace: manifest.Namespace,
		TokenKey:        "token",
	}

	if err := r.apiKeyClient.Reconcile(ctx, key); err != nil {
		logger.Info("failed reconciling API key")

		r.updateStatus(ctx, manifest, err)
		r.Recorder.Event(manifest, "Warning", "Error", "could not reconcile API key with Grafana")

		return ctrl.Result{}, err
	}

	r.updateStatus(ctx, manifest, nil)
	r.Recorder.Event(manifest, "Normal", "Synchronized", "API key reconciled")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *APIKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.APIKey{}).
		Complete(r)
}

func (r *APIKeyReconciler) updateStatus(ctx context.Context, manifest *v1alpha1.APIKey, err error) {
	logger := log.FromContext(ctx)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep manifestCopy of original object and modify this manifestCopy
	// Or create a manifestCopy manually for better performance
	manifestCopy := manifest.DeepCopy()

	if err == nil {
		manifestCopy.Status.Status = "OK"
		manifestCopy.Status.Message = "Synchronized"
	} else {
		manifestCopy.Status.Status = "Error"
		manifestCopy.Status.Message = err.Error()
	}

	if err := r.Status().Update(ctx, manifestCopy); err != nil {
		logger.Error(err, "unable to update APIKey status")
	}
}
