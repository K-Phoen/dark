package main

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/K-Phoen/dark/internal"
	"github.com/K-Phoen/dark/internal/pkg/dashboards"
	"github.com/K-Phoen/grabana"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	// enables GCP auth
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	clientset "github.com/K-Phoen/dark/internal/pkg/generated/clientset/versioned"
	informers "github.com/K-Phoen/dark/internal/pkg/generated/informers/externalversions"
	"github.com/K-Phoen/dark/internal/pkg/signals"
	"github.com/caarlos0/env"
	"gopkg.in/go-playground/validator.v9"
)

type config struct {
	MasterURL  string `env:"K8S_MASTER_URL"`
	KubeConfig string `env:"K8S_CONFIG"`

	GrafanaHost  string `env:"GRAFANA_HOST" validate:"required"`
	GrafanaToken string `env:"GRAFANA_TOKEN" validate:"required"`

	InsecureSkipVerify bool `env:"INSECURE_SKIP_VERIFY"`
}

func (cfg *config) loadFromEnv() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}
	if err := validator.New().Struct(*cfg); err != nil {
		return err
	}

	return nil
}

func main() {
	cfg := config{}
	if err := cfg.loadFromEnv(); err != nil {
		klog.Fatalf("Error loading configuration: %s", err.Error())
	}

	restCfg, err := clientcmd.BuildConfigFromFlags(cfg.MasterURL, cfg.KubeConfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	darkClient, err := clientset.NewForConfig(restCfg)
	if err != nil {
		klog.Fatalf("Error building dark clientset: %s", err.Error())
	}

	darkInformerFactory := informers.NewSharedInformerFactory(darkClient, time.Second*30)

	grabanaClient := grabana.NewClient(makeHTTPClient(cfg), cfg.GrafanaHost, grabana.WithAPIToken(cfg.GrafanaToken))
	dashboardCreator := dashboards.NewCreator(grabanaClient)

	controller := internal.NewController(kubeClient, darkClient, darkInformerFactory.Controller().V1().GrafanaDashboards(), dashboardCreator)

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	darkInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}

func makeHTTPClient(cfg config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify},
		},
		Timeout: 10 * time.Second, // Large, but better than no timeout.
	}
}
