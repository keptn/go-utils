package models

// Metadata metadata
type Metadata struct {

	// bridgeversion
	Bridgeversion string `json:"bridgeversion,omitempty"`

	// keptnlabel
	Keptnlabel string `json:"keptnlabel,omitempty"`

	// keptnservices
	Keptnservices interface{} `json:"keptnservices,omitempty"`

	// keptnversion
	Keptnversion string `json:"keptnversion,omitempty"`

	// namespace
	Namespace string `json:"namespace,omitempty"`
}
