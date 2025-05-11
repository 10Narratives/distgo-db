package main

import (
	"fmt"

	"githib.com/10Narratives/distgo-db/internal/master/config"
)

func main() {
	// TODO: read configuration file
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: initialize logger

	// TODO: initialize gRPC-server app

	// TODO: initialize master app
}
