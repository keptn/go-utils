package models

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
