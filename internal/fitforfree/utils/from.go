package utils

import "time"

// From returns the start time to list for on the given day
func From(day int) time.Time {
	now := time.Now()
	if day == 0 {
		return now
	}

	t := now.AddDate(0, 0, day)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}
