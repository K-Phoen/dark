package main

import (
	"crypto/tls"
	"flag"
	"net/http"
	"os"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/K-Phoen/grabana"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	k8skevingomezfrv1 "github.com/K-Phoen/dark/api/v1"
	k8skevingomezfrv1alpha1 "github.com/K-Phoen/dark/api/v1alpha1"
	"github.com/K-Phoen/dark/internal/pkg/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(k8skevingomezfrv1.AddToScheme(scheme))
	utilruntime.Must(k8skevingomezfrv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	// config definition
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var grafanaHost string
	var grafanaToken string
	var insecureSkipVerify bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller operator. "+
			"Enabling this will ensure there is only one active controller operator.")
	flag.StringVar(&grafanaHost, "grafana-host", "http://localhost:3000", "The host to use to reach Grafana.")
	flag.StringVar(&grafanaToken, "grafana-api-key", "", "The API key to use to authenticate to Grafana.")
	flag.BoolVar(&insecureSkipVerify, "insecure-skip-verify", false, "Skips SSL certificates verification. Useful when self-signed certificates are used, but can be insecure. Enabled at your own risks.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)

	must(viper.BindEnv("grafana-host", "GRAFANA_HOST"))
	must(viper.BindEnv("grafana-token", "GRAFANA_TOKEN"))
	must(viper.BindEnv("insecure-skip-verify", "INSECURE_SKIP_VERIFY"))

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	must(viper.BindPFlags(pflag.CommandLine))

	logger := zap.New(zap.UseFlagOptions(&opts))
	ctrl.SetLogger(logger)

	httpClient := makeHTTPClient(&tls.Config{
		//nolint:gosec
		InsecureSkipVerify: viper.GetBool("insecure-skip-verify"),
	})
	grabanaClient := grabana.NewClient(
		httpClient,
		viper.GetString("grafana-host"),
		grabana.WithAPIToken(viper.GetString("grafana-token")),
	)

	// controllers setup
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "2351aaad.k8s.kevingomez.fr",
	})
	if err != nil {
		setupLog.Error(err, "unable to start operator")
		os.Exit(1)
	}

	if err = controllers.StartGrafanaDashboardReconciler(mgr, grabanaClient); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "GrafanaDashboard")
		os.Exit(1)
	}
	if err = controllers.StartDatasourceReconciler(logger, mgr, grabanaClient); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Datasource")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	// liveness and readiness probes
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// main runtime loop
	setupLog.Info("starting operator")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running operator")
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		setupLog.Error(err, "")
		os.Exit(1)
	}
}

func makeHTTPClient(tlsConfig *tls.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 10 * time.Second, // Large, but better than no timeout.
	}
}
