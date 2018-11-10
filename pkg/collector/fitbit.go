package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ajbosco/statboard/pkg/config"
	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/pkg/errors"
	"github.com/sajal/fitbitclient"
)

var _ Collector = &FitbitCollector{}

const (
	fitbitURI = "https://api.fitbit.com/1/user/-"
)

// FitbitCollector is used to collect metrics from Fitbit API and implements Collector interface
type FitbitCollector struct {
	baseURI string
	client  *http.Client
}

// NewFitbitCollector parses config file and creates a new FitbitCollector
func NewFitbitCollector(cfg config.Config) (*FitbitCollector, error) {
	if cfg.Fitbit.ClientID == "" {
		return nil, errors.New("'fitbit.client_id' must be present in config")
	}
	if cfg.Fitbit.ClientSecret == "" {
		return nil, errors.New("'fitbit.client_secret' must be present in config")
	}
	if cfg.Fitbit.CacheFile == "" {
		return nil, errors.New("'fitbit.cache_file' must be present in config")
	}

	clientCfg := &fitbitclient.Config{
		ClientID:     cfg.Fitbit.ClientID,
		ClientSecret: cfg.Fitbit.ClientSecret,
		Scopes:       []string{"activity"},
		CredFile:     cfg.Fitbit.CacheFile,
	}
	client, err := fitbitclient.NewFitBitClient(clientCfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not create fitbit client")
	}

	return &FitbitCollector{baseURI: fitbitURI, client: client}, nil
}

// Collect returns metric from Fitbit API
func (c *FitbitCollector) Collect(metricName string, monthsBack int) ([]statboard.Metric, error) {
	var m []statboard.Metric
	var err error

	switch metricName {
	case "steps":
		m, err = c.getSteps(monthsBack)
	default:
		err = fmt.Errorf("unsupported metric: %s", metricName)
	}

	return m, err
}

func (c *FitbitCollector) getSteps(monthsBack int) ([]statboard.Metric, error) {
	var a FitbitActivities

	// set range for which we will collect steps
	end := time.Now().AddDate(0, 0, -1)
	start := end.AddDate(0, -monthsBack, 0)

	// create metric for each month in range
	metrics := generateEmptyMetrics("fitbit.steps", start, end)

	endpoint := fmt.Sprintf("activities/steps/date/%s/%s.json", start.Format("2006-01-02"), end.Format("2006-01-02"))
	resp, err := doRequest(c.client, c.baseURI, endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "collecting steps failed")
	}

	if err = json.Unmarshal(resp, &a); err != nil {
		return nil, errors.Wrap(err, "unmarshaling steps failed")
	}

	metrics, err = aggregateSteps(a.Steps, metrics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to aggregate step counts")
	}
	return metrics, nil
}

func doRequest(client *http.Client, baseURI string, endpoint string) ([]byte, error) {
	// Create the request.
	uri := fmt.Sprintf("%s/%s", baseURI, strings.Trim(endpoint, "/"))
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("creating request to %s failed", uri))
	}

	// Do the request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("performing request to %s failed", uri))
	}
	defer resp.Body.Close()

	// Check that the response status code was OK.
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("bad response code: %d", resp.StatusCode)
	}

	// Read the body of the response.
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading the response body failed")
	}

	return b, nil
}

// aggregateSteps loops through daily step counts and aggregates them by month
func aggregateSteps(steps []FitbitSteps, metrics []statboard.Metric) ([]statboard.Metric, error) {
	for _, s := range steps {
		for i := 0; i < len(metrics); i++ {
			met := &metrics[i]
			// parse Activity Date into time.Time
			dt, err := time.Parse("2006-01-02", s.ActivityDate)
			// convert Steps string to float
			steps, err := strconv.ParseFloat(s.Steps, 64)
			if err != nil {
				return nil, errors.Wrap(err, "converting steps to float failed")
			}
			// get first of day of month for date of step activity
			stepDate := time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.UTC)
			if stepDate.Equal(met.Date) {
				met.Value = met.Value + steps
			}
		}
	}
	return metrics, nil
}
