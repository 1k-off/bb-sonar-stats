package sonarqube

// SonarConfig is a structire with Sonarqube server config. SonarConfig is a part of main (Sonar) struct.
type SonarConfig struct {
	Token      string `json:"token"`
	Server     string `json:"server"`
	ProjectKey string `json:"project_key"`
}

// SonarStats is a structure with stats collected from different endpoints that we need to show in bitbucket. SonarStats is a part of main (Sonar) struct.
type SonarStats struct {
	MeasureStats MeasureStats
	IssueStats   IssueStats
	QualityGate  QualityGate
}

// Endpoint is a structure that stores all endpoints that we need to connect.
type Endpoint struct {
	Component Component
	Issues    Issues
}

// Component is a structure with component endpoint (/api/measures/component) data.
type Component struct {
	Coverage struct {
		Path      string
		Method    string
		MetricKey MeasureMetricKey
	}
}

// MeasureStats is a structure for values from sonar /api/measures/component endpoint. MeasureStats is a part of SonarStats struct.
type MeasureStats struct {
	TestCoverageValue       string `json:"testCoverageValue"`
	LineCoverageValue       string `json:"lineCoverageValue"`
	ConditionsToCoverValue  string `json:"conditions_to_cover_value"`
	LinesOfCodeValue        string `json:"linesOfCodeValue"`
	FilesValue              string `json:"filesValue"`
	FunctionsValue          string `json:"functionsValue"`
	UnitTestSuccessValue    string `json:"unit_test_success_value"`
	UnitTestFailuresValue   string `json:"unit_test_failures_value"`
	UnitTestErrorsValue     string `json:"unit_test_errors_value"`
	UnitTestTestsValue      string `json:"unit_test_tests_value"`
	UnitTestSkippedValue    string `json:"unit_test_skipped_value"`
	UnitTestTimeValue       string `json:"unit_test_time_value"`
	DuplicationsValue       string `json:"duplications_value"`
	DuplicationsLinesValue  string `json:"duplications_lines_value"`
	DuplicationsBlocksValue string `json:"duplications_blocks_value"`
	DuplicationsFilesValue  string `json:"duplications_files_value"`
	Sqale                   Sqale  `json:"sqale"`
	DebtRatioValue          string `json:"debt_ratio_value"`
	DebtRatioIndexValue     string `json:"debt_ratio_index_value"`
	IssuesValue             string `json:"issues_value"`
}

// MeasureMetricKey is a structure that stores predefined metric keys for measure stats. Metric keys are defined in http.go
type MeasureMetricKey struct {
	TestCoverage       string
	LineCoverage       string
	ConditionsToCover  string
	LinesOfCode        string
	Files              string
	Functions          string
	UnitTestSuccess    string
	UnitTestFailures   string
	UnitTestErrors     string
	UnitTestTests      string
	UnitTestSkipped    string
	UnitTestTime       string
	Duplications       string
	DuplicationsLines  string
	DuplicationsBlocks string
	DuplicationsFiles  string
	Sqale              string
	DebtRatio          string
	DebtRatioIndex     string
	Issues             string
}

// Sqale is a structure for sqale metric.
type Sqale struct {
	Value string `json:"value"`
	Color string
}

// Issues is a structure with issues endpoint (/api/issues/search) data.
type Issues struct {
	Path     string
	Method   string
	Severity IssueSeverity
}

// IssueStats is a structure for values from sonar /api/issues/search endpoint. IssueStats is a part of SonarStats struct.
type IssueStats struct {
	BlockerValue  string
	CriticalValue string
	MajorValue    string
	MinorValue    string
	InfoValue     string
}

// IssueSeverity is a structure that stores predefined severity keys for issues stats. Severity keys are defined in http.go.
type IssueSeverity struct {
	Blocker  string
	Critical string
	Major    string
	Minor    string
	Info     string
}

// QualityGate is a structure for quality gate status
type QualityGate struct {
	Text  string
	Color string
}
