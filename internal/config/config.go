package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	OpenAIAPIKey string
}

func NewConfig() *Config {
	return &Config{
		OpenAIAPIKey: GetAPIKey(),
	}
}

func Init() {
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("OPENAI_API_KEY", "demo")

	// Read from .env file
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func GetAPIKey() string {
	return viper.GetString("API_KEY")
}
