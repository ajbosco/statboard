package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ajbosco/statboard/pkg/config"
	"github.com/ajbosco/statboard/pkg/reporter"
	"github.com/ajbosco/statboard/pkg/storage"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// EnvConfig contains environment variables for the metric reporter
type EnvConfig struct {
	ConfigFilePath string `required:"true"`
	DbFilePath     string `required:"true"`
}

type dashboard struct {
	charts []reporter.Chart
	cfg    config.Config
	store  storage.Store
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

	dashboard := &dashboard{cfg: cfg, store: s}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/static/", fs)
	http.HandleFunc("/", dashboard.handler)

	// start web server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (d *dashboard) handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	charts, err := d.getChartJs()
	if err != nil {
		logrus.Fatal(err)
	}

	err = tmpl.Execute(w, charts)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (d *dashboard) getChartJs() ([]reporter.Chart, error) {
	var charts []reporter.Chart

	// Render charts for all metrics
	for metType, metCfgs := range d.cfg.Metrics {
		for metName, metCfg := range metCfgs {

			// Fetch metric values from database
			metricName := fmt.Sprintf("%s.%s", metType, metName)
			met, err := d.store.GetMetric(metricName, time.Now().AddDate(0, -metCfg.ChartMonthsBack, 0))
			if err != nil {
				return nil, errors.Wrap(err, "failed to get metric")
			}
			if len(met) == 0 {
				logrus.Info(fmt.Sprintf("no metrics returned for %q", metricName))
				continue
			}

			// Render chart for  metric values
			chart, err := reporter.NewChart(metricName, metCfg.ChartName, metCfg.ChartColor, met)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create new chart")
			}
			chartString, err := chart.RenderChart()
			if err != nil {
				return nil, errors.Wrap(err, "failed to render chart")
			}
			chart.ChartJS = template.HTML(chartString)

			// append charts to dashboard
			charts = append(charts, chart)
		}
	}
	return charts, nil
}
