package env

import (
	"github.com/caarlos0/env/v11"
	"net"
	"time"
)

type orderHTTPEnvConfig struct {
	Host        string        `env:"HTTP_HOST,required"`
	Port        string        `env:"HTTP_PORT,required"`
	ReadTimeout time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout time.Duration   `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
}

type orderHTTPConfig struct {
	raw orderHTTPEnvConfig
}

func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderHTTPConfig{raw: raw}, nil
}

func (cfg *orderHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *orderHTTPConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}

func (cfg *orderHTTPConfig) WriteTimeout() time.Duration {
	return cfg.raw.WriteTimeout
}

func (cfg *orderHTTPConfig) IdleTimeout() time.Duration {
	return cfg.raw.IdleTimeout
}
