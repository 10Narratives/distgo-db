package main

import (
	"fmt"

	mastercfg "github.com/10Narratives/distgo-db/internal/config/master"
)

func main() {
	cfg := mastercfg.MustLoad()
	fmt.Println(cfg)
}
