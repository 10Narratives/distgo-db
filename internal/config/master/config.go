package mastercfg

import (
	"flag"
	"fmt"
	"os"

	"github.com/10Narratives/distgo-db/internal/config"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPC    config.GRPCConfig
	Logging config.LoggingConfig
}

func MustLoad() *Config {
	var path string

	flag.StringVar(&path, "config", "", "path to configuration file")
	flag.Parse()

	if path == "" {
		panic("cannot run node: path to configuration is missing")
	}

	return MustLoadFromFile(path)
}

func MustLoadFromFile(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file not found: %s", path))
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}

	return &cfg
}
