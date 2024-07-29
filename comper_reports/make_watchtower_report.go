package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
)

type OutputLog struct {
    Version     string    `json:"version"`
    PluginsUsed []Plugin  `json:"plugins_used"`
    FiltersUsed []Filter  `json:"filters_used"`
    Results     struct{}  `json:"results"`
}

type Plugin struct {
    Name  string  `json:"name"`
}

type Filter struct {
    Path     string `json:"path"`
    MinLevel int    `json:"min_level,omitempty"`
}

type Tool struct {
    Tool      string      `json:"tool"`
    OutputLog json.RawMessage `json:"output_log"`
}

type NotebookIssues map[string][]Tool


func readJSON(filename string) (*NotebookIssues, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var result NotebookIssues
    err = json.Unmarshal(data, &result)
    if err != nil {
        return nil, err
    }

    return &result, nil
}

func main() {
    jsonData, err := readJSON("severity_mapped_detailed_reports_1722173437.json")
    if err != nil {
        log.Fatalf("Error reading JSON file: %v", err)
    }

    file, err := os.Create("watchtower_report.txt")
    if err != nil {
        log.Fatalf("Error creating file: %v", err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)

    for notebookPath, tools := range *jsonData {
        writer.WriteString(fmt.Sprintf("Notebook Name: %s\n", notebookPath))
        for _, tool := range tools {
            writer.WriteString(fmt.Sprintf("\n  Tool: %s\n", tool.Tool))
            switch tool.Tool {
            case "Detect-Secret":
                var outputLog OutputLog
                err = json.Unmarshal(tool.OutputLog, &outputLog)
                if err != nil {
                    log.Printf("Error unmarshalling OutputLog: %v", err)
                    continue
                }
                writer.WriteString(fmt.Sprintf("  Version: %s\n", outputLog.Version))
                writer.WriteString("  Plugins Used:\n")
                for _, plugin := range outputLog.PluginsUsed {
                    writer.WriteString(fmt.Sprintf("      Name: %s\n", plugin.Name))
                }
                writer.WriteString("\n  Filters Used:\n")
                for _, filter := range outputLog.FiltersUsed {
                    writer.WriteString(fmt.Sprintf("      Path: %s\n", filter.Path))
                    if filter.MinLevel != 0 {
                        writer.WriteString(fmt.Sprintf("      Min Level: %d\n", filter.MinLevel))
                    }
                }
                writer.WriteString(fmt.Sprintf("  Results: %v\n", outputLog.Results))
            case "Whisper":
                var outputLog []interface{}
                err = json.Unmarshal(tool.OutputLog, &outputLog)
                if err != nil {
                    log.Printf("Error unmarshalling OutputLog: %v", err)
                    continue
                }
                writer.WriteString(fmt.Sprintf("  Output Log: %v\n", outputLog))
            case "Presidio-Analyzer", "Safety":
                var outputLog string
                err = json.Unmarshal(tool.OutputLog, &outputLog)
                if err != nil {
                    log.Printf("Error unmarshalling OutputLog: %v", err)
                    continue
                }
                writer.WriteString(fmt.Sprintf("  Output Log: \n\n%s\n", outputLog))
            default:
                writer.WriteString(fmt.Sprintf("  Output Log: \n\n%s\n", string(tool.OutputLog)))
            }
            writer.WriteString("\n")
        }
        writer.WriteString("\n")
    }

    writer.Flush()

    fmt.Println("Output written to output.txt")
}
