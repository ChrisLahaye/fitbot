package login

// LoginData contains the data
type LoginData struct {
	MemberID      string `json:"memberid"`
	Postcode      string `json:"postcode"`
	TermsAccepted bool   `json:"terms_accepted"`
}

// LoginResultData contains the result data
type LoginResultData struct {
	FirstName string `json:"firstname"`
	SessionID string `json:"sessionid"`
	Surname   string `json:"surname"`
}

// LoginResult contains the result
type LoginResult struct {
	Data LoginResultData `json:"data"`
}

// Login logs in
func (api *API) Login(data LoginData) (*LoginResult, error) {
	res := &LoginResult{}
	if err := api.Request("POST", "/v0/login/", nil, data, res); err != nil {
		return nil, err
	}
	return res, nil
}
