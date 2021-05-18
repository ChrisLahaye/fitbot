package utils

import "time"

// Date returns the time to list for on the given day
func Date(day int) time.Time {
	now := time.Now()
	t := now.AddDate(0, 0, day)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}
