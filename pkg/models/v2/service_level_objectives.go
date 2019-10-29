package v2

type SLOComparison struct {
	CompareWith               string `yaml:"compare_with"`              // single_result|several_results
	IncludeResultWithScore    string `yaml:"include_result_with_score"` // all|pass|pass_or_warn
	NumberOfComparisonResults int    `yaml:"number_of_comparison_results"`
	AggregateFunction         string `yaml:"aggregate_function"`
}

type SLOCriteria struct {
	Criteria []string `yaml:"criteria"`
}

type SLO struct {
	SLI     string         `yaml:"sli"`
	Pass    []*SLOCriteria `yaml:"pass"`
	Warning []*SLOCriteria `yaml:"warning"`
	Weight  int            `yaml:"weight"`
	KeySLI  bool           `yaml:"key_sli"`
}

type SLOScore struct {
	Pass    string `yaml:"pass"`
	Warning string `yaml:"warning"`
}

//ServiceLevelObjectives describes SLO requirements
type ServiceLevelObjectives struct {
	SpecVersion string            `yaml:"spec_version"`
	Filter      map[string]string `yaml:"filter"`
	Comparison  *SLOComparison    `yaml:"comparison"`
	Objectives  []*SLO            `yaml:"objectives"`
	TotalScore  *SLOScore         `yaml:"total_score"`
}
