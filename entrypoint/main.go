package main

import (
	"log"

	"github.com/pulzeai-oss/tritonserver/entrypoint/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
