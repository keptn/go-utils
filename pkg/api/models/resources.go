package models

import "encoding/json"

// Resources resources
type Resources struct {

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// resources
	Resources []*Resource `json:"resources"`

	// Total number of resources
	TotalCount float64 `json:"totalCount,omitempty"`
}

// ToJSON converts object to JSON string
func (r *Resources) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON converts JSON string to object
func (r *Resources) FromJSON(b []byte) error {
	var res Resources
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*r = res
	return nil
}
