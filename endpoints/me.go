package endpoints

import api "github.com/1CRM/api-client-go"

// UserInfo is the result of calling /me endpoint
type UserInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Timezone  string `json:"timezone"`
	IsPartner bool   `json:"is_partner"`
}

// Me calls /me endpoint
func (cl *Client) Me(options ...api.RequestOption) (*UserInfo, error) {
	res, err := cl.Get("me", options...)
	if err != nil {
		return nil, err
	}
	var info UserInfo
	if err = res.ParseJSON(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
