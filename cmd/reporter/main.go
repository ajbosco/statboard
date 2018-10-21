package main

import (
	"fmt"
	"time"

	"github.com/ajbosco/statboard/pkg/config"
	"github.com/ajbosco/statboard/pkg/metric"
	"github.com/ajbosco/statboard/pkg/reporter"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// EnvConfig contains environment variables for the metric reporter
type EnvConfig struct {
	ConfigFilePath string `required:"true"`
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

	var cfg config.Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	//Create metric store
	s, err := metric.NewStormStore(cfg.Db.FilePath)
	if err != nil {
		logrus.Fatal(err)
	}

	// Render charts for all metrics
	for metType, metCfgs := range cfg.Metrics {
		for metName, metCfg := range metCfgs {
			metricName := fmt.Sprintf("%s.%s", metType, metName)
			met, err := s.GetMetric(metricName, time.Now().AddDate(0, 0, -metCfg.DaysBack))
			if err != nil {
				logrus.Fatal(err)
			}

			if len(met) == 0 {
				logrus.Info(fmt.Sprintf("no metrics returned for %q", metricName))
				return
			}

			if err := reporter.RenderChart(metricName, metCfg.Color, cfg.Charts.DirPath, met); err != nil {
				logrus.Fatal(err)
			}
		}
	}
}
