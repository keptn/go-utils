package models

// GitAuthCredentials stores git credentials
type GitAuthCredentials struct {

	// git remote URL
	RemoteURL string `json:"remoteURL"`

	// git user
	User string `json:"user,omitempty"`

	// https git credentials
	HttpsAuth *HttpsGitAuth `json:"https,omitempty"`

	//ssh git credentials
	SshAuth *SshGitAuth `json:"ssh,omitempty"`
}

// HttpsGitAuth stores HTTPS git credentials
type HttpsGitAuth struct {
	// Git token
	Token string `json:"token"`

	//git PEM Certificate
	Certificate string `json:"certificate,omitempty"`

	// insecure skip tls
	InsecureSkipTLS bool `json:"insecureSkipTLS"`

	// git proxy credentials
	Proxy *ProxyGitAuth `json:"proxy,omitempty"`
}

// SshGitAuth stores SSH git credentials
type SshGitAuth struct {
	// git private key
	PrivateKey string `json:"privateKey"`

	// git private key passphrase
	PrivateKeyPass string `json:"privateKeyPass,omitempty"`
}

// ProxyGitAuth stores proxy git credentials
type ProxyGitAuth struct {
	// git proxy URL
	URL string `json:"url"`

	// git proxy scheme
	Scheme string `json:"scheme"`

	// git proxy user
	User string `json:"user,omitempty"`

	// git proxy password
	Password string `json:"password,omitempty"`
}

// GitAuthCredentialsSecure stores git credentials without secure information
// model for retrieving credentials data with GET request
type GitAuthCredentialsSecure struct {
	// git remote URL
	RemoteURL string `json:"remoteURL"`

	// git user
	User string `json:"user,omitempty"`

	// https git credentials
	HttpsAuth *HttpsGitAuthSecure `json:"https,omitempty"`
}

// HttpsGitAuthSecure stores HTTPS git credentials without secure information
// model for retrieving credentials data with GET request
type HttpsGitAuthSecure struct {
	// insecure skip tls
	InsecureSkipTLS bool `json:"insecureSkipTLS"`

	// git proxy credentials
	Proxy *ProxyGitAuthSecure `json:"proxy,omitempty"`
}

// ProxyGitAuthSecure stores proxy git credentials without secure information
// model for retrieving credentials data with GET request
type ProxyGitAuthSecure struct {
	// git proxy URL
	URL string `json:"url"`

	// git proxy scheme
	Scheme string `json:"scheme"`

	// git proxy user
	User string `json:"user,omitempty"`
}
