package util

import "github.com/spf13/viper"

type Config struct {
	RabbitMQ rabbitMQConfig `mapstructure:",squash"`
	Postgres postgresConfig `mapstructure:",squash"`
}

type rabbitMQConfig struct {
	Username string `mapstructure:"RABBITMQ_USER"`
	Password string `mapstructure:"RABBITMQ_PASS"`
	Host     string `mapstructure:"RABBITMQ_HOST"`
	Port     string `mapstructure:"RABBITMQ_PORT"`
	Vhost    string `mapstructure:"RABBITMQ_VHOST"`
}

type postgresConfig struct {
	Username string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	Database string `mapstructure:"POSTGRES_DB"`
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
