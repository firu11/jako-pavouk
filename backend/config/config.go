package config

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type EmailConfig struct {
	To              string `env:"EMAIL_TO"`
	From            string `env:"EMAIL_FROM"`
	Password        string `env:"EMAIL_PASSWORD"`
	Host            string `env:"EMAIL_HOST"`
	Port            int    `env:"EMAIL_PORT,default=465"`
	NotificationURL string `env:"MOBIL_NOTIFIKACE_URL"`
}

func (c EmailConfig) Enabled() bool {
	return c.Host != "" && c.From != "" && c.Password != ""
}

type Config struct {
	Production  bool   `env:"PRODUCTION,default=false"`
	Host        string `env:"HOST,default=0.0.0.0"`
	Port        int    `env:"PORT,default=8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	PublicDir   string `env:"PUBLIC_DIR"`
	Email       EmailConfig
}

func (c Config) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func LoadEnvConfig() *Config {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatalf("failed to process env config: %v", err)
	}
	return &c
}
