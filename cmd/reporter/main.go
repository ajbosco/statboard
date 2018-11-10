package main

import (
	"github.com/ajbosco/statboard/pkg/config"
	"github.com/ajbosco/statboard/pkg/reporter"
	"github.com/ajbosco/statboard/pkg/storage"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// EnvConfig contains environment variables for the metric reporter
type EnvConfig struct {
	ConfigFilePath string `required:"true"`
	DbFilePath     string `required:"true"`
}

func main() {
	var envCfg EnvConfig
	err := envconfig.Process("statboard", &envCfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	viper.SetConfigFile(envCfg.ConfigFilePath)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	// Create metric store
	s, err := storage.NewStormStore(envCfg.DbFilePath)
	if err != nil {
		logrus.Fatal(err)
	}

	var cfg config.Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := reporter.NewServer(cfg, ":8080", s)
	srv.ListenAndServe()
}
