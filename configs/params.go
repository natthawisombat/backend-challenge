package configs

import (
	"backend-challenge/pkg/logging"
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

var (
	App = new(config)
)

type config struct {
	Host    string        `env:"APP_HOST,default=localhost" json:",omitempty"`
	Port    string        `env:"APP_PORT,default=8080" json:",omitempty"`
	Timeout time.Duration `env:"APP_TIMEOUT,default=1m" json:",omitempty"`
	Prefix  string        `env:"APP_PREFIX,default=/develop" json:",omitempty"`
}

func SetEnv(ctx context.Context) error {
	logger := logging.FromContext(ctx).Named("set environment")
	configs := []interface{}{
		App,
		logging.L,
	}
	for _, cfg := range configs {
		if err := envconfig.Process(ctx, cfg); err != nil {
			logger.Warnw("failed to process config file: %v", err)
			return fmt.Errorf("failed to process config file: %v", err)
		}
	}

	return nil
}
