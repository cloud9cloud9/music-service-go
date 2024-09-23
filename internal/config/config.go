package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

const (
	cfgPath = "config/config.yml"
)

type Config struct {
	MySql struct {
		DBUser     string `yaml:"username"`
		DBPassword string `yaml:"password" env:"MYSQL_ROOT_PASSWORD"`
		DBAddress  string `yaml:"addr"`
		DBName     string `yaml:"db_name" env:"MYSQL_DATABASE"`
		DBNet      string `yaml:"net"`
	} `yaml:"mysql"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	JWT struct {
		Secret     string `yaml:"secret"`
		Expiration int64  `yaml:"expiration"`
	} `yaml:"jwt"`
	Spotify struct {
		ClientID     string `yaml:"client_id" env:"CLIENT_ID"`
		ClientSecret string `yaml:"client_secret" env:"CLIENT_SECRET"`
	} `yaml:"spotify"`
}

var Instance *Config

var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		Instance = &Config{}

		if err := cleanenv.ReadConfig(cfgPath, Instance); err != nil {
			help, _ := cleanenv.GetDescription(Instance, nil)
			log.Fatalf("Config error: %s", help)
		}
	})
	return Instance
}
