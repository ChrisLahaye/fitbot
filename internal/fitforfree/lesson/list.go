package lesson

import (
	"fmt"
	"time"
)

// ListParams contains the params
type ListParams struct {
	Venue string
	From  time.Time
	To    time.Time
}

// ListResult contains the result
type ListResult struct {
	Data struct {
		Lessons []Lesson `json:"lessons"`
	} `json:"data"`
}

// List lists
func (api *API) List(params ListParams) (*ListResult, error) {
	res := &ListResult{}
	if err := api.Request("GET", "/v0/lessons", listQuery{
		Venues: fmt.Sprintf("[\"%s\"]", params.Venue),
		From:   params.From.Unix(),
		To:     params.To.Unix(),
	}, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

type listQuery struct {
	Venues string `url:"venues"`
	From   int64  `url:"from"`
	To     int64  `url:"to"`
}
