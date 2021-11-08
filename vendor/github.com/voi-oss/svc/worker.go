package svc

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// Worker defines a SVC worker.
type Worker interface {
	Init(*zap.Logger) error
	Run() error
	Terminate() error
}

// Healther defines a worker that can report his healthz status.
type Healther interface {
	Healthy() error
}

// Gatherer is a place for workers to return a prometheus.Gatherer
// for SVC to serve on the metrics endpoint
type Gatherer interface {
	Gatherer() prometheus.Gatherer
}
