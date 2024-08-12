package comper_scanners

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

const (
	safetyTool       = "Safety"
	detectSecretTool = "Detect-Secret"
	presidioAnalyzer = "Presidio-Analyzer"
)

func WriteOutputs(outputFile1, outputFile2, outputFile3, outputFile4 string, trivyOutput TrivyOutput, data JSONStructure, notebookIssues NotebookIssues, fileOutputs FileOutputs) error {
	err := WriteOutput(outputFile1, func(writer *tabwriter.Writer) {
		WriteTrivyOutput(writer, trivyOutput)
	})
	if err != nil {
		return fmt.Errorf("error writing Trivy output: %w", err)
	}

	err = WriteOutput(outputFile2, func(writer *tabwriter.Writer) {
		WriteSafetyOutput(writer, data)
	})
	if err != nil {
		return fmt.Errorf("error writing Safety output: %w", err)
	}

	err = WriteOutput(outputFile3, func(writer *tabwriter.Writer) {
		WriteNotebookIssues(writer, notebookIssues)
	})
	if err != nil {
		return fmt.Errorf("error writing Notebook Issues output: %w", err)
	}

	err = WriteOutput(outputFile4, func(writer *tabwriter.Writer) {
		WriteFileOutputs(writer, fileOutputs)
	})
	if err != nil {
		return fmt.Errorf("error writing File Outputs: %w", err)
	}

	return nil
}

func WriteTrivyOutput(writer *tabwriter.Writer, trivyOutput TrivyOutput) {
	fmt.Fprintln(writer, "Severity\tVulnerability ID\tPkg Name\tInstalled Version\tFixed Version")
	fmt.Fprintln(writer, "---------\t----------------\t--------\t----------------\t-------------")

	for _, issue := range trivyOutput.RootIssues {
		if details := issue.Details; details.Results != nil {
			result := details.Results
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
				issue.Severity,
				result.VulnerabilityID,
				result.PkgName,
				result.InstalledVersion,
				result.FixedVersion,
			)
		}
	}
}

func WriteSafetyOutput(writer *tabwriter.Writer, data JSONStructure) {
	fmt.Fprintln(writer, "Severity\tVulnerability ID\tPkg Name\tInstalled Version\tFixed Version")
	fmt.Fprintln(writer, "---------\t----------------\t--------\t----------------\t-------------")

	// Create a map to track printed lines
	printedLines := make(map[string]bool)

	for _, tools := range data {
		for _, tool := range tools {
			if tool.Tool == safetyTool {
				var outputLog SafetyOutputLog
				if err := json.Unmarshal(tool.OutputLog, &outputLog); err != nil {
					continue
				}

				for _, vuln := range outputLog.Vulnerabilities {
					fixedVersions := strings.Join(vuln.FixedVersions, ", ")
					line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n",
						vuln.Severity.Cvssv3.BaseSeverity,
						vuln.CVE,
						vuln.PackageName,
						vuln.AnalyzedVersion,
						fixedVersions,
					)

					// Check if the line has been printed before
					if !printedLines[line] {
						fmt.Fprint(writer, line)
						printedLines[line] = true
					}
				}
			}
		}
	}
}

func WriteNotebookIssues(writer *tabwriter.Writer, notebookIssues NotebookIssues) {
	fmt.Fprintln(writer, "File name\tCode\tSeverity\tDescription\tSummary Field")
	fmt.Fprintln(writer, "----------\t----\t--------\t-----------\t-------------")

	for _, notebook := range notebookIssues.NotebookIssues {
		for _, issue := range notebook.Issues {
			description := issue.Details.Description
			summaryField, _ := json.Marshal(issue.Details.SummaryField)
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
				filepath.Base(notebook.Path),
				issue.Code,
				issue.Severity,
				description,
				string(summaryField),
			)

			for _, subIssue := range issue.SubIssues {
				subDescription := subIssue.Details.Description
				subSummaryField, _ := json.Marshal(subIssue.Details.SummaryField)
				fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
					filepath.Base(notebook.Path),
					subIssue.Code,
					subIssue.Severity,
					subDescription,
					string(subSummaryField),
				)
			}
		}
	}
}

func WriteFileOutputs(writer *tabwriter.Writer, fileOutputs FileOutputs) {
	fmt.Fprintln(writer, "File name\tCode\tSeverity\tDescription\tSummary Field")
	fmt.Fprintln(writer, "----------\t----\t--------\t-----------\t-------------")

	for fileName, toolOutputs := range fileOutputs {
		for _, toolOutput := range toolOutputs {
			switch toolOutput.Tool {
			case detectSecretTool:
				WriteDetectSecretOutput(writer, toolOutput.OutputLog)
			case presidioAnalyzer:
				WritePresidioAnalyzerOutput(writer, toolOutput.OutputLog, fileName)
			}
		}
	}
}

func WriteDetectSecretOutput(writer *tabwriter.Writer, outputLog interface{}) {
	var detectSecretOutputLog DetectSecretOutputLog
	data, _ := json.Marshal(outputLog)
	json.Unmarshal(data, &detectSecretOutputLog)

	for _, results := range detectSecretOutputLog.Results {
		for _, result := range results {
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
				filepath.Base(result.Filename),
				detectSecretTool,
				result.VulnerabilitySeverity,
				result.Type,
				"null",
			)
		}
	}
}

func WritePresidioAnalyzerOutput(writer *tabwriter.Writer, outputLog interface{}, fileName string) {
	if outputLogStr, ok := outputLog.(string); ok {
		lines := strings.Split(outputLogStr, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			var presidioResult PresidioAnalyzerResult
			fields := strings.Split(line, ", ")
			for _, field := range fields {
				kv := strings.Split(field, ": ")
				if len(kv) != 2 {
					continue
				}
				key, value := kv[0], kv[1]
				switch key {
				case "Type":
					presidioResult.Type = value
				case "Line":
					fmt.Sscanf(value, "%d", &presidioResult.Line)
				case "Start":
					fmt.Sscanf(value, "%d", &presidioResult.Start)
				case "End":
					fmt.Sscanf(value, "%d", &presidioResult.End)
				case "Score":
					fmt.Sscanf(value, "%f", &presidioResult.Score)
				case "vulnerability_severity":
					presidioResult.VulnerabilitySeverity = value
				}
			}
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n",
				filepath.Base(fileName),
				presidioAnalyzer,
				presidioResult.VulnerabilitySeverity,
				"null",
				presidioResult.Type,
			)
		}
	}
}
