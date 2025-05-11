package config

import (
	"flag"
	"os"

	"githib.com/10Narratives/distgo-db/pkg/config"
	"github.com/ilyakaznacheev/cleanenv"
)

type MasterConfig struct {
	GRPC    config.GRPCConfig    `yaml:"grpc"`
	Logging config.LoggingConfig `yaml:"logging"`
}

func MustLoad() *MasterConfig {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to configuration file")
	flag.Parse()

	if configPath == "" {
		panic("cannot run worker node: path to configuration is missing")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("configuration file does not exist: " + configPath)
	}

	var cfg MasterConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read configuration: " + err.Error())
	}

	return &cfg
}
