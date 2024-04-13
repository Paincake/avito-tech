package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env                      string  `env:"ENV" env_default:"local"`
	Port                     string  `env:"DB_PORT" env_default:"5432"`
	Host                     string  `env:"DB_HOST" env_default:"localhost"`
	Name                     string  `env:"DB_NAME" env_default:"postgres"`
	User                     string  `env:"DB_USER" env_default:"user"`
	Password                 string  `env:"DB_PASSWORD" env_default:"password"`
	CacheSchedulerRate       int64   `env:"SCHEDULER_RATE_MINUTE" env_default:"1"`
	CacheKeyInvalidationTime float64 `env:"CACHE_KEY_INVALIDATION_MINUTES" env_default:"5"`
}

func MustLoad(configPath string) (*Config, error) {
	if configPath == "" {
		return &Config{}, errors.New("not a valid config file path")
	}
	if _, err := os.Stat(configPath); err != nil {
		return &Config{}, errors.New("error reading config file")
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return &Config{}, errors.New("error reading config file")
	}
	return &cfg, nil
}
