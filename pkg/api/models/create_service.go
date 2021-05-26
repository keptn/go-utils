package models

// CreateService create service
type CreateService struct {

	// service name
	// Required: true
	ServiceName *string `json:"serviceName"`
}
