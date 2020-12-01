package utils

import "time"

// To returns the end time to list for on the given day
func To(day int) time.Time {
	t := From(day)
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local)
}
