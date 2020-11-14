package sonarqube

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// SonarHTTPClient is a structure that extends http.Client. This need for creating methods of http client.
type SonarHTTPClient struct {
	http.Client
	Config SonarHTTPClientConfig
}

// SonarHTTPClientConfig is a structure with sonarqube config data that need for connecting sonarqube server.
type SonarHTTPClientConfig struct {
	Token  string
	Server string
}

// NewHTTPClient is a wrapper for http client. Here we configuring http client once to use it in the application.
func NewHTTPClient(server, token string) *SonarHTTPClient {
	client := &SonarHTTPClient{
		Config: SonarHTTPClientConfig{
			Token:  token,
			Server: server,
		},
	}
	return client
}

// NewReq is a wrapper for http NewRequest method to set headers once
func (c *SonarHTTPClient) NewHTTPRequest(method, path string, body io.Reader) *http.Request {
	url := c.Config.Server + path
	req, _ := http.NewRequest(method, url, body)
	if c.Config.Token != "" {
		sonarAuthString := Base64Encode(c.Config.Token + ":")
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", sonarAuthString))
	}
	return req
}

// SonarEndpoint is a function to fill all metric keys and endpoints that will be used in future to collect sonarqube project stats.
// TODO: Unit tests metrics
func SonarEndpoint() *Endpoint {
	// Define measure metric keys
	mmk := MeasureMetricKey{
		TestCoverage:       "coverage",
		LineCoverage:       "line_coverage",
		ConditionsToCover:  "new_conditions_to_cover",
		LinesOfCode:        "ncloc",
		Files:              "files",
		Functions:          "functions",
		UnitTestSuccess:    "",
		UnitTestFailures:   "",
		UnitTestErrors:     "",
		UnitTestTests:      "",
		UnitTestSkipped:    "",
		UnitTestTime:       "",
		Duplications:       "duplicated_lines_density",
		DuplicationsLines:  "duplicated_lines",
		DuplicationsBlocks: "duplicated_blocks",
		DuplicationsFiles:  "duplicated_files",
		Sqale:              "sqale_rating",
		DebtRatio:          "sqale_debt_ratio",
		DebtRatioIndex:     "sqale_index",
		Issues:             "violations",
	}
	// Define issues severity
	is := IssueSeverity{
		Blocker:  "BLOCKER",
		Critical: "CRITICAL",
		Major:    "MAJOR",
		Minor:    "MINOR",
		Info:     "INFO",
	}
	c := Component{Coverage: struct {
		Path      string
		Method    string
		MetricKey MeasureMetricKey
	}{
		Path:      "/api/measures/component",
		Method:    http.MethodGet,
		MetricKey: mmk},
	}

	i := Issues{
		Path:     "/api/issues/search",
		Method:   http.MethodGet,
		Severity: is,
	}

	return &Endpoint{
		Component: c,
		Issues:    i,
	}
}

// AddQueryParams is a function to add query params to url.
func AddQueryParams(path, projectKey, metricKey string) string {
	pathStr, _ := url.Parse(path)
	query, _ := url.ParseQuery(pathStr.RawQuery)
	query.Add("component", projectKey)
	query.Add("metricKeys", metricKey)
	pathStr.RawQuery = query.Encode()
	return pathStr.String()
}

// AddQueryParamsIssues is a function to add query params to issues url.
func AddQueryParamsIssues(path, projectKey, severity string) string {
	pathStr, _ := url.Parse(path)
	query, _ := url.ParseQuery(pathStr.RawQuery)
	query.Add("componentKeys", projectKey)
	query.Add("severities", severity)
	query.Add("resolved", "false")
	pathStr.RawQuery = query.Encode()
	return pathStr.String()
}
