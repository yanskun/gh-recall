package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Days     int    `mapstructure:"days"`
	Locale   string `mapstructure:"locale"`
	Model    string `mapstructure:"model"`
	Port     int    `mapstructure:"port"`
	Sections int    `mapstructure:"sections"`
}

func LoadConfig() Config {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.config/gh-recall")

	setDefaultValues()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(); err != nil {
				log.Fatalf("Error creating config file: %v", err)
			}
		} else {
			log.Fatalf("Error reading config")
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %s\n", err)
	}

	return config
}

func createDefaultConfig() error {
	configPath := fmt.Sprintf("%s/.config/gh-recall", os.Getenv("HOME"))
	configFile := fmt.Sprintf("%s/config.toml", configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return err
		}
	}

	setDefaultValues()

	if err := viper.WriteConfigAs(configFile); err != nil {
		return err
	}

	return nil
}

func setDefaultValues() {
	viper.SetDefault("days", 7)
	viper.SetDefault("locale", "en")
	viper.SetDefault("model", "phi4")
	viper.SetDefault("port", 11434)
	viper.SetDefault("sections", 3)
}
