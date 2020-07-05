package main

import (
	"fmt"
	"os"

	"github.com/K-Phoen/dark/cmd/converter/cmd"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger, err := createLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create logger: %s", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmd.ToYamlCommand(logger))
	rootCmd.AddCommand(cmd.ToManifestCommand(logger))

	_ = rootCmd.Execute()
}

func createLogger() (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}
