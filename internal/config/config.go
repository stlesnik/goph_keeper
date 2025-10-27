package config

import (
	"errors"
	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	ServerAddress string `mapstructure:"server_address"`
	PostgresDSN   string `mapstructure:"postgres_dsn"`
	Environment   string `mapstructure:"environment"`
	JWTSecret     string `mapstructure:"jwt_secret"`
	TLSCertFile   string `mapstructure:"tls_cert_file"`
	TLSKeyFile    string `mapstructure:"tls_key_file"`
	EnableHTTPS   bool   `mapstructure:"enable_https"`
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
	err = config.Validate()
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) Validate() error {
	if c.PostgresDSN == "" {
		return errors.New("postgres DSN is required")
	}
	if c.JWTSecret == "" {
		return errors.New("JWT secret is required")
	}
	if len(c.JWTSecret) < 32 {
		return errors.New("JWT secret must be at least 32 characters")
	}
	return nil
}
