package infrastructure

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	Production  string = "Production"
	Development        = "Development"
)

type Config struct {
	Env              string `mapstructure:"ENV"`
	ServerPort       int    `mapstructure:"SERVER_PORT"`
	DatabaseHost     string `mapstructure:"DATABASE_HOST"`
	DatabasePort     int    `mapstructure:"DATABASE_PORT"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName     string `mapstructure:"DATABASE_NAME"`
	JwtSecret        string `mapstructure:"JWT_SECRET"`
	ElasticSearch    string `mapstructure:"ELASTIC_SEARCH_HOST"`
}

func ParseConfig() Config {
	var config Config

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("got error while reading config file: %w", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("got error while parsing config: %w", err))
	}

	return config
}
