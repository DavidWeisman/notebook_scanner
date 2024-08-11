package comper_scanners

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	outputFile1Name = "output1.txt"
	outputFile2Name = "output2.txt"
	outputFile3Name = "output3.txt"
	outputFile4Name = "output4.txt"
	firsst_name     = "NBdefense"
	second_name     = "whatchtower"
	outputFile1     = "final_result1"
	outputFile2     = "final_result2"
	usageMessage    = "Usage: go run main.go <path_to_json_file1> <path_to_json_file2>"
)

func Comper() {

	jsonFileName1, err := Find_nb_file("/Users/david/desktop/notebook_scaner/scand_reports/nbdefense")
	if err != nil {
		log.Fatalf("Error finding newest file: %v", err)
	}

	jsonFileName2 := Find_watch_file("/Users/david/desktop/notebook_scaner/scand_reports/watchtower")

	trivyOutput, err := ReadAndUnmarshal[TrivyOutput](jsonFileName1)
	if err != nil {
		fmt.Printf("Error reading Trivy JSON: %v\n", err)
		return
	}

	data, err := ReadAndUnmarshal[JSONStructure](jsonFileName2)
	if err != nil {
		fmt.Printf("Error reading JSON structure: %v\n", err)
		return
	}

	notebookIssues, err := ReadAndUnmarshal[NotebookIssues](jsonFileName1)
	if err != nil {
		fmt.Printf("Error reading Notebook Issues JSON: %v\n", err)
		return
	}

	fileOutputs, err := ReadAndUnmarshal[FileOutputs](jsonFileName2)
	if err != nil {
		fmt.Printf("Error reading File Outputs JSON: %v\n", err)
		return
	}

	err = WriteOutputs(outputFile1Name, outputFile2Name, outputFile3Name, outputFile4Name, trivyOutput, data, notebookIssues, fileOutputs)
	if err != nil {
		fmt.Printf("Error writing outputs: %v\n", err)
		return
	}

	data1, err := ReadFile(outputFile1Name)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", outputFile1Name, err)
		return
	}

	data2, err := ReadFile(outputFile2Name)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", outputFile2Name, err)
		return
	}

	data3, err := ReadFile(outputFile3Name)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", outputFile1Name, err)
		return
	}

	data4, err := ReadFile(outputFile4Name)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", outputFile2Name, err)
		return
	}

	err = WriteTableToFile(outputFile1, firsst_name, data1, second_name, data2, data3, data4)
	if err != nil {
		fmt.Printf("Error writing table to file %s: %v\n", outputFile1, err)
	}

	//	err = WriteTableToFile(outputFile2, firsst_name, data3, second_name, data4)
	//	if err != nil {
	//		fmt.Printf("Error writing table to file %s: %v\n", outputFile2, err)
	//	}

	to_delete_files := []string{outputFile1Name, outputFile2Name, outputFile3Name, outputFile4Name}
	if err := DeleteFiles(to_delete_files); err != nil {
		fmt.Println(err)
	}
}

func DeleteFiles(files []string) error {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return fmt.Errorf("error deleting file %s: %w", file, err)
		}
	}
	return nil
}

func Find_nb_file(root string) (string, error) {
	var newestFile string
	var newestModTime time.Time

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().After(newestModTime) {
			newestModTime = info.ModTime()
			newestFile = path
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return newestFile, nil
}

func Find_watch_file(root string) string {
	newestDir, err := findNewestDir(root)
	if err != nil {
		log.Fatalf("Error finding newest directory: %v", err)
	}

	if newestDir == "" {
		log.Fatalf("No subdirectories found in %s", root)
	}

	filePrefix := "severity_mapped_detailed_reports"
	filePath, err := findFileWithPrefix(newestDir, filePrefix)
	if err != nil {
		log.Fatalf("Error finding file with prefix %q: %v", filePrefix, err)
	}

	return filePath
}

func findNewestDir(root string) (string, error) {
	var newestDir string
	var newestModTime time.Time

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() || path == root {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().After(newestModTime) {
			newestModTime = info.ModTime()
			newestDir = path
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return newestDir, nil
}

func findFileWithPrefix(dir, prefix string) (string, error) {
	var matchedFile string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if strings.HasPrefix(d.Name(), prefix) {
			matchedFile = path
			return filepath.SkipDir // Stop the walk once we found the file
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if matchedFile == "" {
		return "", fmt.Errorf("no file with prefix %q found in directory %q", prefix, dir)
	}

	return matchedFile, nil
}
