package models

import "encoding/json"

// Services services
type Services struct {

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// services
	Services []*Service `json:"services"`

	// Total number of services
	TotalCount float64 `json:"totalCount,omitempty"`
}

// ToJSON converts object to JSON string
func (s *Services) ToJSON() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// FromJSON converts JSON string to object
func (s *Services) FromJSON(b []byte) error {
	var res Services
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*s = res
	return nil
}
