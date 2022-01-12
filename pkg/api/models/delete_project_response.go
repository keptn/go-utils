package models

import "encoding/json"

// DeleteProjectResponse delete project response
type DeleteProjectResponse struct {

	// message
	Message string `json:"message,omitempty"`
}

// ToJSON converts object to JSON string
func (d *DeleteProjectResponse) ToJSON() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

// FromJSON converts JSON string to object
func (d *DeleteProjectResponse) FromJSON(b []byte) error {
	var res DeleteProjectResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*d = res
	return nil
}
