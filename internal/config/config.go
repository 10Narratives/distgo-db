package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type BaseConfig struct {
	Logging LoggingConfig `yaml:"logging"`
	GRPC    GRPCConfig    `yaml:"grpc"`
}

type WorkerConfig struct {
	BaseConfig `yaml:",inline"`
}

type MasterConfig struct {
	BaseConfig `yaml:",inline"`
}

type LoggingConfig struct {
	Level    string `yaml:"level" env-default:"info"`
	Format   string `yaml:"format" env-default:"json"`
	Output   string `yaml:"output" env-default:"stdout"`
	FilePath string `yaml:"file_path" env-default:"./logs/worker.log"`
	MaxSize  int    `yaml:"max_size" env-default:"100"`
	MaxAge   int    `yaml:"max_age" env-default:"7"`
	Compress bool   `yaml:"compress" env-default:"true"`
}

type GRPCConfig struct {
	BindAddress string        `yaml:"bind_address" env-default:"0.0.0.0:50052"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
}

func mustConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to configuration file")
	flag.Parse()

	if configPath == "" {
		panic("cannot run node: path to configuration is missing")
	}

	return configPath
}

func mustLoadConfig(path string, cfg any) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file not found: %s", path))
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}
}

func MustLoadWorker() *WorkerConfig {
	var configPath string = mustConfigPath()
	var cfg WorkerConfig
	mustLoadConfig(configPath, &cfg)
	return &cfg
}

func MustLoadMaster() *MasterConfig {
	var configPath string = mustConfigPath()
	var cfg MasterConfig
	mustLoadConfig(configPath, &cfg)
	return &cfg
}
