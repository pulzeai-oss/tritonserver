package exec

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/pulzeai-oss/tensorrt-llm-deployment/entrypoint/internal/config"
)

const MpiRunExecutable = "/usr/local/mpi/bin/mpirun"
const TritonServerExecutable = "/opt/tritonserver/bin/tritonserver"

func exec(argv []string) error {
	prettyArgs, err := json.Marshal(argv[1:])
	if err != nil {
		return err
	}
	log.Printf("Executing command %q with args: %s", argv[0], prettyArgs)
	return syscall.Exec(argv[0], argv, os.Environ())
}

func Run(opts *config.Opts) error {
	// Read configuration from environment variables and fill repo templates
	if err := config.FillConfigTemplatesFromEnv(opts.ModelRepo); err != nil {
		return fmt.Errorf("failed to fill config templates from environment variables: %v", err)
	}

	// Build args
	args := []string{MpiRunExecutable, "--allow-run-as-root"}
	for i := 0; i < opts.WorldSize; i++ {
		workerArgs := []string{
			"-n", "1",
			TritonServerExecutable,
			fmt.Sprintf("--backend-config=python,shm-region-prefix-name=prefix%d_", i),
			"--disable-auto-complete-config",
			fmt.Sprintf("--model-repository=%s", opts.ModelRepo),
		}

		// Additional args for first worker
		if i == 0 {
			workerArgs = append(workerArgs,
				fmt.Sprintf("--grpc-port=%d", opts.GrpcPort),
				fmt.Sprintf("--http-port=%d", opts.HttpPort),
				"--log-file=/dev/stdout",
				fmt.Sprintf("--log-verbose=%d", opts.Verbosity),
				fmt.Sprintf("--metrics-port=%d", opts.MetricsPort),
			)
		} else {
			workerArgs = append(workerArgs,
				fmt.Sprintf("--load-model=%s", config.DIR_ENGINE_IDENTIFIER),
				"--model-control-mode=explicit",
			)
		}
		args = append(args, append(workerArgs, ":")...)
	}

	// Run Triton inference server
	return exec(args)
}
