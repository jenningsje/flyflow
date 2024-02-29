package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	OpenAIAPIKey string
}

func NewConfig() *Config {
	Init()
	return &Config{
		OpenAIAPIKey: GetAPIKey(),
	}
}

func Init() {
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("OPENAI_API_KEY", "demo")
}

func GetAPIKey() string {
	return viper.GetString("OPENAI_API_KEY")
}
