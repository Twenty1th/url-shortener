package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" binding:"required"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Addr        string        `yaml:"address" binding:"required"`
	Timout      time.Duration `yaml:"timout" binding:"required"`
	IdleTimeout time.Duration `yaml:"idle_timeout" binding:"required"`
	User        string        `yaml:"user" binding:"required"`
	Password    string        `yaml:"password" binding:"required"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	if _, err := os.Stat(configPath); err != nil {
		log.Fatal("CONFIG_PATH does not exist")
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config: %s", err)
	}
	return &cfg
}
