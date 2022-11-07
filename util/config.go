package util

import "github.com/spf13/viper"

type Config struct {
	RabbitMQ RabbitMQConfig `mapstructure:",squash"`
	Postgres PostgresConfig `mapstructure:",squash"`
}

type RabbitMQConfig struct {
	Username string `mapstructure:"RABBITMQ_USER"`
	Password string `mapstructure:"RABBITMQ_PASS"`
	Host     string `mapstructure:"RABBITMQ_HOST"`
	Port     uint16 `mapstructure:"RABBITMQ_PORT"`
	Vhost    string `mapstructure:"RABBITMQ_VHOST"`
}

type PostgresConfig struct {
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     uint16 `mapstructure:"POSTGRES_PORT"`
	Database string `mapstructure:"POSTGRES_DB"`
}

func LoadConfig() (*Config, error) {
	config := Config{}
	viper.SetConfigFile(".env")

	setDefaults()
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return &config, err
	}

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("RABBITMQ_HOST", "localhost")
	viper.SetDefault("RABBITMQ_PORT", 5672)

	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
}
