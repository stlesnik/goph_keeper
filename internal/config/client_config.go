package config

import (
	"errors"
	"github.com/spf13/viper"
)

// ClientConfig holds client configuration
type ClientConfig struct {
	ServerURL string `mapstructure:"server_url"`
	TokenFile string `mapstructure:"token_file"`
	CertFile  string `mapstructure:"cert_file"`
}

// LoadClientConfig creates a new ClientConfig by parsing configuration file, environment variables, and flags.
func LoadClientConfig() (ClientConfig, error) {
	viper.SetConfigName("client_config")
	viper.AddConfigPath("configs")

	viper.SetDefault("server_url", ":8080")

	err := viper.ReadInConfig()
	if err != nil {
		return ClientConfig{}, err
	}
	var config ClientConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		return ClientConfig{}, err
	}
	err = config.ValidateClientConfig()
	if err != nil {
		return ClientConfig{}, err
	}
	return config, nil
}

// ValidateClientConfig checks if all required parameters are set
func (c *ClientConfig) ValidateClientConfig() error {
	if c.CertFile == "" {
		return errors.New("JWT secret is required")
	}
	return nil
}
