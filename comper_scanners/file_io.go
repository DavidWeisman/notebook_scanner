package comper_scanners

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func ReadFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func WriteTableToFile(filePath, name1 string, data1 []string, name2 string, data2 []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Calculate the maximum width of the columns
	maxWidth := 0
	for _, line := range data1 {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	for _, line := range data2 {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	maxWidth += 4 // Add some padding

	// Write headers
	fmt.Fprintf(writer, "%-*s | %-*s\n", maxWidth, name1, maxWidth, name2)
	writer.WriteString(strings.Repeat("-", maxWidth*2+3) + "\n")

	// Write rows
	maxLines := len(data1)
	if len(data2) > maxLines {
		maxLines = len(data2)
	}

	for i := 0; i < maxLines; i++ {
		var line1, line2 string
		if i < len(data1) {
			line1 = data1[i]
		}
		if i < len(data2) {
			line2 = data2[i]
		}
		fmt.Fprintf(writer, "%-*s | %-*s\n", maxWidth, line1, maxWidth, line2)
	}

	return writer.Flush()
}

func WriteOutput(filename string, writeFunc func(*tabwriter.Writer)) error {
	outputFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := tabwriter.NewWriter(outputFile, 0, 8, 2, ' ', 0)
	defer writer.Flush()

	writeFunc(writer)
	return nil
}
