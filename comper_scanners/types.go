package comper_scanners

import "encoding/json"

type Details struct {
	Results *Results `json:"results,omitempty"`
}

type Results struct {
	VulnerabilityID  string `json:"VulnerabilityID"`
	PkgName          string `json:"PkgName"`
	InstalledVersion string `json:"InstalledVersion"`
	FixedVersion     string `json:"FixedVersion"`
}

type Issue struct {
	Severity string  `json:"severity"`
	Details  Details `json:"details"`
}

type TrivyOutput struct {
	Root       string  `json:"root"`
	RootIssues []Issue `json:"root_issues"`
}

type Severity struct {
	Cvssv3 struct {
		BaseSeverity string `json:"base_severity"`
	} `json:"cvssv3"`
}

type Vulnerability struct {
	VulnerabilityID string   `json:"vulnerability_id"`
	PackageName     string   `json:"package_name"`
	AnalyzedVersion string   `json:"analyzed_version"`
	FixedVersions   []string `json:"fixed_versions"`
	CVE             string   `json:"CVE"`
	Severity        Severity `json:"severity"`
}

type SafetyOutputLog struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type ToolData struct {
	Tool      string          `json:"tool"`
	OutputLog json.RawMessage `json:"output_log"`
}

type JSONStructure map[string][]ToolData

type Cell struct {
	CellIndex       int    `json:"cell_index"`
	CellType        string `json:"cell_type"`
	ScrubbedContent string `json:"scrubbed_content"`
}

type IssueDetail struct {
	Description  string                 `json:"description"`
	SummaryField map[string]interface{} `json:"summary_field"`
}

type Issue2 struct {
	Code      string      `json:"code"`
	Severity  string      `json:"severity"`
	Cell      Cell        `json:"cell"`
	Location  string      `json:"location"`
	Details   IssueDetail `json:"details"`
	SubIssues []Issue2    `json:"issues"`
}

type NotebookIssue struct {
	Path   string   `json:"path"`
	Issues []Issue2 `json:"issues"`
}

type NotebookIssues struct {
	NotebookIssues []NotebookIssue `json:"notebook_issues"`
}

type DetectSecretResult struct {
	Type                  string `json:"type"`
	Filename              string `json:"filename"`
	HashedSecret          string `json:"hashed_secret"`
	IsVerified            bool   `json:"is_verified"`
	LineNumber            int    `json:"line_number"`
	VulnerabilitySeverity string `json:"vulnerability_severity"`
}

type DetectSecretOutputLog struct {
	Version     string                          `json:"version"`
	PluginsUsed []map[string]interface{}        `json:"plugins_used"`
	FiltersUsed []map[string]interface{}        `json:"filters_used"`
	Results     map[string][]DetectSecretResult `json:"results"`
	GeneratedAt string                          `json:"generated_at"`
}

type ToolOutput struct {
	Tool      string      `json:"tool"`
	OutputLog interface{} `json:"output_log"`
}

type FileOutputs map[string][]ToolOutput

type PresidioAnalyzerResult struct {
	Type                  string  `json:"Type"`
	Line                  int     `json:"Line"`
	Content               string  `json:"Content"`
	Start                 int     `json:"Start"`
	End                   int     `json:"End"`
	Score                 float64 `json:"Score"`
	VulnerabilitySeverity string  `json:"vulnerability_severity"`
}
