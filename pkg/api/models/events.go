package models

import "encoding/json"

// Events events
type Events struct {

	// events
	Events []*KeptnContextExtendedCE `json:"events"`

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// Total number of resources
	TotalCount float64 `json:"totalCount,omitempty"`
}

// ToJSON converts object to JSON string
func (e *Events) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON converts JSON string to object
func (e *Events) FromJSON(b []byte) error {
	var res Events
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*e = res
	return nil
}
