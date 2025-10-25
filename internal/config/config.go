package config

import (
	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	ServerAddress string `mapstructure:"server_address"`
	PostgresDSN   string `mapstructure:"postgres_dsn"`
	Environment   string `mapstructure:"environment"`
	JWTSecret     string `mapstructure:"jwt_secret"`
}

// Load creates a new Config by parsing configuration file, environment variables, and flags.
func Load() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("configs")

	viper.SetDefault("server_address", ":8080")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
