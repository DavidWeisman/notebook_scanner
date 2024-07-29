package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
)

type OutputLog struct {
    Version     string   `json:"version"`
    PluginsUsed []Plugin `json:"plugins_used"`
    FiltersUsed []Filter `json:"filters_used"`
    Results     struct{} `json:"results"`
}

type Plugin struct {
    Name string `json:"name"`
}

type Filter struct {
    Path     string `json:"path"`
    MinLevel int    `json:"min_level,omitempty"`
}

type Tool struct {
    Tool      string          `json:"tool"`
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

    dirPath := "/Users/david/desktop/notebook_scaner/scand_reports"

    // Walk through the directory to find subdirectories with JSON files
    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Check if the path is a directory and not the root directory
        if info.IsDir() && path != dirPath {
            // Read the directory contents
            files, err := ioutil.ReadDir(path)
            if err != nil {
                return err
            }

            for _, f := range files {
                if filepath.Ext(f.Name()) == ".json" {
                    jsonFilePath := filepath.Join(path, f.Name())

                    jsonData, err := readJSON(jsonFilePath)
                    if err != nil {
                        continue
                    }

                    // Create a new text file for each JSON file
                    outputFileName := f.Name()[:len(f.Name())-len(filepath.Ext(f.Name()))] + "_report.txt"
                    outputFilePath := filepath.Join(path, outputFileName)

                    file, err := os.Create(outputFilePath)
                    if err != nil {
                        log.Fatalf("Error creating file %s: %v", outputFilePath, err)
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
                                    continue
                                }
                                writer.WriteString(fmt.Sprintf("  Output Log: %v\n", outputLog))
                            case "Presidio-Analyzer", "Safety":
                                var outputLog string
                                err = json.Unmarshal(tool.OutputLog, &outputLog)
                                if err != nil {
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

                    fmt.Printf("Output written to %s\n", outputFilePath)
                }
            }
        }

        return nil
    })

    if err != nil {
        log.Fatalf("Error walking the path %s: %v", dirPath, err)
    }
}
