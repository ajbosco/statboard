package reporter

import (
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) handleDashboard() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		charts, err := s.getChartJs()
		if err != nil {
			logrus.Fatal(err)
		}

		err = tmpl.Execute(w, charts)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func (s *Server) handleFavicon() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/static/favicon.ico")
	}
}
