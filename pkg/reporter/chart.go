package reporter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ajbosco/statboard/pkg/metric"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// RenderChart formats the Statboard metrics and renders a go-chart time series plot
func RenderChart(name string, chartColor string, filePath string, metrics []metric.Metric) error {
	if len(metrics) <= 1 {
		return fmt.Errorf("time series charts require more than 1 value, only found %d. collect more data", len(metrics))
	}
	var series []chart.Series
	var x []time.Time
	var y []float64

	for _, met := range metrics {
		x = append(x, met.Date)
		y = append(y, met.Value)
	}

	ts := chart.TimeSeries{
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex(chartColor),
			FillColor:   drawing.ColorFromHex(chartColor),
		},
		XValues: x,
		YValues: y,
	}
	if err := ts.Validate(); err != nil {
		return errors.Wrap(err, "invalid time series chart")
	}

	series = append(series, ts)

	fileName := fmt.Sprintf("%v-chart.png", name)

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
			ValueFormatter: func(v interface{}) string {
				if typed, isTyped := v.(float64); isTyped {
					return time.Unix(0, int64(typed)).UTC().Format("2006-01-02")
				}
				return fmt.Sprint(v)
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Series: series,
	}

	return writeChart(graph, filePath, fileName)
}

func writeChart(c chart.Chart, filePath string, fileName string) error {
	if err := os.MkdirAll(filePath, 0777); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create chart directory:%q", filePath))
	}

	buf := bytes.NewBuffer([]byte{})
	fullPath := fmt.Sprintf("%v/%v", filePath, fileName)

	if err := c.Render(chart.PNG, buf); err != nil {
		return errors.Wrap(err, "failed to render chart")
	}

	if err := ioutil.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write chart to filePath:%q", fullPath))
	}
	logrus.Info(fmt.Sprintf("chart written to %q", fullPath))

	return nil
}
