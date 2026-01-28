package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local" env:"ENV"`
	DatabaseURL string        `yaml:"database_url" env-required:"true" env:"DATABASE_URL"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h" env:"TOKEN_TTL"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080" env:"HTTP_SERVER_ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s" env:"HTTP_SERVER_TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s" env:"HTTP_SERVER_IDLE_TIMEOUT"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("config file %s does not exist, reading from env", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Printf("cannot read config from file: %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
