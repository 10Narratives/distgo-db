package config

import "time"

type GRPCConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type LoggingConfig struct {
	Level  string `yaml:"level" env-default:"info"`
	Format string `yaml:"format" env-default:"json"`
	Output string `yaml:"output" env-default:"stdout"`
}
