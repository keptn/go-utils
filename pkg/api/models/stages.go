package models

// Stages stages
type Stages struct {

	// Pointer to next page, base64 encoded
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of returned page
	PageSize float64 `json:"pageSize,omitempty"`

	// stages
	Stages []*Stage `json:"stages"`

	// Total number of stages
	TotalCount float64 `json:"totalCount,omitempty"`
}
