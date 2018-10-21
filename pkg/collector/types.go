package collector

// FitbitActivities defines the struct for the Fitbit activies object
type FitbitActivities struct {
	Steps []FitbitSteps `json:"activities-steps"`
}

// FitbitSteps defines the struct for the Fitbit steps object
type FitbitSteps struct {
	ActivityDate string `json:"dateTime"`
	Steps        string `json:"value"`
}
