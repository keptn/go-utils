package models

import "encoding/json"

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

// ToJSON converts object to JSON string
func (p *GitAuthCredentials) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *GitAuthCredentials) FromJSON(b []byte) error {
	var res GitAuthCredentials
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
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

// ToJSON converts object to JSON string
func (p *HttpsGitAuth) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *HttpsGitAuth) FromJSON(b []byte) error {
	var res HttpsGitAuth
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
}

// SshGitAuth stores SSH git credentials
type SshGitAuth struct {
	// git private key
	PrivateKey string `json:"privateKey"`

	// git private key passphrase
	PrivateKeyPass string `json:"privateKeyPass,omitempty"`
}

// ToJSON converts object to JSON string
func (p *SshGitAuth) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *SshGitAuth) FromJSON(b []byte) error {
	var res SshGitAuth
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
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

// ToJSON converts object to JSON string
func (p *ProxyGitAuth) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *ProxyGitAuth) FromJSON(b []byte) error {
	var res ProxyGitAuth
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
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

// ToJSON converts object to JSON string
func (p *GitAuthCredentialsSecure) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *GitAuthCredentialsSecure) FromJSON(b []byte) error {
	var res GitAuthCredentialsSecure
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
}

// HttpsGitAuthSecure stores HTTPS git credentials without secure information
// model for retrieving credentials data with GET request
type HttpsGitAuthSecure struct {
	// insecure skip tls
	InsecureSkipTLS bool `json:"insecureSkipTLS"`

	// git proxy credentials
	Proxy *ProxyGitAuthSecure `json:"proxy,omitempty"`
}

// ToJSON converts object to JSON string
func (p *HttpsGitAuthSecure) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *HttpsGitAuthSecure) FromJSON(b []byte) error {
	var res HttpsGitAuthSecure
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
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

// ToJSON converts object to JSON string
func (p *ProxyGitAuthSecure) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// FromJSON converts JSON string to object
func (p *ProxyGitAuthSecure) FromJSON(b []byte) error {
	var res ProxyGitAuthSecure
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}
	*p = res
	return nil
}
