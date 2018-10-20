package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sajal/fitbitclient"
	"github.com/spf13/viper"
)

var _ Collector = &FitbitCollector{}

const (
	fitbitURI = "https://api.fitbit.com/1/user/-"
)

type activitiesSteps struct {
	Steps []steps `json:"activities-steps"`
}

type steps struct {
	ActivityDate string `json:"dateTime"`
	Steps        string `json:"value"`
}

// FitbitCollector is used to collect metrics from Fitbit API and implements Collector interface
type FitbitCollector struct {
	baseURI string
	client  *http.Client
}

// NewFitbitCollector parses config file and creates a new FitbitCollector
func NewFitbitCollector(configFilePath string) (*FitbitCollector, error) {
	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("could not read config; filepath:%q", configFilePath))
	}

	clientID := viper.GetString("fitbit.client_id")
	if clientID == "" {
		return nil, errors.New("'fitbit.client_id' must be present in config")
	}
	clientSecret := viper.GetString("fitbit.client_secret")
	if clientSecret == "" {
		return nil, errors.New("'fitbit.client_secret' must be present in config")
	}
	cacheFile := viper.GetString("fitbit.cache_file")
	if cacheFile == "" {
		return nil, errors.New("'fitbit.cache_file' must be present in config")
	}

	cfg := &fitbitclient.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"activity"},
		CredFile:     cacheFile,
	}
	client, err := fitbitclient.NewFitBitClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not create fitbit client")
	}

	return &FitbitCollector{baseURI: fitbitURI, client: client}, nil
}

// Collect returns metric from Fitbit API
func (c *FitbitCollector) Collect(metricName string, daysBack int) ([]Metric, error) {
	var m []Metric
	var err error

	switch metricName {
	case "steps":
		m, err = getSteps(c.client, c.baseURI, daysBack)
	default:
		err = fmt.Errorf("Unsupported metric: %s", metricName)
	}

	return m, err
}

func getSteps(client *http.Client, baseURI string, daysBack int) ([]Metric, error) {
	var a activitiesSteps
	var m []Metric

	end := time.Now().AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -daysBack)

	endpoint := fmt.Sprintf("activities/steps/date/%s/%s.json", start.Format("2006-01-02"), end.Format("2006-01-02"))
	resp, err := doRequest(client, baseURI, endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "collecting steps failed")
	}

	if err = json.Unmarshal(resp, &a); err != nil {
		return nil, errors.Wrap(err, "unmarshaling steps failed")
	}

	for _, s := range a.Steps {
		dt, err := time.Parse("2006-01-02", s.ActivityDate)
		if err != nil {
			return nil, errors.Wrap(err, "parsing activity date failed")
		}
		v, err := strconv.ParseFloat(s.Steps, 64)
		if err != nil {
			return nil, errors.Wrap(err, "converting steps to float failed")
		}

		met := Metric{
			Name:  "fitbit.steps",
			Date:  dt,
			Value: v,
		}

		m = append(m, met)
	}

	return m, nil
}

func doRequest(client *http.Client, baseURI string, endpoint string) ([]byte, error) {
	// Create the request.
	uri := fmt.Sprintf("%s/%s", baseURI, strings.Trim(endpoint, "/"))
	fmt.Println(uri)
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
