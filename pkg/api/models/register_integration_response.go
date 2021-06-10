package models

import "encoding/json"

type RegisterIntegrationResponse struct {
	ID string `json:"id"`
}

func (i *RegisterIntegrationResponse) ToJSON() ([]byte, error) {
	if i == nil {
		return nil, nil
	}
	return json.Marshal(i)
}

func (i *RegisterIntegrationResponse) FromJSON(b []byte) error {
	var res RegisterIntegrationResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*i = res
	return nil
}
