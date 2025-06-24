package config

import "time"

type GRPCConfig struct {
	Host    string        `yaml:"host" env-required:"true"`
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type LoggingConfig struct {
	Level  string `yaml:"level" env-default:"error"`
	Format string `yaml:"format" env-default:"json"`
	Output string `yaml:"output" env-default:"file"`
}
