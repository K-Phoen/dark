package svc

import (
	"os"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (s *SVC) newLogger(level zapcore.Level, encoder zapcore.Encoder) (*zap.Logger, zap.AtomicLevel) {
	atom := zap.NewAtomicLevel()
	atom.SetLevel(level)

	s.zapOpts = append(s.zapOpts, zap.ErrorOutput(zapcore.Lock(os.Stderr)), zap.AddCaller())

	logger := zap.New(zapcore.NewSamplerWithOptions(zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stdout),
		atom,
	), time.Second, 100, 10),
		s.zapOpts...,
	)

	return logger, atom
}

// WithZapMetrics will add a hook to the zap logger and emit metrics to prometheus
// based on log level and log name.
func WithZapMetrics() Option {
	return func(s *SVC) error {
		requestCount := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "logger_emitted_entries",
				Help: "Number of log messages emitted.",
			},
			[]string{"level", "logger_name"},
		)
		if err := prometheus.Register(requestCount); err != nil {
			return err
		}

		s.zapOpts = append(s.zapOpts,
			zap.Hooks(func(e zapcore.Entry) error {
				counter, err := requestCount.GetMetricWithLabelValues(e.Level.String(), e.LoggerName)
				if err != nil {
					return err
				}
				counter.Inc()
				return nil
			}))
		return nil
	}
}

// WithLogger is an option that allows you to provide your own customized logger.
func WithLogger(logger *zap.Logger, atom zap.AtomicLevel) Option {
	return func(s *SVC) error {
		return assignLogger(s, logger, atom)
	}
}

// WithDevelopmentLogger is an option that uses a zap Logger with
// configurations set meant to be used for development.
func WithDevelopmentLogger(opts ...zap.Option) Option {
	return func(s *SVC) error {
		s.zapOpts = append(s.zapOpts, opts...)
		logger, atom := s.newLogger(
			zapcore.DebugLevel,
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		)
		logger = logger.With(zap.String("app", s.Name), zap.String("version", s.Version))
		return assignLogger(s, logger, atom)
	}
}

// WithProductionLogger is an option that uses a zap Logger with configurations
// set meant to be used for production.
func WithProductionLogger(opts ...zap.Option) Option {
	return func(s *SVC) error {
		s.zapOpts = append(s.zapOpts, opts...)
		logger, atom := s.newLogger(
			zapcore.InfoLevel,
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		)
		logger = logger.With(zap.String("app", s.Name), zap.String("version", s.Version))
		return assignLogger(s, logger, atom)
	}
}

// WithConsoleLogger is an option that uses a zap Logger with configurations
// set meant to be used for debugging in the console.
func WithConsoleLogger(level zapcore.Level, opts ...zap.Option) Option {
	return func(s *SVC) error {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.RFC3339TimeEncoder
		s.zapOpts = append(s.zapOpts, opts...)

		logger, atom := s.newLogger(
			level,
			zapcore.NewConsoleEncoder(config),
		)
		return assignLogger(s, logger, atom)
	}
}

// WithStackdriverLogger is an option that uses a zap Logger with configurations
// set meant to be used for production and is compliant with the GCP/Stackdriver format.
func WithStackdriverLogger(level zapcore.Level, opts ...zap.Option) Option {
	return func(s *SVC) error {
		s.zapOpts = append(s.zapOpts, opts...)
		logger, atom := s.newLogger(
			level,
			zapcore.NewJSONEncoder(zapdriver.NewProductionEncoderConfig()),
		)
		logger = logger.With(zapdriver.ServiceContext(s.Name), zapdriver.Label("version", s.Version))
		return assignLogger(s, logger, atom)
	}
}

func assignLogger(s *SVC, logger *zap.Logger, atom zap.AtomicLevel) error {
	stdLogger, err := zap.NewStdLogAt(logger, zapcore.ErrorLevel)
	if err != nil {
		return err
	}
	undo, err := zap.RedirectStdLogAt(logger, zapcore.ErrorLevel)
	if err != nil {
		return err
	}

	s.logger = logger
	s.stdLogger = stdLogger
	s.atom = atom
	s.loggerRedirectUndo = undo

	return nil
}
