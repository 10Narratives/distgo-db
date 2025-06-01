package main

import (
	"fmt"

	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
)

func main() {
	cfg := workercfg.MustLoad()
	fmt.Println(cfg)

	// TODO: initialize logger
	// TODO: initialize gRPC application
	// TODO: initialize application
}
