package reporter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/pkg/errors"
	chart "github.com/wcharczuk/go-chart"
)

// RenderChart formats the Statboard metrics and renders a go-chart time series plot
func RenderChart(name string, filePath string, metrics []statboard.Metric) error {
	var series []chart.Series
	var x []time.Time
	var y []float64

	for _, met := range metrics {
		x = append(x, met.Date)
		y = append(y, met.Value)
	}

	series = append(series, chart.TimeSeries{
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: x,
		YValues: y,
	})

	fileName := fmt.Sprintf("%v-chart.png", name)

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
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

	return nil
}
