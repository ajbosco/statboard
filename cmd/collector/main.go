package main

import (
	"fmt"
	"sync"

	"github.com/ajbosco/statboard/pkg/collector"
	"github.com/ajbosco/statboard/pkg/config"
	"github.com/ajbosco/statboard/pkg/storage"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// EnvConfig contains environment variables for the metric collector
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

	var collectWg sync.WaitGroup
	collectWg.Add(len(cfg.Metrics))

	// Create collectors
	for metType, metCfgs := range cfg.Metrics {
		go func(metType string, metCfgs map[string]config.MetricConfig) {
			defer collectWg.Done()
			var metWg sync.WaitGroup
			metWg.Add(len(metCfgs))

			c, err := createCollector(metType, cfg)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Info(fmt.Sprintf("collecting metrics for %q", metType))

			// Collect and write metrics
			for metName, metCfg := range metCfgs {
				go func(metName string, metCfg config.MetricConfig) {
					defer metWg.Done()
					metricName := fmt.Sprintf("%s.%s", metType, metName)
					logrus.Info(fmt.Sprintf("collecting %q", metricName))
					metrics, err := c.Collect(metName, metCfg.CollectMonthsBack)
					if err != nil {
						logrus.Fatal(errors.Wrap(err, fmt.Sprintf("failed to collect metric:%q", metricName)))
					}
					logrus.Info(fmt.Sprintf("collected %d %q records", len(metrics), metricName))

					for _, met := range metrics {
						err = s.WriteMetric(met)
						if err != nil {
							logrus.Fatal(errors.Wrap(err, fmt.Sprintf("failed to write metric:%q", metricName)))
						}
					}
					logrus.Info(fmt.Sprintf("wrote %d %q records to database", len(metrics), metricName))
				}(metName, metCfg)
			}
			metWg.Wait()
			logrus.Info(fmt.Sprintf("completed collecting metrics for %q", metType))
		}(metType, metCfgs)
	}
	collectWg.Wait()
}

func createCollector(collectorType string, cfg config.Config) (collector.Collector, error) {
	var c collector.Collector
	var err error

	switch collectorType {
	case "fitbit":
		c, err = collector.NewFitbitCollector(cfg)
	case "github":
		c, err = collector.NewGithubCollector(cfg)
	case "goodreads":
		c, err = collector.NewGoodreadsCollector(cfg)
	default:
		err = fmt.Errorf("Unsupported collector type:%q", collectorType)
	}
	return c, err
}
