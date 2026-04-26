package config

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
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
	if err != nil {
		return nil, err
	}

	if envSecret := os.Getenv("GOPHER_JWT_SECRET"); envSecret != "" {
		log.Println("Secreto JWT cargado desde el entorno")
	}

	return conf, err
}
func WatchConfig(path string, onChange func(*Config)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("Detectado cambio en config.yaml, recargando...")
					newCfg, err := LoadConfig(path)
					if err == nil {
						onChange(newCfg)
					} else {
						log.Printf("Error al recargar configuración: %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
}
