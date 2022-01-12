package models

import "encoding/json"

// DeleteServiceResponse delete service response
type DeleteServiceResponse struct {

	// message
	Message string `json:"message,omitempty"`
}

func (d *DeleteServiceResponse) ToJSON() ([]byte, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

func (d *DeleteServiceResponse) FromJSON(b []byte) error {
	var res DeleteServiceResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*d = res
	return nil
}
