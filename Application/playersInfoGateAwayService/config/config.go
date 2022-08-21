package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Network                  string `mapstructure:"NETWORK_TYPE"`
	HostGrpcPort             string `mapstructure:"HOST_GRPC_PORT"`
	HostRestPort             string `mapstructure:"HOST_REST_PORT"`
	PlayerInfoServiceAddress string `mapstructure:"PLAYER_INFO_SERVICE_ADDRESS"`
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
	if err != nil {
		_ = fmt.Errorf("do not parse config file:%v", err)
	}

	return
}
