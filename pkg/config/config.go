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
	DSN       string `mapstructure:"dsn"`
	IsCluster bool   `mapstructure:"is_cluster"`
}

type DatabaseConfig struct {
	Engine   string `mapstructure:"engine"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"user"`
	Password string `mapstructure:"pass"`
	Database string `mapstructure:"database"`
	Instance string `mapstructure:"instance"`
}

type LoggerConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
	Format   string `mapstructure:"format"`
}

type Config struct {
	ServerConfig    ServerConfig                `mapstructure:"server"`
	RedisConfig     RedisConfig                 `mapstructure:"redis"`
	DatabasesConfig map[string][]DatabaseConfig `mapstructure:"databases"`
	LoggerConfig    LoggerConfig                `mapstructure:"logger"`
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
	}, true)
	return nil
}

func GetConfig() Config {
	return configSingleton.Get()
}
