package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ApiServer struct {
		Host         string        `mapstructure:"host"`
		Port         string        `mapstructure:"port"`
		ReadTimeout  time.Duration `mapstructure:"readTimeout"`
		WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	} `mapstructure:"apiServer"`
	ProxyServer struct {
		Host         string        `mapstructure:"host"`
		Port         string        `mapstructure:"port"`
		ReadTimeout  time.Duration `mapstructure:"readTimeout"`
		WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	} `mapstructure:"proxyServer"`
	Logfile struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"logfile"`
	Postgres struct {
		Name     string `mapstructure:"DB_NAME"`
		User     string `mapstructure:"DB_USER"`
		Password string `mapstructure:"DB_PASSWORD"`
		Host     string `mapstructure:"DB_HOST"`
		Port     int    `mapstructure:"DB_PORT"`
	}
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

func InitialConfig(nameConfig, pathConfig string) (*Config, error) {
	var config Config
	viper.SetConfigName(nameConfig)
	viper.AddConfigPath(pathConfig)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return &config, fmt.Errorf("Fatal error config file: %s \n", err)

	}

	viper.SetConfigFile(".env")

	err = viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return &config, fmt.Errorf("unable to decode into struct: %v", err)
	}
	config.Postgres.Name = viper.GetString("DB_NAME")
	config.Postgres.Host = viper.GetString("DB_HOST")
	config.Postgres.User = viper.GetString("DB_USER")
	config.Postgres.Port = viper.GetInt("DB_PORT")
	config.Postgres.Password = viper.GetString("DB_PASSWORD")
	return &config, nil
}
