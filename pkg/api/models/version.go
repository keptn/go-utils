package models

// Version version
type Version struct {

	// Branch in repository containing the resource
	Branch string `json:"branch,omitempty"`

	// Upstream respository containing the resource
	UpstreamURL string `json:"upstreamURL,omitempty"`

	// Version identifier
	Version string `json:"version,omitempty"`
}
