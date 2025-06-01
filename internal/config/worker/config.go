package workercfg

import (
	"flag"
	"fmt"
	"os"

	"github.com/10Narratives/distgo-db/internal/config"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Name    string               `yaml:"name" env-required:"true"`
	GRPC    config.GRPCConfig    `yaml:"grpc"`
	Logging config.LoggingConfig `yaml:"logging"`
}

func MustLoad() *Config {
	var path string

	flag.StringVar(&path, "config", "", "path to configuration file")
	flag.Parse()

	if path == "" {
		panic("cannot load config: path to file is missing")
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
