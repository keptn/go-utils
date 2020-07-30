package v0_2_0

///// v0.2.0 Shipyard Spec ///////

// Shipyard describes a shipyard specification according to Keptn spec 0.2.0
type Shipyard struct {
	ApiVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   ShipyardMetadata `json:"metadata" yaml:"metadata"`
	Spec       ShipyardSpec     `json:"spec" yaml:"spec"`
}

// ShipyardMetadata describes Shipyayrd metadata
type ShipyardMetadata struct {
	Name string `json:"name" yaml:"name"`
}

// ShipyardSpec consists of any number of stages
type ShipyardSpec struct {
	Stages []Stage `json:"stages" yaml:"stages"`
}

// Stage defines a stage by its name and list of workflows
type Stage struct {
	Name      string     `json:"name" yaml:"name"`
	Workflows []Workflow `json:"workflow" yaml:"workflow"`
}

// Workflow defines a workflow by its name and tasks. The triggers property is optional
type Workflow struct {
	Name     string   `json:"name" yaml:"name"`
	Triggers []string `json:"triggers" yaml:"triggers"`
	Tasks    []Task   `json:"tasks" yaml:"tasks"`
}

// Task defines a task by its name and optional properties
type Task struct {
	Name       string      `json:"name" yaml:"name"`
	Properties interface{} `json:"properties" yaml:"properties"`
}
