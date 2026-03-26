package client

import (
	"log/slog"
	"os"
	"time"
)

// Config controls timeout, retry and logging behavior.
type Config struct {
	Timeout     time.Duration
	RetryMax    int
	BaseBackoff time.Duration
	Logger      *slog.Logger
}

func defaultConfig() Config {
	return Config{
		Timeout:     5 * time.Second,
		RetryMax:    3,
		BaseBackoff: 200 * time.Millisecond,
		Logger:      slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}
