package cmd

import (
	"log"

	"github.com/pulzeai-oss/tensorrt-llm-deployment/entrypoint/internal/config"
	"github.com/pulzeai-oss/tensorrt-llm-deployment/entrypoint/internal/exec"
	"github.com/spf13/cobra"
)

func Execute() error {
	opts := config.Opts{}

	rootCmd := &cobra.Command{
		Use:   "entrypoint [flags]",
		Short: "An entrypoint for the TensorRT LLM deployment",
		Run: func(cmd *cobra.Command, args []string) {
			if err := exec.Run(&opts); err != nil {
				log.Fatalf("failed to execute entrypoint: %v", err)
			}
		},
	}
	rootCmd.Flags().IntVar(&opts.GrpcPort, "grpc-port", 8001, "Port to bind GRPC server to")
	rootCmd.Flags().IntVar(&opts.HttpPort, "http-port", 8000, "Port to bind HTTP server to")
	rootCmd.Flags().
		IntVar(&opts.MetricsPort, "metrics-port", 8002, "Port to bind metrics server to")
	rootCmd.Flags().
		StringVar(&opts.ModelRepo, "model-repo", "/srv/run/repo", "Path to Triton model repository")
	rootCmd.Flags().IntVar(&opts.WorldSize, "world-size", 1, "Number of GPUs to use for inference")

	return rootCmd.Execute()
}
