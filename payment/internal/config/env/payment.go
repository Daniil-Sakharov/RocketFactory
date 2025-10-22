package env

import (
	"github.com/caarlos0/env/v11"
	"net"
)

type paymentEnvConfig struct{
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type paymentConfig struct{
	raw paymentEnvConfig
}

func NewPaymentConfig() (*paymentConfig, error) {
	var raw paymentEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentConfig{raw: raw}, nil
}

func (p *paymentConfig)Address() string {
	return net.JoinHostPort(p.raw.Host,p.raw.Port)
}
