package lesson

// BookData contains the data
type BookData struct {
	ID string `json:"id"`
}

// BookResultData contains the result data
type BookResultData struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// BookResult contains the result
type BookResult struct {
	Data BookResultData `json:"data"`
}

// Book books
func (api *API) Book(data BookData) (*BookResult, error) {
	res := &BookResult{}
	if err := api.Request("POST", "/v0/lessons/", nil, data, res); err != nil {
		return nil, err
	}
	return res, nil
}
