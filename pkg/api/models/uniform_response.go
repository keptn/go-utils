package models

import "encoding/json"

type RegisterIntegrationResponse struct {
	ID string `json:"id"`
}

// ToJSON converts object to JSON string
func (i *RegisterIntegrationResponse) ToJSON() ([]byte, error) {
	return json.Marshal(i)
}

// FromJSON converts JSON string to object
func (i *RegisterIntegrationResponse) FromJSON(b []byte) error {
	var res RegisterIntegrationResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*i = res
	return nil
}

type CreateSubscriptionResponse struct {
	ID string `json:"id"`
}

// ToJSON converts object to JSON string
func (s *CreateSubscriptionResponse) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON converts JSON string to object
func (s *CreateSubscriptionResponse) FromJSON(b []byte) error {
	var res CreateSubscriptionResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*s = res
	return nil
}
