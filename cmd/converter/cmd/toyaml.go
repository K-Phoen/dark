package cmd

import (
	"os"

	"github.com/K-Phoen/dark/internal/pkg/converter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func ToYamlCommand(logger *zap.Logger) *cobra.Command {
	var inputFile, outputFile string

	var cmd = &cobra.Command{
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

			conv := converter.NewJSON(logger)
			if err := conv.ToYAML(input, output); err != nil {
				logger.Fatal("Could not convert dashboard", zap.Error(err))
			}
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "input file")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagFilename("input")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "input file")
	_ = cmd.MarkFlagRequired("output")
	_ = cmd.MarkFlagFilename("output")

	return cmd
}
