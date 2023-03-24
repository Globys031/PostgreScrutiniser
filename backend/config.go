package main

import (
	"fmt"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/spf13/viper"
)

type Config struct {
	JWT_secret_key string `mapstructure:"JWT_SECRET_KEY"`
	Backend_port   int    `mapstructure:"BACKEND_PORT"`
}

func LoadConfig(logger *utils.Logger) (c Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		logger.LogFatal(fmt.Errorf("Could not load .env file configs: %v", err))
		return
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		logger.LogFatal(fmt.Errorf("Could not unmarshal .env file configs: %v", err))
	}

	return
}
