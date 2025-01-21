package config

import (
	"errors"
	"fmt"
	"io/fs"
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
	PublicKey  []byte
	PrivateKey []byte
}

var config *Config

func Init() (*Config, error) {

	c := Config{}

	dir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	keysDirectory := os.DirFS(fmt.Sprintf("%s/keys", dir))

	keyFiles, err := fs.Glob(keysDirectory, "*.pem")

	if err != nil {
		return nil, err
	}

	var publicKey string
	var privateKey string

	for _, v := range keyFiles {
		if v == "public.pem" {
			publicKey = v
		}
		if v == "private.pem" {
			privateKey = v
		}
	}

	if publicKey == "" || privateKey == "" {
		return nil, errors.New("public.pem and private.pem files are required")
	}

	publicKeyContent, err := os.ReadFile(fmt.Sprintf("%s/keys/%s", dir, publicKey))

	if err != nil {
		return nil, err
	}

	c.PublicKey = publicKeyContent

	privateKeyContent, err := os.ReadFile(fmt.Sprintf("%s/keys/%s", dir, privateKey))

	if err != nil {
		return nil, err
	}

	c.PrivateKey = privateKeyContent

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

	config = &c

	return &c, nil
}

func Get() *Config {

	if config == nil {
		panic("config hasn't initialized yet")
	}

	return config
}
