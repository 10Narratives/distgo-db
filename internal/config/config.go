package config

import (
	"time"

	_ "github.com/ilyakaznacheev/cleanenv"
)

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"50052"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
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
