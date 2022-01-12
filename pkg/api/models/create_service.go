package models

import "encoding/json"

// CreateService create service
type CreateService struct {

	// service name
	// Required: true
	ServiceName *string `json:"serviceName"`
}

// ToJSON converts object to JSON string
func (c *CreateService) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// FromJSON converts JSON string to object
func (c *CreateService) FromJSON(b []byte) error {
	var res CreateService
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*c = res
	return nil
}
