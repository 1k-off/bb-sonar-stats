package sonarqube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Sonar is a main structure with all necessary data.
type Sonar struct {
	Config       *SonarConfig
	Endpoint     *Endpoint
	Client       *SonarHTTPClient
	ClientConfig SonarHTTPClientConfig
	Stats        SonarStats
}

// SonarNewClient creates new instance of http client to connect with sonarqube server.
func SonarNewClient(server, token, projectKey string) *Sonar {
	e := SonarEndpoint()
	c := NewHTTPClient(server, token)
	sc := &SonarConfig{
		Server:     server,
		Token:      token,
		ProjectKey: projectKey,
	}

	return &Sonar{
		Config:   sc,
		Endpoint: e,
		Client:   c,
	}
}

// GetStats is a function to collect all stats from different endpoints to one structure.
func (s *Sonar) GetStats() (SonarStats, error){
	if err := s.checkIfProjectExist(); err != nil {
		return SonarStats{}, err
	}
	// Measure stats
	s.Stats.MeasureStats.TestCoverageValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.TestCoverage)
	s.Stats.MeasureStats.LineCoverageValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.LineCoverage)
	s.Stats.MeasureStats.ConditionsToCoverValue = s.GetPossibleEmptyStatValue(s.Endpoint.Component.Coverage.MetricKey.ConditionsToCover)
	s.Stats.MeasureStats.LinesOfCodeValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.LinesOfCode)
	s.Stats.MeasureStats.FilesValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.Files)
	s.Stats.MeasureStats.FunctionsValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.Functions)
	s.Stats.MeasureStats.UnitTestSuccessValue = s.GetPossibleEmptyStatValue(s.Endpoint.Component.Coverage.MetricKey.UnitTestSuccess)
	s.Stats.MeasureStats.UnitTestFailuresValue = s.GetPossibleEmptyStatValue(s.Endpoint.Component.Coverage.MetricKey.UnitTestFailures)
	s.Stats.MeasureStats.UnitTestErrorsValue = s.GetPossibleEmptyStatValue(s.Endpoint.Component.Coverage.MetricKey.UnitTestErrors)
	s.Stats.MeasureStats.UnitTestTestsValue = s.GetPossibleEmptyStatValue(s.Endpoint.Component.Coverage.MetricKey.UnitTestTests)
	s.Stats.MeasureStats.UnitTestSkippedValue = s.GetPossibleEmptyStatValue(s.Endpoint.Component.Coverage.MetricKey.UnitTestSkipped)
	s.Stats.MeasureStats.UnitTestTimeValue = s.GetPossibleEmptyStatValue(s.Endpoint.Component.Coverage.MetricKey.UnitTestTime)
	s.Stats.MeasureStats.DuplicationsValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.Duplications)
	s.Stats.MeasureStats.DuplicationsLinesValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.DuplicationsLines)
	s.Stats.MeasureStats.DuplicationsBlocksValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.DuplicationsBlocks)
	s.Stats.MeasureStats.DuplicationsFilesValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.DuplicationsFiles)
	s.Stats.MeasureStats.DebtRatioValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.DebtRatio)
	s.Stats.MeasureStats.DebtRatioIndexValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.DebtRatioIndex)
	s.Stats.MeasureStats.IssuesValue = s.GetMeasureStats(s.Endpoint.Component.Coverage.MetricKey.Issues)
	s.Stats.MeasureStats.Sqale.Value, s.Stats.MeasureStats.Sqale.Color = s.GetSqaleValue(s.Stats.MeasureStats.DebtRatioValue)
	// Issues
	s.Stats.IssueStats.BlockerValue = s.GetIssuesCountBySeverity(s.Endpoint.Issues.Severity.Blocker)
	s.Stats.IssueStats.CriticalValue = s.GetIssuesCountBySeverity(s.Endpoint.Issues.Severity.Critical)
	s.Stats.IssueStats.MajorValue = s.GetIssuesCountBySeverity(s.Endpoint.Issues.Severity.Major)
	s.Stats.IssueStats.MinorValue = s.GetIssuesCountBySeverity(s.Endpoint.Issues.Severity.Minor)
	s.Stats.IssueStats.InfoValue = s.GetIssuesCountBySeverity(s.Endpoint.Issues.Severity.Info)
	s.Stats.QualityGate.Text, s.Stats.QualityGate.Color = s.GetQualityGateStatus()

	// Idea for dynamic variables. But we haven't dynamic variables in go.

	//vMetricKey := reflect.ValueOf(s.Endpoint.Component.Coverage.MetricKey)
	//vMeasureStats := reflect.ValueOf(s.Stats.MeasureStats)
	//metricKeyValues := make([]interface{}, vMetricKey.NumField())
	//measureStatsValues := make([]interface{}, vMeasureStats.NumField())
	//for i := 0; i < vMetricKey.NumField(); i++ {
	//	metricKeyValues[i] = vMetricKey.Field(i).Interface()
	//	measureStatsValues[i] = s.GetMeasureStats(projectKey, metricKeyValues[i].(string))
	//}

	return s.Stats, nil
}

func (s *Sonar) checkIfProjectExist() error {
	type ProjectResponce struct {
		Components []struct {
			Key              string `json:"key"`
		} `json:"components"`
	}
	var (
		response ProjectResponce
	)
	path := "/api/projects/search?projects=" + s.Config.ProjectKey
	req := s.Client.NewHTTPRequest(http.MethodGet, path, nil)
	resp, err := s.Client.Do(req)
	if resp.StatusCode == 401 {
		return errors.New(fmt.Sprintf("Can't pass authorization to server %s with provided token", s.Config.Server))
	}
	if err != nil {
		return err
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		_ = json.Unmarshal(body, &response)
		if len(response.Components) == 0 {
			return errors.New(fmt.Sprintf("Project %s is not exist", s.Config.ProjectKey))
		} else {
			return nil
		}
	}
}

// GetMeasureStats is a function to collect measure stats (/api/measures/component endpoint) from sonarqube server.
func (s *Sonar) GetMeasureStats(metricKey string) string {
	type ComponentResponse struct {
		Component struct {
			Measures []struct {
				Metric string `json:"metric"`
				Value  string `json:"value"`
			} `json:"measures"`
		} `json:"component"`
	}

	var (
		response ComponentResponse
		result   string = ""
	)

	path := AddQueryParams(s.Endpoint.Component.Coverage.Path, s.Config.ProjectKey, metricKey)
	req := s.Client.NewHTTPRequest(s.Endpoint.Component.Coverage.Method, path, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		HandleError(err)
		return result
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		_ = json.Unmarshal(body, &response)
		for _, m := range response.Component.Measures {
			if m.Metric == metricKey {
				result = m.Value
			}
		}
		return result
	}
}

// GetPossibleEmptyStatValue is a function to process empty-string responce from server. It adds dash to result stat string to make stat form view better.
func (s *Sonar) GetPossibleEmptyStatValue(metricKey string) string {
	if res := s.GetMeasureStats(metricKey); res != "" {
		return res
	} else {
		return "-"
	}
}

// GetIssuesCountBySeverity is a function to get sonar project issues count by issue severity
func (s *Sonar) GetIssuesCountBySeverity(severity string) string {
	type IssuesResponce struct {
		Total int `json:"total"`
	}
	var (
		response IssuesResponce
	)
	path := AddQueryParamsIssues(s.Endpoint.Issues.Path, s.Config.ProjectKey, severity)
	req := s.Client.NewHTTPRequest(s.Endpoint.Issues.Method, path, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		HandleError(err)
		return ""
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		_ = json.Unmarshal(body, &response)
		return strconv.Itoa(response.Total)
	}
}

// GetSqaleValue is a function to get sonar project sqale rating based on debt ratio
func (s *Sonar) GetSqaleValue(debtRatioStr string) (string, string) {
	/*
		https://docs.sonarqube.org/latest/user-guide/metric-definitions/
		(Formerly the SQALE rating.) Rating given to your project related to the value of your Technical Debt Ratio. The default Maintainability Rating grid is:
		A=0-0.05, B=0.06-0.1, C=0.11-0.20, D=0.21-0.5, E=0.51-1
	*/

	// SqaleColor is a struct with predefined colors for sqale rating.
	type SqaleColor struct {
		A string
		B string
		C string
		D string
		E string
	}
	var (
		Amin, Amax, Bmin, Bmax, Cmin, Cmax, Dmin, Dmax float64 = 0, 5, 6, 10, 11, 20, 21, 50
	)

	// Define sqale colors
	sqaleColor := SqaleColor{
		A: "sqale-A",
		B: "sqale-B",
		C: "sqale-C",
		D: "sqale-D",
		E: "sqale-E",
	}

	debtRatio, err := strconv.ParseFloat(debtRatioStr, 64)
	if err != nil {
		HandleError(err)
		return sqaleColor.E, "-"
	}
	if IsBetween(debtRatio, Amin, Amax) {
		return "A", sqaleColor.A
	} else if IsBetween(debtRatio, Bmin, Bmax) {
		return "B", sqaleColor.B
	} else if IsBetween(debtRatio, Cmin, Cmax) {
		return "C", sqaleColor.C
	} else if IsBetween(debtRatio, Dmin, Dmax) {
		return "D", sqaleColor.D
	} else {
		return "E", sqaleColor.E
	}
}

// GetQualityGateStatus is a function to get project quality gate status and generate text for bitbucket form
func (s *Sonar) GetQualityGateStatus() (string, string) {
	type QualityGateResponce struct {
		ProjectStatus struct {
			Status string `json:"status"`
		} `json:"projectStatus"`
	}
	var (
		statusOk, statusFailed string = "OK", "ERROR"
		path                   string = "/api/qualitygates/project_status?projectKey=" + s.Config.ProjectKey
		response               QualityGateResponce
	)

	req := s.Client.NewHTTPRequest(s.Endpoint.Issues.Method, path, nil)
	resp, err := s.Client.Do(req)
	if err != nil {
		HandleError(err)
		return "", ""
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		_ = json.Unmarshal(body, &response)
		if response.ProjectStatus.Status == statusOk {
			return "The project has passed the quality gate.", statusOk
		} else {
			return "The project has not passed the quality gate.", statusFailed
		}
	}
}
