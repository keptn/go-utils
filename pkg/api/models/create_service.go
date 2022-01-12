package models

import "encoding/json"

// CreateService create service
type CreateService struct {

	// service name
	// Required: true
	ServiceName *string `json:"serviceName"`
}

func (c *CreateService) ToJSON() ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

func (c *CreateService) FromJSON(b []byte) error {
	var res CreateService
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*c = res
	return nil
}
