package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	AppName    string `mapstructure:"APP_NAME"`
	AppPort    string `mapstructure:"APP_PORT"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbName     string `mapstructure:"DB_NAME"`
	DbUsername string `mapstructure:"DB_USERNAME"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
}

func Init() (*Config, error) {

	c := Config{}

	dir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	viper.AddConfigPath(dir)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&c)

	if err != nil {
		return nil, err
	}

	fmt.Println(c)

	return &c, nil
}
