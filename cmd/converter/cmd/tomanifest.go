package cmd

import (
	"os"

	"github.com/K-Phoen/dark/internal/pkg/converter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func ToManifestCommand(logger *zap.Logger) *cobra.Command {
	var inputFile, outputFile string
	var options converter.K8SManifestOptions

	var cmd = &cobra.Command{
		Use:   "convert-k8s-manifest name",
		Args:  cobra.ExactArgs(1),
		Short: "Converts a JSON dashboard into a k8s manifest",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := os.Open(inputFile)
			if err != nil {
				logger.Fatal("Could not open input file", zap.Error(err))
			}

			output, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				logger.Fatal("Could not open output file", zap.Error(err))
			}

			options.Name = args[0]

			conv := converter.NewJSON(logger)
			if err := conv.ToK8SManifest(input, output, options); err != nil {
				logger.Fatal("Could not convert dashboard", zap.Error(err))
			}
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagFilename("input")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Input file")
	_ = cmd.MarkFlagRequired("output")
	_ = cmd.MarkFlagFilename("output")
	cmd.Flags().StringVar(&options.Folder, "folder", "Dark", "Dashboard folder")
	cmd.Flags().StringVarP(&options.Namespace, "namespace", "n", "", "Manifest namespace")

	return cmd
}
