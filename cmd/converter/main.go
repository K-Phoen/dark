package main

import (
	"fmt"
	"os"

	"github.com/K-Phoen/dark/internal/pkg/converter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	var inputFile, outputFile string
	var folder string
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create logger: %s", err)
		os.Exit(1)
	}

	conv := converter.NewJSON(logger)

	var cmdYaml = &cobra.Command{
		Use:   "convert-yaml",
		Short: "Converts a JSON dashboard into YAML",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := os.Open(inputFile)
			if err != nil {
				logger.Fatal("Could not open input file", zap.Error(err))
			}

			output, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
			if err != nil {
				logger.Fatal("Could not open output file", zap.Error(err))
			}

			if err := conv.ToYAML(input, output); err != nil {
				logger.Fatal("Could not convert dashboard", zap.Error(err))
			}
		},
	}
	var cmdK8sManifest = &cobra.Command{
		Use:   "convert-k8s-manifest",
		Args:  cobra.ExactArgs(1),
		Short: "Converts a JSON dashboard into a k8s manifest",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := os.Open(inputFile)
			if err != nil {
				logger.Fatal("Could not open input file", zap.Error(err))
			}

			output, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
			if err != nil {
				logger.Fatal("Could not open output file", zap.Error(err))
			}

			if err := conv.ToK8SManifest(input, output, folder, args[0]); err != nil {
				logger.Fatal("Could not convert dashboard", zap.Error(err))
			}
		},
	}

	cmdYaml.Flags().StringVarP(&inputFile, "input", "i", "", "input file")
	cmdYaml.MarkFlagRequired("input")
	cmdYaml.MarkFlagFilename("input")
	cmdYaml.Flags().StringVarP(&outputFile, "output", "o", "", "input file")
	cmdYaml.MarkFlagRequired("output")
	cmdYaml.MarkFlagFilename("output")

	cmdK8sManifest.Flags().StringVarP(&inputFile, "input", "i", "", "input file")
	cmdK8sManifest.MarkFlagRequired("input")
	cmdK8sManifest.MarkFlagFilename("input")
	cmdK8sManifest.Flags().StringVarP(&outputFile, "output", "o", "", "input file")
	cmdK8sManifest.MarkFlagRequired("output")
	cmdK8sManifest.MarkFlagFilename("output")
	cmdK8sManifest.Flags().StringVar(&folder, "folder", "General", "dashboard folder")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdYaml)
	rootCmd.AddCommand(cmdK8sManifest)
	rootCmd.Execute()
}
