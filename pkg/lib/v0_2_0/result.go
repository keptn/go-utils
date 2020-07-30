package keptn

type ResultType string

const (
	ResultPass    ResultType = "pass"
	ResultWarning ResultType = "warning"
	ResultFailed  ResultType = "fail"
)
