package models

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
