package lesson

// BookData contains the data
type BookData struct {
	ID string `json:"id"`
}

// BookResult contains the result
type BookResult struct {
	Data Booking `json:"data"`
}

// Book books
func (api *API) Book(data BookData) (*BookResult, error) {
	res := &BookResult{}
	if err := api.Request("POST", "/v0/lessons/", nil, data, res); err != nil {
		return nil, err
	}
	return res, nil
}
