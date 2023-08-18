package config

import (
	"fmt"

	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type RedisConfig struct {
	Addrs    []string `mapstructure:"addrs"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	DB       int      `mapstructure:"db"`
}
type TelegramConfig struct {
	APIKey string `mapstructure:"api_key"`
}

type Config struct {
	ServerConfig   ServerConfig   `mapstructure:"server"`
	TelegramConfig TelegramConfig `mapstructure:"telegram"`
	RedisConfig    RedisConfig    `mapstructure:"redis"`
}

var (
	configSingleton *singleton.Singleton[Config]
)

func InitConfig(path string) error {
	fmt.Println("Init config with file path:", path)
	configSingleton = singleton.NewSingleton[Config](func() (config Config) {
		viper := viper.New()
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
		return config
	})
	return nil
}

func GetConfig() Config {
	return configSingleton.Get()
}
