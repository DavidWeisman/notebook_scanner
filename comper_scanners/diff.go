package comper_scanners

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func diff() {

	begginingPath, err := GetParentDirBeforeNotebookScaner()
	if err != nil {
		fmt.Println(err)
		return
	}

	file1 := filepath.Join(begginingPath, "notebook_scanner/scanner/output1.txt")
	file2 := filepath.Join(begginingPath, "notebook_scanner/scanner/output2.txt")
	file3 := filepath.Join(begginingPath, "notebook_scanner/scanner/output3.txt")
	file4 := filepath.Join(begginingPath, "notebook_scanner/scanner/output4.txt")

	// Read the contents of the input files
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

	// Read additional files
	firstLines, firstHeader, err := readLinesWithHeader(file3)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", file3, err)
		return
	}

	secondLines, secondHeader, err := readLinesWithHeader(file4)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", file4, err)
		return
	}

	// Create sets for comparisons
	firstFileSet := createFileSet(firstLines)
	secondFileSet := createFileSet(secondLines)

	// Find unique lines in each file
	uniqueInFirst := filterUniqueLines(firstLines, secondFileSet)
	uniqueInSecond := filterUniqueLines(secondLines, firstFileSet)

	// Open the combined output file
	outputFile := "unique results.txt"
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", outputFile, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write results to the combined file
	writeSection(writer, "Unique in  NBdefense \n -------------------------------------------------------------------------", append([]string{linesFile1[0]}, uniqueInFile1...))
	writeSection(writer, "Unique in  Watchtower \n -------------------------------------------------------------------------", append([]string{linesFile2[0]}, uniqueInFile2...))
	writeSection(writer, "Unique in NBdefense \n -------------------------------------------------------------------------", append([]string{firstHeader}, uniqueInFirst...))
	writeSection(writer, "Unique in Watchtower \n -------------------------------------------------------------------------", append([]string{secondHeader}, uniqueInSecond...))

	writer.Flush()

	fmt.Println("Processing completed. Results written to:", outputFile)
}

func writeSection(writer *bufio.Writer, title string, lines []string) {
	fmt.Fprintf(writer, "\n\n%s\n", title)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
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

func readLinesWithHeader(filename string) ([]string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		header := scanner.Text()
		lines = append(lines, header)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return lines, header, scanner.Err()
	}
	return nil, "", scanner.Err()
}

func createFileSet(lines []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, line := range lines[1:] {
		key1 := extractKeyBySummaryField(line)
		key2 := extractKeyByDescription(line)
		set[key1] = struct{}{}
		set[key2] = struct{}{}
	}
	return set
}

func extractKeyBySummaryField(line string) string {
	fields := strings.Fields(line)
	if len(fields) >= 5 {
		fileName := fields[0]
		summaryField := fields[len(fields)-1]
		return fileName + " " + summaryField
	}
	return ""
}

func extractKeyByDescription(line string) string {
	fields := strings.Fields(line)
	if len(fields) >= 4 {
		fileName := fields[0]
		description := fields[3]
		return fileName + " " + description
	}
	return ""
}

func filterUniqueLines(lines []string, comparisonSet map[string]struct{}) []string {
	var filtered []string
	for _, line := range lines[1:] {
		key1 := extractKeyBySummaryField(line)
		key2 := extractKeyByDescription(line)
		if _, found1 := comparisonSet[key1]; !found1 {
			if _, found2 := comparisonSet[key2]; !found2 {
				filtered = append(filtered, line)
			}
		}
	}
	return filtered
}
