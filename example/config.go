package main

import (
	"os"

	"github.com/boxgo/box/config"
	"github.com/boxgo/config/source/env"
	"github.com/boxgo/config/source/file"
)

const (
	envFlag = "ENV"
	envDev  = "dev"
	envTest = "test"
	envUat  = "uat"
	envProd = "prod"
)

func newConfig() config.Config {
	cfg := config.NewConfig(
		file.NewSource(file.WithPath(configPath())),
		env.NewSource(),
	)

	return cfg
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
