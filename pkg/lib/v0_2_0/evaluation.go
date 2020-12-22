package v0_2_0

const EvaluationTaskName = "evaluation"

type EvaluationTriggeredEventData struct {
	EventData
	Test       Test       `json:"test"`
	Evaluation Evaluation `json:"evaluation"`
	Deployment Deployment `json:"deployment"`
}

type Test struct {
	// Start indicates the starting timestamp of the tests
	Start string `json:"start"`
	// End indicates the end timestamp of the tests
	End string `json:"end"`
}

type Evaluation struct {
	// Start indicates the starting timestamp of the tests
	Start string `json:"start"`
	// End indicates the end timestamp of the tests
	End string `json:"end"`
}

type Deployment struct {
	// DeploymentNames gives the names of the deployments
	DeploymentNames []string `json:"deploymentNames"`
}

type EvaluationStartedEventData struct {
	EventData
}

type EvaluationStatusChangedEventData struct {
	EventData
}

type EvaluationFinishedEventData struct {
	EventData
	Evaluation EvaluationDetails `json:"evaluation"`
}

type EvaluationDetails struct {
	TimeStart        string                 `json:"timeStart"`
	TimeEnd          string                 `json:"timeEnd"`
	Result           string                 `json:"result"`
	Score            float64                `json:"score"`
	SLOFileContent   string                 `json:"sloFileContent"`
	IndicatorResults []*SLIEvaluationResult `json:"indicatorResults"`
	ComparedEvents   []string               `json:"comparedEvents"`
	// GitCommit indicates the version which should be deployed
	GitCommit string `json:"gitCommit"`
}

type SLIResult struct {
	Metric  string  `json:"metric"`
	Value   float64 `json:"value"`
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
}

type SLIEvaluationResult struct {
	Score   float64      `json:"score"`
	Value   *SLIResult   `json:"value"`
	Targets []*SLITarget `json:"targets"`
	Status  string       `json:"status"` // pass | warning | fail
}

type SLITarget struct {
	Criteria    string  `json:"criteria"`
	TargetValue float64 `json:"targetValue"`
	Violated    bool    `json:"violated"`
}
