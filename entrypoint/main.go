package main

import (
	"log"

	"github.com/pulzeai-oss/tensorrt-llm-deployment/entrypoint/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("failed to execute entrypoint: %v", err)
	}
}
