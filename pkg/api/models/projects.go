package models

import "encoding/json"

// Projects projects
type Projects struct {

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// projects
	Projects []*Project `json:"projects"`

	// Total number of projects
	TotalCount float64 `json:"totalCount,omitempty"`
}

func (p *Projects) ToJSON() ([]byte, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

func (p *Projects) FromJSON(b []byte) error {
	var res Projects
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
}
