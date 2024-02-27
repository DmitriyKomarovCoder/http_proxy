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

	MongoDB struct {
		Url         string `mapstructure:"url"`
		DbName      string `mapstructure:"dbname"`
		ColRequest  string `mapstructure:"colrequest"`
		ColResponse string `mapstructure:"colresponse"`
	} `mapstructure:"mongodb"`

	Certificate struct {
		Key     string `mapstructure:"key"`
		Cert    string `mapstructure:"cert"`
		Subject string `mapstructure:"subject"`
	} `mapstructure:"certificateinfo"`

	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

func InitialConfig(nameConfig, pathConfig string) (*Config, error) {
	var config Config
	viper.SetConfigName(nameConfig)
	viper.AddConfigPath(pathConfig)
	err := viper.ReadInConfig()
	if err != nil {
		return &config, fmt.Errorf("Fatal error config file: %s \n", err)

	}

	err = viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return &config, fmt.Errorf("unable to decode into struct: %v", err)
	}

	return &config, nil
}
