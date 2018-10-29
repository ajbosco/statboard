package reporter

import (
	"fmt"
	"html/template"
	"strings"

	chartjs "github.com/ajbosco/goChartjs"
	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/pkg/errors"
	colors "gopkg.in/go-playground/colors.v1"
)

// Chart contains information for creating a new metric chart
type Chart struct {
	metricName string
	ChartName  string
	color      string
	metrics    []statboard.Metric
	ChartJS    template.HTML
}

// NewChart returns a chart object
func NewChart(metricName string, chartName string, color string, metrics []statboard.Metric) (Chart, error) {
	validName := strings.Replace(metricName, ".", "_", -1)
	hex, err := colors.ParseHEX(color)
	if err != nil {
		return Chart{}, errors.Wrap(err, fmt.Sprintf("failed to parse color %q", color))
	}

	return Chart{metricName: validName, ChartName: chartName, color: hex.ToRGB().String(), metrics: metrics}, nil
}

// RenderChart formats the Statboard metrics and returns chart.js string
func (c *Chart) RenderChart() (string, error) {
	chartData := metricsToPoints(c.metrics)

	chart := getChart(chartData, c.metricName, c.color)

	s, err := chart.Render()
	if err != nil {
		return "", errors.Wrap(err, "rendering chart failed")
	}

	return s, nil
}

// metricsToPoints converts statboard Metrics to chart.js points
func metricsToPoints(metrics []statboard.Metric) []chartjs.Point {
	var data []chartjs.Point

	for _, met := range metrics {
		point := chartjs.Point{X: met.Date.Format("02-Jan-2006"), Y: met.Value}
		data = append(data, point)
	}

	return data
}

// getChart formats and returns Chartjs object
func getChart(chartData []chartjs.Point, metricName string, chartColor string) chartjs.Chart {
	lineTension := 0

	dataset := []chartjs.Dataset{{
		Label:           metricName,
		LineTension:     &lineTension,
		Data:            chartData,
		BackgroundColor: chartColor,
	},
	}

	c := chartjs.Chart{
		Name:      metricName,
		ChartType: "line",
		Options: &chartjs.Options{

			Scales: chartjs.Scales{
				XAxes: []chartjs.Axes{
					{
						Time: &chartjs.Time{
							Unit: "day",
						},
						Type: "time",
					},
				},
			},
			MaintainAspectRatio: chartjs.False,
			Responsive:          chartjs.True,
		},
	}

	c.Data.Datasets = dataset

	return c
}
