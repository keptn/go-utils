package models

import "encoding/json"

// DeleteServiceResponse delete service response
type DeleteServiceResponse struct {

	// message
	Message string `json:"message,omitempty"`
}

// ToJSON converts object to JSON string
func (d *DeleteServiceResponse) ToJSON() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

// FromJSON converts JSON string to object
func (d *DeleteServiceResponse) FromJSON(b []byte) error {
	var res DeleteServiceResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*d = res
	return nil
}
