package main

import (
	"fmt"

	"github.com/K-Phoen/dark/internal/pkg/worker"
	"github.com/voi-oss/svc"
	"go.uber.org/zap"

	// enables GCP auth
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// serviceVersion will be populated by the build script with the sha of the last git commit.
var serviceVersion = "snapshot"

func main() {
	cfg := worker.Config{}

	// Read up global configs
	if err := svc.LoadFromEnv(&cfg); err != nil {
		panic(fmt.Sprintf("could not load configuration: %s", err))
	}

	// SVC supervisor Init
	options := []svc.Option{
		svc.WithMetrics(),
		svc.WithHealthz(),
		svc.WithMetricsHandler(),
		svc.WithHTTPServer("9090"),
		svc.WithStackdriverLogger(zap.InfoLevel),
	}

	// SVC supervisor Init
	service, err := svc.New("dark", serviceVersion, options...)
	svc.MustInit(service, err)

	// Workers definition
	service.AddWorker("dashboards-controller", worker.New(cfg))

	// Service main loop
	service.Run()
}
