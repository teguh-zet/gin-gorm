package helper

import (
	"github.com/spf13/viper"
)

type Config struct {
	// Server Configuration
	PORT     string `mapstructure:"PORT"` // Changed from AppPort to PORT to match user example
	GIN_MODE string `mapstructure:"GIN_MODE"`
	SCHEMA   string `mapstructure:"SCHEMA"`

	// Database
	DB string `mapstructure:"DB"`

	// CORS & LOG
	ALLOW_ORIGIN string `mapstructure:"ALLOW_ORIGIN"`
	LOG_FILE     string `mapstructure:"LOG_FILE"`
	AUTO_MIGRATE string `mapstructure:"AUTO_MIGRATE"`

	// NATS
	NatsServers string `mapstructure:"NATS_SERVERS"` // Matches user example
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
