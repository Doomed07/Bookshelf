package statistics_service

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MinRatings int `envconfig:"MIN_RATINGS" default:"30"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("STATS", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	if config.MinRatings < 1 {
		return Config{}, fmt.Errorf("STATS_MIN_RATINGS must be >= 1, got %d", config.MinRatings)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get statistics config: %w", err)
		panic(err)
	}

	return config
}
