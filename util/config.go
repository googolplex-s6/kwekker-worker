package util

import "github.com/spf13/viper"

type Config struct {
	RabbitMQURI string `mapstructure:"RABBITMQ_URI"`
	DatabaseURI string `mapstructure:"DATABASE_URI"`
}

func LoadConfig() (Config, error) {
	config := Config{}
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)

	return config, nil
}
