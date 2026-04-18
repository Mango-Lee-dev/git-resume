package export

import (
	"time"

	"github.com/wootaiklee/git-resume/pkg/models"
)

// getDateRange returns the min and max dates from results
func getDateRange(results []models.AnalysisResult) (time.Time, time.Time) {
	if len(results) == 0 {
		return time.Time{}, time.Time{}
	}

	minDate, maxDate := results[0].Date, results[0].Date
	for _, r := range results {
		if r.Date.Before(minDate) {
			minDate = r.Date
		}
		if r.Date.After(maxDate) {
			maxDate = r.Date
		}
	}

	return minDate, maxDate
}
