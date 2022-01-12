package models

import "encoding/json"

// Service service
type Service struct {

	// Creation date of the service
	CreationDate string `json:"creationDate,omitempty"`

	// Currently deployed image
	DeployedImage string `json:"deployedImage,omitempty"`

	// last event types
	LastEventTypes map[string]EventContextInfo `json:"lastEventTypes,omitempty"`

	// open approvals
	OpenApprovals []*Approval `json:"openApprovals"`

	// Service name
	ServiceName string `json:"serviceName,omitempty"`
}

// ToJSON converts object to JSON string
func (s *Service) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

// FromJSON converts JSON string to object
func (s *Service) FromJSON(b []byte) error {
	var res Service
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*s = res
	return nil
}
