package config

import (
	"fmt"

	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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
	Schema   string `mapstructure:"schema"`
	SSLMode  string `mapstructure:"ssl_mode"`
	LogMode  string `mapstructure:"log_mode"`
}

type LoggerConfig struct {
	Level        string `mapstructure:"level"`
	FilePath     string `mapstructure:"file_path"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups"`
	MaxAge       int    `mapstructure:"max_age"`
	Format       string `mapstructure:"format"`
	ReportCaller bool   `mapstructure:"report_caller"`
	Service      string `mapstructure:"service"`
}

type MinioConfig struct {
	Addr       string `mapstructure:"addr"`
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
	SSL        bool   `mapstructure:"ssl"`
	BucketName string `mapstructure:"bucket_name"`
}

type APMConfig struct {
	ServerURL   string `mapstructure:"server_url"`
	SecretToken string `mapstructure:"secret_token"`
	ServiceName string `mapstructure:"service_name"`
	Environment string `mapstructure:"environment"`
}

type TranslationConfig struct {
	Folder string `mapstructure:"folder"`
}

type Config struct {
	ServerConfig      ServerConfig                `mapstructure:"server"`
	RedisConfig       RedisConfig                 `mapstructure:"redis"`
	DatabasesConfig   map[string][]DatabaseConfig `mapstructure:"databases"`
	LoggerConfig      LoggerConfig                `mapstructure:"logger"`
	MinioConfig       MinioConfig                 `mapstructure:"minio"`
	APMConfig         APMConfig                   `mapstructure:"apm"`
	TranslationConfig TranslationConfig           `mapstructure:"translation"`
}

var (
	configSingleton *singleton.Singleton[Config]
)

func InitConfig(path string) error {
	fmt.Println("Init config with file path:", path)
	configSingleton = singleton.NewSingleton(func() (config Config) {
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
