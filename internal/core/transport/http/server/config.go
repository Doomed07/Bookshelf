package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr              string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout   time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	ReadHeaderTimeout time.Duration `envconfig:"READ_HEADER_TIMEOUT" default:"5s"`
	ReadTimeout       time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
	WriteTimeout      time.Duration `envconfig:"WRITE_TIMEOUT" default:"15s"`
	IdleTimeout       time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get HTTP server config: %w", err)
		panic(err)
	}

	return config
}

type metricsEnv struct {
	Addr string `envconfig:"ADDR" default:":2112"`
}

// NewMetricsConfigMust — конфиг сервера метрик (METRICS_ADDR). Возвращает обычный Config,
// чтобы для /metrics переиспользовать тот же HTTPServer с корректной остановкой.
func NewMetricsConfigMust() Config {
	var env metricsEnv

	if err := envconfig.Process("METRICS", &env); err != nil {
		panic(fmt.Errorf("get metrics server config: %w", err))
	}

	return Config{
		Addr:              env.Addr,
		ShutdownTimeout:   5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
