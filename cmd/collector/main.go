package main

import (
	"fmt"

	"github.com/ajbosco/statboard/pkg/collector"
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

	// Create collectors
	for k, v := range metrics {
		c, err := createCollector(k, cfg.ConfigFilePath)
		if err != nil {
			logrus.Fatal(err)
		}

		// Collect and write metrics
		for _, m := range v {
			metrics, err := c.Collect(m, 10)
			if err != nil {
				logrus.Fatal(errors.Wrap(err, fmt.Sprintf("failed to collect metric:%q", m)))
			}

			for _, met := range metrics {
				err = s.WriteMetric(met)
				if err != nil {
					logrus.Fatal(errors.Wrap(err, fmt.Sprintf("failed to write metric:%q", m)))
				}
			}

		}
	}
}

func createCollector(collectorType string, cfgFilePath string) (collector.Collector, error) {
	var c collector.Collector
	var err error

	switch collectorType {
	case "fitbit":
		c, err = collector.NewFitbitCollector(cfgFilePath)
	default:
		err = fmt.Errorf("Unsupported collector type:%q", collectorType)
	}
	return c, err
}
