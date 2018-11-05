package collector

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ajbosco/reads/goodreads"
	"github.com/ajbosco/statboard/pkg/config"
	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/pkg/errors"
)

var (
	_ Collector = &GoodreadsCollector{}
)

// GoodreadsCollector is used to collect metrics from Goodreads API and implements Collector interface
type GoodreadsCollector struct {
	client *goodreads.Client
}

// NewGoodreadsCollector parses config file and creates a new GoodreadsCollector
func NewGoodreadsCollector(cfg config.Config) (*GoodreadsCollector, error) {
	if cfg.Goodreads.DeveloperKey == "" {
		return nil, errors.New("'goodreads.developer_key' must be present in config")
	}
	if cfg.Goodreads.DeveloperSecret == "" {
		return nil, errors.New("'goodreads.developer_secret' must be present in config")
	}

	// Get OAuth tokens if they are not in config.
	if cfg.Goodreads.AccessToken == "" || cfg.Goodreads.AccessSecret == "" {
		token, err := goodreads.GetAccessToken(cfg.Goodreads.DeveloperKey, cfg.Goodreads.DeveloperSecret)
		if err != nil {
			return nil, errors.Wrap(err, "fetching Goodreads access token failed")
		}
		cfg.Goodreads.AccessToken = token.Token
		cfg.Goodreads.AccessSecret = token.Secret
		cfg.Write()
	}

	goodreadsCfg := goodreads.Config{
		DeveloperKey:    cfg.Goodreads.DeveloperKey,
		DeveloperSecret: cfg.Goodreads.DeveloperSecret,
		AccessToken:     cfg.Goodreads.AccessToken,
		AccessSecret:    cfg.Goodreads.AccessSecret}

	client, err := goodreads.NewClient(&goodreadsCfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Goodreads client")
	}

	return &GoodreadsCollector{client: client}, nil
}

// Collect returns metric from Goodreads API
func (c *GoodreadsCollector) Collect(metricName string, daysBack int, granularity string) ([]statboard.Metric, error) {
	var m []statboard.Metric
	var err error

	switch metricName {
	case "books_read":
		m, err = c.getBooksRead(daysBack, granularity)
	case "pages_read":
		m, err = c.getPagesRead(daysBack, granularity)
	default:
		err = fmt.Errorf("unsupported metric: %s", metricName)
	}

	return m, err
}

func (c *GoodreadsCollector) getBooksRead(daysBack int, granularity string) ([]statboard.Metric, error) {
	end := time.Now().UTC().AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -daysBack)

	metrics := generateEmptyMetrics("goodreads.books_read", start, end, granularity)

	books, err := c.fetchBooks()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch books")
	}

	metrics, err = aggregateBooks(books, metrics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to aggregate book counts")
	}

	return metrics, nil
}

func (c *GoodreadsCollector) getPagesRead(daysBack int, granularity string) ([]statboard.Metric, error) {
	end := time.Now().UTC().AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -daysBack)

	metrics := generateEmptyMetrics("goodreads.pages_read", start, end, granularity)

	books, err := c.fetchBooks()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch books")
	}

	metrics, err = aggregatePages(books, metrics)
	if err != nil {
		return nil, errors.Wrap(err, "failed to aggregate page counts")
	}

	return metrics, nil
}

func (c *GoodreadsCollector) fetchBooks() ([]goodreads.Book, error) {

	user, err := c.client.GetCurrentUserID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Goodreads user_id")
	}

	books, err := c.client.ListShelfBooks("read", user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Goodreads books")
	}

	return books, nil
}

func aggregateBooks(books []goodreads.Book, metrics []statboard.Metric) ([]statboard.Metric, error) {
	for _, book := range books {
		for i := 0; i < len(metrics); i++ {
			met := &metrics[i]
			if book.ReadAt == "" {
				continue
			}
			// convert goodreads string date into time.Time
			readDate, err := time.Parse(time.RubyDate, book.ReadAt)
			if err != nil {
				return metrics, errors.Wrap(err, "failed to parse ReadAt date")
			}
			// get first of read month since this is a monthly metric
			readMonth := time.Date(readDate.Year(), readDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			if readMonth.Equal(met.Date) {
				met.Value++
			}
		}
	}
	return metrics, nil
}

func aggregatePages(books []goodreads.Book, metrics []statboard.Metric) ([]statboard.Metric, error) {
	for _, book := range books {
		for i := 0; i < len(metrics); i++ {
			met := &metrics[i]
			if book.ReadAt == "" || book.Pages == "" {
				continue
			}
			// convert goodreads string date into time.Time
			readDate, err := time.Parse(time.RubyDate, book.ReadAt)
			if err != nil {
				return metrics, errors.Wrap(err, "failed to parse ReadAt date")
			}
			// get first of read month since this is a monthly metric
			readMonth := time.Date(readDate.Year(), readDate.Month(), 1, 0, 0, 0, 0, time.UTC)
			// convert pages to float for metric Value
			pages, err := strconv.ParseFloat(book.Pages, 64)
			if err != nil {
				return metrics, errors.Wrap(err, "failed to convert Pages to float64")
			}
			if readMonth.Equal(met.Date) {
				met.Value = met.Value + pages
			}
		}
	}
	return metrics, nil
}
