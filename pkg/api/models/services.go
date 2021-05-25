package models

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
