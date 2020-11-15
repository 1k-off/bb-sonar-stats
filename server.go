package main

import (
	bb "bb-sonarqube-integration/bitbucket"
	"bb-sonarqube-integration/sonarqube"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type TenantInfo struct {
	Producttype string      `json:"productType"`
	Principal   interface{} `json:"principal"`
	Eventtype   string      `json:"eventType"`
	Baseurl     string      `json:"baseUrl"`
	Publickey   string      `json:"publicKey"`
	User        interface{} `json:"user"`
	Key         string      `json:"key"`
	Baseapiurl  string      `json:"baseApiUrl"`
	Clientkey   string      `json:"clientKey"`
	Consumer    struct {
		Description string      `json:"description"`
		Links       interface{} `json:"links"`
		URL         string      `json:"url"`
		Secret      string      `json:"secret"`
		Key         string      `json:"key"`
		ID          int         `json:"id"`
		Name        string      `json:"name"`
	} `json:"consumer"`
	Sharedsecret string `json:"sharedSecret"`
}

func (c *Context) healthcheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]string{"OK"})
}

func (c *Context) atlassianConnect(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("./templates", "atlassian-connect.json")
	vals := map[string]string{
		"Organization": c.Config.Organization,
		"BaseUrl":     c.Config.BaseUrl,
		"ConsumerKey": c.Config.BitbucketOauthKey,
	}
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		c.Logger.Errorf("%v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	tmpl.ExecuteTemplate(w, "config", vals)
}

func (c *Context) installed(w http.ResponseWriter, r *http.Request) {
	c.Logger.Infoln("Received /installed call")
	ti := &TenantInfo{}
	requestContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Logger.Errorf("Can't read request: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	json.Unmarshal(requestContent, ti)

	c.Logger.Infof("Parsed /installed: %#v", ti)
	json.NewEncoder(w).Encode([]string{"OK"})
}

func (c *Context) uninstalled(w http.ResponseWriter, r *http.Request) {
	c.Logger.Infoln("Received /uninstalled call")
	json.NewEncoder(w).Encode([]string{"OK"})
}

func (c *Context) stats(w http.ResponseWriter, r *http.Request) {
	c.PrintDump(w, r, false)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	lp := path.Join(templatesPath, "panel.hbs")
	repo := strings.Split(r.URL.Query().Get("repoPath"), "/")
	repoOwner, repoName := repo[0], repo[1]
	c.SonarServerUrl, c.SonarToken, c.SonarProjectKey = c.GetSonarData(repoOwner, repoName, c.Config.RepoBranch, c.Config.SonarConfigPath)
	if c.SonarServerUrl != "" && c.SonarProjectKey != "" {
		s := sonarqube.SonarNewClient(c.SonarServerUrl, c.SonarToken, c.SonarProjectKey)
		stats, err := s.GetStats()
		if err != nil {
			c.showErrorForm(err.Error(), w)
		}
		vals := map[string]string{
			"SonarServerUrl":               c.SonarServerUrl,
			"SonarProjectKey":              c.SonarProjectKey,
			"SonarTestCoverageValue":       stats.MeasureStats.TestCoverageValue,
			"SonarLineCoverageValue":       stats.MeasureStats.LineCoverageValue,
			"SonarConditionsToCoverValue":  stats.MeasureStats.ConditionsToCoverValue,
			"SonarUnitTestSuccessValue":    stats.MeasureStats.UnitTestSuccessValue,
			"SonarUnitTestFailuresValue":   stats.MeasureStats.UnitTestFailuresValue,
			"SonarUnitTestErrorsValue":     stats.MeasureStats.UnitTestErrorsValue,
			"SonarUnitTestTestsValue":      stats.MeasureStats.UnitTestTestsValue,
			"SonarUnitTestSkippedValue":    stats.MeasureStats.UnitTestSkippedValue,
			"SonarUnitTestTimeValue":       stats.MeasureStats.UnitTestTimeValue,
			"SonarLinesOfCodeValue":        stats.MeasureStats.LinesOfCodeValue,
			"SonarFilesValue":              stats.MeasureStats.FilesValue,
			"SonarFunctionsValue":          stats.MeasureStats.FunctionsValue,
			"SonarDuplicationsValue":       stats.MeasureStats.DuplicationsValue,
			"SonarDuplicationsLinesValue":  stats.MeasureStats.DuplicationsLinesValue,
			"SonarDuplicationsBlocksValue": stats.MeasureStats.DuplicationsBlocksValue,
			"SonarDuplicationsFilesValue":  stats.MeasureStats.DuplicationsFilesValue,
			"SonarSqaleValue":              stats.MeasureStats.Sqale.Value,
			"SonarSqaleColor":              stats.MeasureStats.Sqale.Color,
			"SonarDebtRatioValue":          stats.MeasureStats.DebtRatioValue,
			"SonarDebtValue":               stats.MeasureStats.DebtRatioIndexValue,
			"SonarIssuesValue":             stats.MeasureStats.IssuesValue,
			"SonarBlockerValue":            stats.IssueStats.BlockerValue,
			"SonarCriticalValue":           stats.IssueStats.CriticalValue,
			"SonarMajorValue":              stats.IssueStats.MajorValue,
			"SonarMinorValue":              stats.IssueStats.MinorValue,
			"SonarInfoValue":               stats.IssueStats.InfoValue,
			"SonarQualityGateStatus":       stats.QualityGate.Text,
			"SonarQualityGateColor":        stats.QualityGate.Color,
		}
		tmpl, err := template.ParseFiles(lp)
		if err != nil {
			c.Logger.Errorf("%v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		tmpl.ExecuteTemplate(w, "panel", vals)
	} else {
		error := "Can't parse sonar.json file, or it is exist."
		if c.SonarToken == "" {
			error = "Can't parse sonar token, or it is not defined."
		}
		if c.SonarProjectKey == "" {
			error = "Can't parse sonar project key, or it is not defined."
		}
		c.showErrorForm(error, w)
	}

}

// routes all URL routes for app add-on
func (c *Context) routes() *mux.Router {
	r := mux.NewRouter()
	r.Path("/healthcheck").Methods("GET").HandlerFunc(c.healthcheck)
	r.Path("/").Methods("GET").HandlerFunc(c.atlassianConnect)
	r.Path("/atlassian-connect.json").Methods("GET").HandlerFunc(c.atlassianConnect)
	r.Path("/installed").Methods("POST").HandlerFunc(c.installed)
	r.Path("/uninstalled").Methods("POST").HandlerFunc(c.uninstalled)
	r.Path("/stats").Methods("GET").HandlerFunc(c.stats)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(staticPath)))
	return r
}

func (c *Context) ListenAndServe() {
	c.Logger.Infoln("Sonarqube stats bitbucket plugin running on port:", c.Config.Port)
	r := c.routes()
	http.Handle("/", r)
	http.ListenAndServe(":"+c.Config.Port, nil)
}

func (c *Context) GetSonarData(owner, repoName, branch, path string) (server, token, projectKey string) {
	var (
		sc *sonarqube.SonarConfig
	)
	bc := bb.NewBitbucketClient(c.Config.BitbucketOauthKey, c.Config.BitbucketOauthSecret)
	content, err := bc.GetFileContent(owner, repoName, branch, path)
	if err != nil {
		c.Logger.Errorf("Can't parse sonar config file. %v", err)
		return "", "", ""
	} else {
		err := json.Unmarshal([]byte(content.String()), &sc)
		if err != nil {
			c.Logger.Errorf("Error while parsing sonar config file. %v", err)
		}
		return sc.Server, sc.Token, sc.ProjectKey
	}
}

func (c *Context) showErrorForm(error string, w http.ResponseWriter) {
	vals := map[string]string {
		"ErrorMessage": error,
	}
	c.Logger.Errorln(error)

	ep := path.Join(templatesPath, "panel_error.hbs")
	tmpl, err := template.ParseFiles(ep)
	if err != nil {
		c.Logger.Errorf("%v", err)
		http.Error(w, error, 500)
		return
	}
	tmpl.ExecuteTemplate(w, "panel_error", vals)
}
