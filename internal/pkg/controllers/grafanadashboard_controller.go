package controllers

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/controller"

	k8skevingomezfrv1 "github.com/K-Phoen/dark/api/v1"
	"github.com/K-Phoen/dark/internal/pkg/grafana"
	"github.com/K-Phoen/grabana"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const grafanaDashboardFinalizerName = "grafanadashboards.k8s.kevingomez.fr/finalizer"
const folderAnnotation = "dark/folder"

type dashboardManager interface {
	FromRawSpec(ctx context.Context, folderName string, uid string, rawJSON []byte) error
	Delete(ctx context.Context, uid string) error
}

// GrafanaDashboardReconciler reconciles a GrafanaDashboard object
type GrafanaDashboardReconciler struct {
	client.Client

	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	Dashboards dashboardManager
}

func StartGrafanaDashboardReconciler(ctrlManager ctrl.Manager, grabanaClient *grabana.Client) error {
	reconciler := &GrafanaDashboardReconciler{
		Client:     ctrlManager.GetClient(),
		Scheme:     ctrlManager.GetScheme(),
		Recorder:   ctrlManager.GetEventRecorderFor("grafanadashboard-controller"),
		Dashboards: grafana.NewCreator(grabanaClient),
	}

	return reconciler.SetupWithManager(ctrlManager)
}

//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=grafanadashboards,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=grafanadashboards/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.kevingomez.fr,resources=grafanadashboards/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *GrafanaDashboardReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("reconciling")

	dashboard := &k8skevingomezfrv1.GrafanaDashboard{}
	if err := r.Get(ctx, req.NamespacedName, dashboard); err != nil {
		logger.Error(err, "unable to fetch GrafanaDashboard")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// examine DeletionTimestamp to determine if object is under deletion
	if dashboard.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(dashboard.GetFinalizers(), grafanaDashboardFinalizerName) {
			controllerutil.AddFinalizer(dashboard, grafanaDashboardFinalizerName)
			if err := r.Update(ctx, dashboard); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		logger.Info("deleting GrafanaDashboard")

		// The object is being deleted
		if containsString(dashboard.GetFinalizers(), grafanaDashboardFinalizerName) {
			logger.Info("finalizer found, deleting dashboard from grafana")

			// our finalizer is present, so lets handle any external dependency
			if err := r.Dashboards.Delete(ctx, dashboard.Name); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(dashboard, grafanaDashboardFinalizerName)
			if err := r.Update(ctx, dashboard); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	folder := dashboard.Annotations[folderAnnotation]
	if dashboard.Folder != "" {
		folder = dashboard.Folder
	}

	// proceed with create/update reconciliation
	if err := r.Dashboards.FromRawSpec(ctx, folder, dashboard.ObjectMeta.Name, dashboard.Spec.Raw); err != nil {
		logger.Error(err, "could not apply GrafanaDashboard in Grafana")

		r.updateStatus(ctx, dashboard, err)
		r.Recorder.Event(dashboard, "Warning", "Error", "could not apply GrafanaDashboard in Grafana")

		return ctrl.Result{}, err
	}

	logger.Info("done!")

	r.updateStatus(ctx, dashboard, nil)
	r.Recorder.Event(dashboard, "Normal", "Synchronized", "GrafanaDashboard synchronized")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GrafanaDashboardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8skevingomezfrv1.GrafanaDashboard{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 5,
		}).
		Complete(r)
}

func (r *GrafanaDashboardReconciler) updateStatus(ctx context.Context, dashboard *k8skevingomezfrv1.GrafanaDashboard, err error) {
	logger := log.FromContext(ctx)

	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	dashboardCopy := dashboard.DeepCopy()

	if err == nil {
		dashboardCopy.Status.Status = "OK"
		dashboardCopy.Status.Message = "Synchronized"
	} else {
		dashboardCopy.Status.Status = "Error"
		dashboardCopy.Status.Message = err.Error()
	}

	if err := r.Status().Update(ctx, dashboardCopy); err != nil {
		logger.Error(err, "unable to update GrafanaDashboard status")
	}
}
