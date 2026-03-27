package client

import (
	"log/slog"
	"time"
)

// Config controls transport connection settings per client instance.
type Config struct {
	Host    string
	Port    int
	Timeout time.Duration
	Logger  *slog.Logger
}

func defaultConfig() Config {
	return Config{
		Host:    "192.168.100.1",
		Port:    9200,
		Timeout: 5 * time.Second,
	}
}
