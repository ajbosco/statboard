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
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// EnvConfig contains environment variables for the metric reporter
type EnvConfig struct {
	ConfigFilePath string `required:"true"`
}

type dashboard struct {
	charts []reporter.Chart
}

func (d *dashboard) handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	err := tmpl.Execute(w, d.charts)
	if err != nil {
		logrus.Fatal(err)
	}
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

	// Create metric store
	s, err := storage.NewStormStore(cfg.Db.FilePath)
	if err != nil {
		logrus.Fatal(err)
	}

	dashboard := &dashboard{}

	// Render charts for all metrics
	for metType, metCfgs := range cfg.Metrics {
		for metName, metCfg := range metCfgs {

			// Fetch metric values from database
			metricName := fmt.Sprintf("%s.%s", metType, metName)
			met, err := s.GetMetric(metricName, time.Now().AddDate(0, 0, -metCfg.DaysBack))
			if err != nil {
				logrus.Fatal(err)
			}
			if len(met) == 0 {
				logrus.Info(fmt.Sprintf("no metrics returned for %q", metricName))
				return
			}

			// Render chart for  metric values
			chart, err := reporter.NewChart(metricName, metCfg.ChartName, metCfg.Color, met)
			if err != nil {
				logrus.Fatal(err)
			}
			chartString, err := chart.RenderChart()
			if err != nil {
				logrus.Fatal(err)
			}
			chart.ChartJS = template.HTML(chartString)

			// append charts to dashboard
			dashboard.charts = append(dashboard.charts, chart)
		}
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/static/", fs)
	http.HandleFunc("/", dashboard.handler)

	// start web server
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Fatal(err)
	}
}
