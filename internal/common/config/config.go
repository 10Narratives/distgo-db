package config

import "time"

// GRPCConfig holds configuration settings for a gRPC server.
type GRPCConfig struct {
	// Port specifies the network port for the gRPC server.
	Port int `yaml:"port"`
	// Timeout defines the maximum duration for gRPC operations.
	Timeout time.Duration `yaml:"timeout"`
}

// LoggingConfig holds configuration options for logging behavior.
type LoggingConfig struct {
	// Level sets the minimum log level (e.g., "debug", "info").
	Level string `yaml:"level"`
	// Format specifies the log output format (e.g., "json", "text").
	Format string `yaml:"format"`
	// File defines the path to the log output file.
	File string `yaml:"file"`
}
