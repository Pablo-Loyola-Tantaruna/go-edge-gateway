package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Backends []BackendConfig `yaml:"backends"`
}

type ServerConfig struct {
	Port                int    `yaml:"port"`
	HealthCheckInterval string `yaml:"health_check_interval"`
}

type BackendConfig struct {
	URL string `yaml:"url"`
}

func LoadConfig(path string) (*Config, error) {
	conf := &Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, conf)
	return conf, err
}
