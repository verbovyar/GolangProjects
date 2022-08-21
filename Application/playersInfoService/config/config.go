package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	ApiKey string `mapstructure:"API_KEY"`

	ConnectionString string `mapstructure:"CONNECTION_STRING"`

	Network string `mapstructure:"NETWORK_TYPE"`
	Port    string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf("do not parse config file:%v", err)
	}

	err = viper.Unmarshal(&config)

	return
}
