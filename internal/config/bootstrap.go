package config

import (
	"os"
)

type Bootstrap struct {
	AppName     string
	Environment string
	ConfigPath  string
}

func InitBootstrap() *Bootstrap {
	env := os.Getenv("APP_PROFILE")
	if env == "" {
		env = "dev"
	}

	return &Bootstrap{
		AppName:     "gopher-guard",
		Environment: env,
		ConfigPath:  "vars/" + env + ".yaml",
	}
}
