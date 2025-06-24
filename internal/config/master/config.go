package mastercfg

import (
	"flag"
	"os"

	"github.com/10Narratives/distgo-db/internal/config"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPC    config.GRPCConfig    `yaml:"grpc"`
	Logging config.LoggingConfig `yaml:"logging"`
}

func MustLoad() *Config {
	var path string

	flag.StringVar(&path, "config", "", "path to configuration file")
	flag.Parse()

	return MustLoadFromFile(path)
}

func MustLoadFromFile(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
