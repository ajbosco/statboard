package main

import (
	"time"

	"github.com/ajbosco/statboard/pkg/reporter"
	"github.com/ajbosco/statboard/pkg/store"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config contains environment variables for the metric collector
type Config struct {
	ConfigFilePath string `required:"true"`
}

func main() {
	var cfg Config
	err := envconfig.Process("statboard", &cfg)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	viper.SetConfigFile(cfg.ConfigFilePath)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	// Check for required config parameters
	filePath := viper.GetString("charts.file_path")
	if filePath == "" {
		logrus.Fatal(errors.New("'charts.file_path' must be present in config"))
	}
	dbFilePath := viper.GetString("db.file_path")
	if dbFilePath == "" {
		logrus.Fatal(errors.New("'db.file_path' must be present in config"))
	}
	metrics := viper.GetStringMapStringSlice("metrics")
	if metrics == nil {
		logrus.Fatal(errors.New("'metrics' must be present in config"))
	}

	// Create statboard store
	s, err := store.NewBoltStore(dbFilePath)
	if err != nil {
		logrus.Fatal(err)
	}

	// Render charts for all metrics
	for _, v := range metrics {
		for _, m := range v {
			met, err := s.GetMetric(m, time.Now().AddDate(0, 0, -10))
			if err != nil {
				logrus.Fatal(err)
			}

			if err := reporter.RenderChart(m, filePath, met); err != nil {
				logrus.Fatal(err)
			}
		}
	}
}
