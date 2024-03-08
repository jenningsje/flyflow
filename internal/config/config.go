package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	OpenAIAPIKey string
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPass       string
	DBName       string
	Env          string
}

func NewConfig() (*Config, error) {
	// Look for a .env file
	viper.SetConfigFile(".env")

	// Read the configuration from the .env file
	err := viper.ReadInConfig()

	// Use AutomaticEnv to override configuration with environment variables
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("OPENAI_API_KEY", "demo")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASS", "password")
	viper.SetDefault("DB_NAME", "flyflow")
	viper.SetDefault("ENV", "local")

	// Return the config
	return &Config{
		OpenAIAPIKey: viper.GetString("OPENAI_API_KEY"),
		Port:         viper.GetString("PORT"),
		DBHost:       viper.GetString("DB_HOST"),
		DBPort:       viper.GetString("DB_PORT"),
		DBUser:       viper.GetString("DB_USER"),
		DBPass:       viper.GetString("DB_PASS"),
		DBName:       viper.GetString("DB_NAME"),
		Env:          viper.GetString("ENV"),
	}, err
}
