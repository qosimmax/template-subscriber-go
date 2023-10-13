// Package config handles environment variables.
package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

// Config contains environment variables.
type Config struct {
	Port                       string  `envconfig:"PORT" default:"8000"`
	ServiceName                string  `envconfig:"SERVICE_NAME" required:"true"`
	Environment                string  `envconfig:"ENVIRONMENT" required:"true"`
	JaegerAgentHost            string  `envconfig:"JAEGER_AGENT_HOST" default:"localhost"`
	JaegerAgentPort            string  `envconfig:"JAEGER_AGENT_PORT" default:"6831"`
	JaegerSamplerType          string  `envconfig:"JAEGER_SAMPLER_TYPE" default:"const"`
	JaegerSamplerParam         float64 `envconfig:"JAEGER_SAMPLER_PARAM" default:"1"`
	DatabasePassword           string  `envconfig:"DATABASE_PASSWORD" required:"true"`
	DatabaseUser               string  `envconfig:"DATABASE_USER" required:"true"`
	DatabaseURL                string  `envconfig:"DATABASE_URL" default:"127.0.0.1"`
	DatabasePort               string  `envconfig:"DATABASE_PORT" default:"5432"`
	DatabaseDB                 string  `envconfig:"DATABASE_DB" default:"postgres"`
	DatabaseOptions            string  `envconfig:"DATABASE_OPTIONS" default:"?sslmode=disable"`
	DatabaseMaxConnections     int     `envconfig:"DATABASE_MAX_CONNECTIONS" default:"12"`
	DatabaseMaxIdleConnections int     `envconfig:"DATABASE_MAX_IDLE_CONNECTIONS" default:"3"`
	NatsURL                    string  `envconfig:"NATS_URL" required:"true"`
}

// LoadConfig reads environment variables and populates Config.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("~/MyTaxi/mytaxi/billing/billing_processor/.env"); err != nil {
			log.Info("No .env file found...")
		}
		log.Info("No .env file found")
	}

	var c Config

	err := envconfig.Process("", &c)

	return &c, err
}
