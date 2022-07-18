package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Environment string `mapstructure:"environment"`

	API api `mapstructure:"api"`
	DB  DB  `mapstructure:"db"`
}

type api struct {
	ServeSwagger bool          `mapstructure:"serve_swagger"`
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	PersonLimit  int           `mapstructure:"person_limit"`
}

type DB struct {
	URL          string `mapstructure:"url"`
	SchemaName   string `mapstructure:"schema_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

var defaults = map[string]interface{}{
	"environment":      "development",
	"shutdown_timeout": time.Second * 5,

	"db.url":            "postgres://root@localhost:5432/root?sslmode=disable",
	"db.schema_name":    "public",
	"db.max_open_conns": 2,
	"db.max_idle_conns": 2,

	"api.serve_swagger": true,
	"api.address":       ":3000",
	"api.read_timeout":  time.Second * 5,
	"api.write_timeout": time.Second * 5,
	"api.person_limit":  3,
}

func New() (*Config, error) {
	viper.AddConfigPath(".")

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	var c Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read config, using defaults: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
