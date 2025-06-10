package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App        string `mapstructure:"APP"`
	Env        string `mapstructure:"ENV"`
	SvcVersion string `mapstructure:"SVC_VERSION"`
	Port       string `mapstructure:"PORT"`

	Database struct {
		ServiceName string `mapstructure:"DATABASE_SERVICE_NAME"`
		DSN         string `mapstructure:"DATABASE_DSN"`
		MaxOpenConn int    `mapstructure:"DATABASE_MAX_OPEN_CONN"`
		MaxIdleConn int    `mapstructure:"DATABASE_MAX_IDLE_CONN"`
	}


	HTTPClient struct {
		TimeoutMS           int  `mapstructure:"HTTP_CLIENT_TIMEOUT_MS"`
		DisableKeepAlives   bool `mapstructure:"HTTP_CLIENT_DISABLE_KEEP_ALIVE"`
		MaxIdleConns        int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS"`
		MaxConnsPerHost     int  `mapstructure:"HTTP_CLIENT_MAX_CONNS_PER_HOST"`
		MaxIdleConnsPerHost int  `mapstructure:"HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST"`
		IdleConnTimeout     int  `mapstructure:"HTTP_CLIENT_IDLE_CONN_TIMEOUT"`
	}

	JWT struct {
		SecretKey string `mapstructure:"JWT_SECRET_KEY"`
	}

}

func NewConfig() (*Config, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config.Database)
	if err != nil {
		return nil, err
	}


	err = viper.Unmarshal(&config.HTTPClient)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config.JWT)
	if err != nil {
		return nil, err
	}

	return &config, nil
}