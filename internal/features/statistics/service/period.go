package statistics_service

import "time"

func inPeriod(t time.Time, from, to *time.Time) bool {
	if from != nil && t.Before(*from) {
		return false
	}
	if to != nil && !t.Before(to.AddDate(0, 0, 1)) {
		return false
	}
	return true
}
