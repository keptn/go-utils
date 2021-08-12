package v0_2_0

type WebHookConfig struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Spec       WebHookConfigSpec `yaml:"spec"`
}

type WebHookConfigSpec struct {
	Webhooks []Webhook `yaml:"webhooks"`
}

type Webhook struct {
	Type     string    `yaml:"type"`
	EnvFrom  []EnvFrom `yaml:"envFrom"`
	Requests []string  `yaml:"requests"`
}

type EnvFrom struct {
	SecretRef WebHookSecretRef `yaml:"secretRef"`
	Name      string           `yaml:"name"`
}

type WebHookSecretRef struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
}
