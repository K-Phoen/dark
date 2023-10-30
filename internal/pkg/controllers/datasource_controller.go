package controllers

import (
	"context"

	"github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/dark/internal/pkg/grafana"
	"github.com/K-Phoen/dark/internal/pkg/kubernetes"
	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/datasource"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const datasourcesFinalizerName = "datasources.k8s.kevingomez.fr/finalizer"

type datasourcesManager interface {
	SpecToModel(ctx context.Context, objectRef types.NamespacedName, spec v1alpha1.DatasourceSpec) (datasource.Datasource, error)
	Upsert(ctx context.Context, model datasource.Datasource) error
	Delete(ctx context.Context, name string) error
}

// DatasourceReconciler reconciles a Datasource object
type DatasourceReconciler struct {
	client.Client

	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	Datasources datasourcesManager
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DatasourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("reconciling")

	datasourceManifest := &v1alpha1.Datasource{}
	if err := r.Get(ctx, req.NamespacedName, datasourceManifest); err != nil {
		logger.Error(err, "unable to fetch Datasource")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if datasourceManifest.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(datasourceManifest.GetFinalizers(), datasourcesFinalizerName) {
			controllerutil.AddFinalizer(datasourceManifest, datasourcesFinalizerName)
			if err := r.Update(ctx, datasourceManifest); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		logger.Info("deleting Datasource")

		// The object is being deleted
		if containsString(datasourceManifest.GetFinalizers(), datasourcesFinalizerName) {
			logger.Info("finalizer found, deleting datasource from grafana")

			// our finalizer is present, so lets handle any external dependency
			if err := r.Datasources.Delete(ctx, datasourceManifest.Name); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(datasourceManifest, datasourcesFinalizerName)
			if err := r.Update(ctx, datasourceManifest); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	datasourceModel, err := r.Datasources.SpecToModel(ctx, req.NamespacedName, datasourceManifest.Spec)
	if err != nil {
		logger.Error(err, "unable to convert Datasource manifest into a Grabana model")

		r.updateStatus(ctx, datasourceManifest, err)
		r.Recorder.Event(datasourceManifest, "Warning", "Error", "could not synchronize Datasource with Grafana")

		return ctrl.Result{}, err
	}

	// proceed with create/update reconciliation
	if err := r.Datasources.Upsert(ctx, datasourceModel); err != nil {
		logger.Error(err, "could not upsert Datasource in Grafana")

		r.updateStatus(ctx, datasourceManifest, err)
		r.Recorder.Event(datasourceManifest, "Warning", "Error", "could not synchronize Datasource with Grafana")

		return ctrl.Result{}, err
	}

	logger.Info("done!")

	r.updateStatus(ctx, datasourceManifest, nil)
	r.Recorder.Event(datasourceManifest, "Normal", "Synchronized", "Datasource synchronized")

	return ctrl.Result{}, nil
}

//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=datasources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=datasources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=datasources/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

func StartDatasourceReconciler(logger logr.Logger, ctrlManager ctrl.Manager, grabanaClient *grabana.Client) error {
	refReader := kubernetes.NewValueRefReader(logger, kubernetes.NewSecrets(logger, ctrlManager.GetClient()))

	reconciler := &DatasourceReconciler{
		Client:      ctrlManager.GetClient(),
		Scheme:      ctrlManager.GetScheme(),
		Recorder:    ctrlManager.GetEventRecorderFor("grafanadashboard-controller"),
		Datasources: grafana.NewDatasources(logger, grabanaClient, refReader),
	}

	return reconciler.SetupWithManager(ctrlManager)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatasourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Datasource{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func (r *DatasourceReconciler) updateStatus(ctx context.Context, manifest *v1alpha1.Datasource, err error) {
	logger := log.FromContext(ctx)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	manifestCopy := manifest.DeepCopy()

	if err == nil {
		manifestCopy.Status.Status = "OK"
		manifestCopy.Status.Message = "Synchronized"
	} else {
		manifestCopy.Status.Status = "Error"
		manifestCopy.Status.Message = err.Error()
	}

	if err := r.Status().Update(ctx, manifestCopy); err != nil {
		logger.Error(err, "unable to update Datasource status")
	}
}
