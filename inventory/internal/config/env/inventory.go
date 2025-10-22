package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryConfig struct {
	raw inventoryEnvConfig
}

func NewInventoryConfig() (*inventoryConfig, error) {
	var raw inventoryEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryConfig{raw: raw}, nil
}

func (cfg *inventoryConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
