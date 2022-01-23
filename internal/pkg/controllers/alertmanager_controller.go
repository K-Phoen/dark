package controllers

import (
	"context"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/dark/internal/pkg/grafana"
	"github.com/K-Phoen/grabana"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const alertManagerFinalizerName = "alertmanagers.k8s.kevingomez.fr/finalizer"

// AlertManagerReconciler reconciles a AlertManager object
type AlertManagerReconciler struct {
	client.Client

	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	alertManager *grafana.AlertManager
}

//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=alertmanagers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=alertmanagers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=alertmanagers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AlertManager object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *AlertManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("reconciling")

	alertManagerManifest := &v1alpha1.AlertManager{}
	if err := r.Get(ctx, req.NamespacedName, alertManagerManifest); err != nil {
		logger.Error(err, "unable to fetch AlertManager")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if alertManagerManifest.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(alertManagerManifest.GetFinalizers(), alertManagerFinalizerName) {
			controllerutil.AddFinalizer(alertManagerManifest, alertManagerFinalizerName)
			if err := r.Update(ctx, alertManagerManifest); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		logger.Info("resetting AlertManager config")

		// The object is being deleted
		if containsString(alertManagerManifest.GetFinalizers(), alertManagerFinalizerName) {
			logger.Info("finalizer found, deleting AlertManager config from grafana")

			// our finalizer is present, so lets handle any external dependency
			if err := r.alertManager.Reset(ctx); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(alertManagerManifest, alertManagerFinalizerName)
			if err := r.Update(ctx, alertManagerManifest); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// handle actual reconciliation
	return r.doReconcileManifest(ctx, alertManagerManifest)
}

func (r *AlertManagerReconciler) doReconcileManifest(ctx context.Context, manifest *v1alpha1.AlertManager) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if err := r.alertManager.Configure(ctx, *manifest); err != nil {
		logger.Info("failed reconciling AlertManager")

		r.updateStatus(ctx, manifest, err)
		r.Recorder.Event(manifest, "Warning", "Error", "could not reconcile AlertManager with Grafana")

		return ctrl.Result{}, err
	}

	r.updateStatus(ctx, manifest, nil)
	r.Recorder.Event(manifest, "Normal", "Synchronized", "AlertManager reconciled")

	return ctrl.Result{}, nil
}

func StartAlertManagerReconciler(logger logr.Logger, ctrlManager ctrl.Manager, grabanaClient *grabana.Client) error {
	reconciler := &AlertManagerReconciler{
		Client:       ctrlManager.GetClient(),
		Scheme:       ctrlManager.GetScheme(),
		Recorder:     ctrlManager.GetEventRecorderFor("alertmanager-controller"),
		alertManager: grafana.NewAlertManager(logger, grabanaClient),
	}

	return reconciler.SetupWithManager(ctrlManager)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AlertManager{}).
		Complete(r)
}

func (r *AlertManagerReconciler) updateStatus(ctx context.Context, manifest *v1alpha1.AlertManager, err error) {
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
		logger.Error(err, "unable to update AlertManager status")
	}
}
