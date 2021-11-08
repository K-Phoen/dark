package worker

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/K-Phoen/dark/internal"
	"github.com/K-Phoen/dark/internal/pkg/dashboards"
	clientset "github.com/K-Phoen/dark/internal/pkg/generated/clientset/versioned"
	informers "github.com/K-Phoen/dark/internal/pkg/generated/informers/externalversions"
	"github.com/K-Phoen/grabana"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const threadiness = 2

type Config struct {
	MasterURL  string `env:"K8S_MASTER_URL"`
	KubeConfig string `env:"K8S_CONFIG"`

	GrafanaHost  string `env:"GRAFANA_HOST" validate:"required"`
	GrafanaToken string `env:"GRAFANA_TOKEN" validate:"required"`

	InsecureSkipVerify bool `env:"INSECURE_SKIP_VERIFY"`
}

type Worker struct {
	config          Config
	stop            chan struct{}
	controller      *internal.Controller
	informerFactory informers.SharedInformerFactory
}

func New(config Config) *Worker {
	return &Worker{
		config: config,
	}
}

func (worker *Worker) Init(logger *zap.Logger) error {
	fmt.Printf("config %#v\n", worker.config)
	restCfg, err := clientcmd.BuildConfigFromFlags(worker.config.MasterURL, worker.config.KubeConfig)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("error building kubernetes clientset: %w", err)
	}

	darkClient, err := clientset.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("error building dark clientset: %w", err)
	}

	grabanaClient := grabana.NewClient(makeHTTPClient(worker.config), worker.config.GrafanaHost, grabana.WithAPIToken(worker.config.GrafanaToken))
	dashboardCreator := dashboards.NewCreator(grabanaClient)

	worker.stop = make(chan struct{})
	worker.informerFactory = informers.NewSharedInformerFactory(darkClient, time.Second*30)
	worker.controller = internal.NewController(logger, kubeClient, darkClient, worker.informerFactory.Controller().V1().GrafanaDashboards(), dashboardCreator)

	return nil
}

func (worker *Worker) Run() error {
	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	worker.informerFactory.Start(worker.stop)

	return worker.controller.Run(threadiness, worker.stop)
}

func (worker *Worker) Terminate() error {
	close(worker.stop)

	return nil
}

func makeHTTPClient(cfg Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify},
		},
		Timeout: 10 * time.Second, // Large, but better than no timeout.
	}
}
