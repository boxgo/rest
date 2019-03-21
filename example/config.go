package main

import (
	"os"

	"github.com/boxgo/box/config"
	"github.com/boxgo/box/config/loader"
)

const (
	envFlag = "APP_MODE"
	envDev  = "dev"
	envTest = "test"
	envUat  = "uat"
	envProd = "prod"
)

func newConfig() config.Config {
	ld := loader.NewFileEnvConfig(configPath())
	config := config.NewConfig(ld)

	return config
}

func configPath() string {
	configFileName := envDev

	if name, ok := os.LookupEnv(envFlag); ok {
		configFileName = name
	}

	return "./config/" + configFileName + ".yaml"
}

func envMode() string {
	return os.Getenv(envFlag)
}
