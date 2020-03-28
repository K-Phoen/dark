package main

import (
	"fmt"
	"os"

	"github.com/K-Phoen/dark/internal/pkg/converter"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create logger: %s", err)
	}

	if len(os.Args) != 3 {
		logger.Fatal("Usage: go run -mod=vendor main.go file.json output.yaml\n")
	}

	input, err := os.Open(os.Args[1])
	if err != nil {
		logger.Fatal("Could not open input file", zap.Error(err))
	}

	output, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		logger.Fatal("Could not open output file", zap.Error(err))
	}

	conv := converter.NewJSON(logger)

	if err := conv.Convert(input, output); err != nil {
		logger.Fatal("Could not convert dashboard", zap.Error(err))
	}
}
