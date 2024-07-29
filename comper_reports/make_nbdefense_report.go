package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "bufio"
)

// Define structs based on the JSON structure
type Root struct {
    Root           string          `json:"root"`
    RootIssues     []interface{}   `json:"root_issues"`
    Plugins        []Plugin        `json:"plugins"`
    Notebooks      []string        `json:"notebooks"`
    NotebookIssues []NotebookIssue `json:"notebook_issues"`
}

type Plugin struct {
    Name     string   `json:"name"`
    Settings Settings `json:"settings"`
}

type Settings struct {
    Enabled             bool                `json:"enabled"`
    ConfidenceThreshold float64             `json:"confidence_threshold,omitempty"`
    Entities            map[string]bool     `json:"entities,omitempty"`
    SecretsPlugins      []SecretsPlugin     `json:"secrets_plugins,omitempty"`
}

type SecretsPlugin struct {
    Name           string  `json:"name"`
    KeywordExclude string  `json:"keyword_exclude,omitempty"`
    Limit          float64 `json:"limit,omitempty"`
}

type NotebookIssue struct {
    Path   string   `json:"path"`
    Issues []Issue  `json:"issues"`
}

type Issue struct {
    Code        string   `json:"code"`
    Severity    string   `json:"severity"`
    Cell        Cell     `json:"cell"`
    Location    string   `json:"location"`
    Details     Details  `json:"details"`
    Issues      []Issue  `json:"issues"`
    LineIndex   int      `json:"line_index,omitempty"`
    StartIndex  int      `json:"character_start_index,omitempty"`
    EndIndex    int      `json:"character_end_index,omitempty"`
}

type Cell struct {
    CellIndex       int    `json:"cell_index"`
    CellType        string `json:"cell_type"`
    ScrubbedContent string `json:"scrubbed_content"`
}

type Details struct {
    Description   string                `json:"description"`
    SummaryField  map[string]int        `json:"summary_field"`
}

// Function to read JSON file and unmarshal into the struct
func readJSON(filename string) (*Root, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var result Root
    err = json.Unmarshal(data, &result)
    if err != nil {
        return nil, err
    }

    return &result, nil
}

func main() {
    // Read the JSON file
    jsonData, err := readJSON("nbdefense0728-1530.json")
    if err != nil {
        log.Fatalf("Error reading JSON file: %v", err)
    }

    // Create a file to write the output
    file, err := os.Create("nbdefense_report.txt")
    if err != nil {
        log.Fatalf("Error creating file: %v", err)
    }
    defer file.Close()

    // Create a buffered writer
    writer := bufio.NewWriter(file)

    writer.WriteString("\nPlugins:\n")
    for _, plugin := range jsonData.Plugins {
        writer.WriteString(fmt.Sprintf("\n  Name: %s\n", plugin.Name))
        writer.WriteString(fmt.Sprintf("  Enabled: %t\n", plugin.Settings.Enabled))
        if plugin.Settings.ConfidenceThreshold != 0 {
            writer.WriteString(fmt.Sprintf("  ConfidenceThreshold: %f\n", plugin.Settings.ConfidenceThreshold))
        }
        if plugin.Settings.Entities != nil {
            writer.WriteString("  Entities:\n")
            for entity, enabled := range plugin.Settings.Entities {
                writer.WriteString(fmt.Sprintf("    %s: %t\n", entity, enabled))
            }
        }
        if plugin.Settings.SecretsPlugins != nil {
            writer.WriteString("  SecretsPlugins:\n")
            for _, secretPlugin := range plugin.Settings.SecretsPlugins {
                writer.WriteString(fmt.Sprintf("    Name: %s\n", secretPlugin.Name))
                if secretPlugin.KeywordExclude != "" {
                    writer.WriteString(fmt.Sprintf("    KeywordExclude: %s\n", secretPlugin.KeywordExclude))
                }
            }
        }
    }

    writer.WriteString("\nNotebooks:\n")
    for _, notebook := range jsonData.Notebooks {
        writer.WriteString(fmt.Sprintf("    Name: %s\n", notebook))
    }    

    writer.WriteString("\nNotebook Issues:\n")
    for _, notebookIssue := range jsonData.NotebookIssues {
        writer.WriteString(fmt.Sprintf("Notebook Name: %s\n", notebookIssue.Path))
        for _, issue := range notebookIssue.Issues {
            printIssue(writer, issue, "  ")
        }
        writer.WriteString("\n\n\n")
    }

    // Flush the buffered writer to ensure all data is written to the file
    writer.Flush()
}

func printIssue(writer *bufio.Writer, issue Issue, indent string) {
    writer.WriteString("———————————————————————————————————————————————————————————————————————————————————\n")
    writer.WriteString(fmt.Sprintf("Code: %s\n", issue.Code))
    writer.WriteString(fmt.Sprintf("Severity: %s\n", issue.Severity))
    writer.WriteString("Cell:\n")
    writer.WriteString(fmt.Sprintf("%s  CellIndex: %d\n", indent, issue.Cell.CellIndex))
    writer.WriteString(fmt.Sprintf("%s  CellType: %s\n", indent, issue.Cell.CellType))
    writer.WriteString("\n\nScrubbedContent:\n")
    writer.WriteString(fmt.Sprintf("%s\n\n", issue.Cell.ScrubbedContent))
    writer.WriteString(fmt.Sprintf("\nLocation: %s\n", issue.Location))
    if issue.Details.Description != "" {
        writer.WriteString("Details:\n")
        writer.WriteString(fmt.Sprintf("%s  Description: %s\n", indent, issue.Details.Description))
        if len(issue.Details.SummaryField) > 0 {
            writer.WriteString(fmt.Sprintf("%s  SummaryField:\n", indent))
            for field, count := range issue.Details.SummaryField {
                writer.WriteString(fmt.Sprintf("%s    %s: %d\n", indent, field, count))
            }
        }
    }
    if len(issue.Issues) > 0 {
        writer.WriteString("\n\n\nNested Issues:\n")
        for _, nestedIssue := range issue.Issues {
            printIssue(writer, nestedIssue, indent+"  ")
        }
    }
    if issue.LineIndex != 0 {
        writer.WriteString(fmt.Sprintf("%sLineIndex: %d\n", indent, issue.LineIndex))
    }
    if issue.StartIndex != 0 || issue.EndIndex != 0 {
        writer.WriteString(fmt.Sprintf("%sCharacter Range: %d-%d\n\n\n", indent, issue.StartIndex, issue.EndIndex))
    }
}
