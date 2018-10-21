package metric

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriteMetric(t *testing.T) {
	b, err := NewStormStore("test.db")
	assert.NoError(t, err)

	defer os.Remove("test.db")
	defer b.Close()
	metric := Metric{
		Name:  "testMetric",
		Date:  time.Now(),
		Value: 3.0,
	}

	err = b.WriteMetric(metric)
	assert.NoError(t, err)
}

func TestGetMetric_NoEntries(t *testing.T) {
	b, err := NewStormStore("test.db")
	assert.NoError(t, err)

	defer os.Remove("test.db")
	defer b.Close()

	metrics, err := b.GetMetric("testMetric", time.Now())
	assert.NoError(t, err)

	var expected []Metric

	assert.Equal(t, expected, metrics)
}

func TestGetMetric_Entries(t *testing.T) {
	b, err := NewStormStore("test.db")
	assert.NoError(t, err)

	defer os.Remove("test.db")
	defer b.Close()

	testTime, err := time.Parse("2006-01-02", "2018-01-01")
	assert.NoError(t, err)

	metric := Metric{
		Name:  "testMetric",
		Date:  testTime,
		Value: 3.0,
	}
	metric.ID = fmt.Sprintf("%s-%s", metric.Date, metric.Name)

	err = b.WriteMetric(metric)
	assert.NoError(t, err)

	metrics, err := b.GetMetric("testMetric", testTime.AddDate(0, 0, -1))
	assert.NoError(t, err)

	expected := []Metric{metric}

	assert.Equal(t, expected, metrics)
}

func TestGetMetric_NoValidEntries(t *testing.T) {
	b, err := NewStormStore("test.db")
	assert.NoError(t, err)

	defer os.Remove("test.db")
	defer b.Close()

	testTime, err := time.Parse("2006-01-02", "2018-01-01")
	assert.NoError(t, err)

	metric := Metric{
		Name:  "testMetric",
		Date:  testTime,
		Value: 3.0,
	}
	metric.ID = fmt.Sprintf("%s-%s", metric.Date, metric.Name)

	err = b.WriteMetric(metric)
	assert.NoError(t, err)

	metrics, err := b.GetMetric("testMetric", time.Now())
	assert.NoError(t, err)

	var expected []Metric

	assert.Equal(t, expected, metrics)
}
