package config

import (
    "github.com/Daniil-Sakharov/RocketFactory/payment/internal/config/env"
    "github.com/joho/godotenv"
)

var appConfig *config

type config struct {
    Payment PaymentConfig
    Logger  LoggerConfig
}

func Load(path ...string) error {
    if err := godotenv.Load(path...); err != nil {
        return err
    }

    paymentCfg, err := env.NewPaymentConfig()
    if err != nil {
        return err
    }

    loggerCfg, err := env.NewLoggerConfig()
    if err != nil {
        return err
    }

    appConfig = &config{
        Payment: paymentCfg,
        Logger:  loggerCfg,
    }

    return nil
}

func AppConfig() *config {
    return appConfig
}