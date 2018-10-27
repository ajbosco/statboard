package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/ajbosco/statboard/pkg/config"
	"github.com/ajbosco/statboard/pkg/statboard"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	_                  Collector = &GithubCollector{}
	contributionEvents           = []string{
		"CommitCommentEvent",
		"CreateEvent",
		"RepositoryEvent",
		"IssuesEvent",
		"IssueCommentEvent",
		"PullRequestEvent",
		"PullRequestReviewEvent",
		"PullRequestReviewCommentEvent",
		"PushEvent"}
)

// GithubCollector is used to collect metrics from Github API and implements Collector interface
type GithubCollector struct {
	username string
	client   *github.Client
}

// NewGithubCollector parses config file and creates a new GithubCollector
func NewGithubCollector(cfg config.Config) (*GithubCollector, error) {
	if cfg.Github.Username == "" {
		return nil, errors.New("'github.username' must be present in config")
	}
	if cfg.Github.AccessToken == "" {
		return nil, errors.New("'github.access_token' must be present in config")
	}

	client := github.NewClient(oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Github.AccessToken})))

	return &GithubCollector{username: cfg.Github.Username, client: client}, nil
}

// Collect returns metric from Github API
func (c *GithubCollector) Collect(metricName string, daysBack int) ([]statboard.Metric, error) {
	var m []statboard.Metric
	var err error

	switch metricName {
	case "contributions":
		m, err = c.getContributions(daysBack)
	default:
		err = fmt.Errorf("unsupported metric: %s", metricName)
	}

	return m, err
}

func (c *GithubCollector) getContributions(daysBack int) ([]statboard.Metric, error) {
	end := time.Now().UTC().AddDate(0, 0, -1)
	start := end.AddDate(0, 0, -daysBack)

	metrics := generateEmptyMetrics("github.contributions", start, end)

	events, err := c.fetchEvents()
	if err != nil {
		return nil, err
	}

	metrics = aggregateEvents(events, metrics)

	return metrics, nil
}

func (c *GithubCollector) fetchEvents() ([]github.Event, error) {
	var events []github.Event

	opt := &github.ListOptions{}
	for page := 1; ; page++ {
		opt.Page = page
		events, resp, err := c.client.Activity.ListEventsPerformedByUser(context.Background(), c.username, true, opt)
		if err != nil {
			return nil, errors.Wrap(err, "collecting github events failed")
		}

		for _, event := range events {
			// Only count activity for contribution events
			if isContribEvent(*event.Type) {
				events = append(events, event)
			}
		}

		if resp.NextPage == 0 {
			break
		}
	}

	return events, nil
}

func aggregateEvents(events []github.Event, metrics []statboard.Metric) []statboard.Metric {
	for _, event := range events {
		// Only count activity for contribution events
		if isContribEvent(*event.Type) {
			for i := 0; i < len(metrics); i++ {
				met := &metrics[i]
				if event.CreatedAt.Truncate(24 * time.Hour).Equal(met.Date) {
					met.Value++
				}
			}
		}
	}

	return metrics
}

func isContribEvent(eventType string) bool {
	for _, v := range contributionEvents {
		if v == eventType {
			return true
		}
	}
	return false
}
