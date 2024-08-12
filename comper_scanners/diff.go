package comper_scanners

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func diff() {

	file1 := "/Users/david/desktop/notebook_scaner/scanner/output1.txt"
	file2 := "/Users/david/desktop/notebook_scaner/scanner/output2.txt"

	outputFile1 := "unique_in_file1.txt"
	outputFile2 := "unique_in_file2.txt"

	linesFile1, err := readLines(file1)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", file1, err)
		return
	}

	linesFile2, err := readLines(file2)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", file2, err)
		return
	}

	setFile1 := extractVulnIDAndPkgNamePairs(linesFile1[1:])
	setFile2 := extractVulnIDAndPkgNamePairs(linesFile2[1:])

	uniqueInFile1 := filterLinesByVulnIDAndPkgName(linesFile1[1:], setFile2)
	uniqueInFile2 := filterLinesByVulnIDAndPkgName(linesFile2[1:], setFile1)

	writeLines(outputFile1, append([]string{linesFile1[0]}, uniqueInFile1...))
	writeLines(outputFile2, append([]string{linesFile2[0]}, uniqueInFile2...))

	fmt.Println("Comparison completed. Results written to:", outputFile1, "and", outputFile2)
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func extractVulnIDAndPkgNamePairs(lines []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, line := range lines {
		key := extractVulnIDAndPkgName(line)
		if key != "" {
			set[key] = struct{}{}
		}
	}
	return set
}

func extractVulnIDAndPkgName(line string) string {
	fields := strings.Fields(line)
	if len(fields) > 2 {
		return fields[1] + " " + fields[2]
	}
	return ""
}

func filterLinesByVulnIDAndPkgName(lines []string, idSet map[string]struct{}) []string {
	var filtered []string
	for _, line := range lines {
		key := extractVulnIDAndPkgName(line)
		if _, found := idSet[key]; !found {
			filtered = append(filtered, line)
		}
	}
	return filtered
}

func writeLines(filename string, lines []string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
}
