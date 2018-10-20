package main

import (
	"fmt"
	"log"

	"github.com/ajbosco/statboard/pkg/collector"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// Config contains environment variables for the metric collector
type Config struct {
	ConfigFilePath string `required:"true"`
	DbFilePath     string `required:"true"`
	CollectorType  string `required:"true"`
	MetricName     string `required:"true"`
}

func main() {
	var cfg Config
	err := envconfig.Process("statboard", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	c, err := createCollector(cfg.CollectorType, cfg.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	s, err := collector.NewBoltStore(cfg.DbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	metricName := cfg.MetricName
	metrics, err := c.Collect(metricName, 10)
	if err != nil {
		log.Fatal(errors.Wrap(err, fmt.Sprintf("failed to collect metric:%q", metricName)))
	}

	for _, met := range metrics {
		err = s.WriteMetric(met)
		if err != nil {
			log.Fatal(errors.Wrap(err, fmt.Sprintf("failed to write metric:%q", metricName)))
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
