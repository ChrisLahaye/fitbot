package scheduled_classes

import (
	"time"
)

// ListQuery contains the query
type ListQuery struct {
	Venues             string    `url:"venues"`
	ExcludeFullyBooked bool      `url:"exclude_fully_booked"`
	Date               time.Time `url:"date" layout:"2006-01-02T15:04:05.000Z07:00"`
	Category           string    `url:"category"` // all | free_practise | outdoor_lesson
}

// ListResult contains the result
type ListResult struct {
	ScheduledClasses []ScheduledClass `json:"scheduled_classes"`
}

// List lists
func (api *API) List(query ListQuery) (*ListResult, error) {
	res := &ListResult{}
	if err := api.Request("GET", "/v1/scheduled_classes", query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}
