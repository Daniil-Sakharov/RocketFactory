package config

import (
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/config/env"
	"github.com/joho/godotenv"
	"os"
)

var appConfig *config

type config struct {
	Inventory InventoryConfig
	Logger    LoggerConfig
	Mongo     MongoConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	inventoryCfg, err := env.NewInventoryConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Inventory: inventoryCfg,
		Logger:    loggerCfg,
		Mongo:     mongoCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
