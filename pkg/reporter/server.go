package reporter

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ajbosco/statboard/pkg/config"
	"github.com/ajbosco/statboard/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg    config.Config
	addr   string
	store  storage.Store
	router *mux.Router
}

func NewServer(cfg config.Config, addr string, store storage.Store) *Server {
	return &Server{cfg: cfg, addr: addr, store: store, router: mux.NewRouter()}
}

func (s *Server) ListenAndServe() {
	s.routes()
	http.ListenAndServe(s.addr, s.router)
}

func (s *Server) getChartJs() ([]chart, error) {
	var charts []chart

	// Render charts for all metrics
	for metType, metCfgs := range s.cfg.Metrics {
		for metName, metCfg := range metCfgs {

			// Fetch metric values from database
			metricName := fmt.Sprintf("%s.%s", metType, metName)
			sinceDt := time.Now().AddDate(0, -metCfg.ChartMonthsBack, 0)
			firstOfMonth := time.Date(sinceDt.Year(), sinceDt.Month(), 1, 0, 0, 0, 0, time.UTC).Truncate(24 * time.Hour)
			met, err := s.store.GetMetric(metricName, firstOfMonth)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get metric")
			}
			if len(met) == 0 {
				logrus.Info(fmt.Sprintf("no metrics returned for %q", metricName))
				continue
			}

			// Render chart for  metric values
			chart, err := newChart(metricName, metCfg.ChartName, metCfg.ChartColor, met)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create new chart")
			}
			chartString, err := chart.renderChart()
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
