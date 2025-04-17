package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	Domain        string `env:"SITE_DOMAIN"`
	UserId        string `env:"ADMIN_USER_ID"`
}

func Get() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return cfg, err
}
