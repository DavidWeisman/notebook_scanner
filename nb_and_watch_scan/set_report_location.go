package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Example usage
	srcDir1 := "/Users/david/desktop/notebook_scaner/nb_and_watch_scan"
	destDir1 := "/Users/david/desktop/notebook_scaner/scand_reports/nbdefense"

	srcDir2 := "/Users/david/Desktop/watchtower/src/scanned_reports"
        destDir2 := "/Users/david/desktop/notebook_scaner/scand_reports/watchtower"

	err1 := moveHTMLFiles(srcDir1, destDir1)
	if err1 != nil {
		fmt.Printf("Failed to move HTML files: %v\n", err1)
	} else {
		fmt.Println("Successfully moved HTML files.")
	}

	err2 := moveDirectories(srcDir2, destDir2)
        if err2 != nil {
                fmt.Printf("Failed to move directories: %v\n", err2)
        } else {
                fmt.Println("Successfully moved directories.")
        }
}

func moveDirectories(srcDir, destDir string) error {
        // Ensure the destination directory exists
        err := os.MkdirAll(destDir, 0755)
        if err != nil {
                return fmt.Errorf("failed to create destination directory: %v", err)
        }

        // Read the source directory
        entries, err := os.ReadDir(srcDir)
        if err != nil {
                return fmt.Errorf("failed to read source directory: %v", err)
        }

        for _, entry := range entries {
                if entry.IsDir() {
                        srcPath := filepath.Join(srcDir, entry.Name())
                        destPath := filepath.Join(destDir, entry.Name())

                        // Move the directory to the destination directory
                        err := os.Rename(srcPath, destPath)
                        if err != nil {
                                return fmt.Errorf("failed to move directory %s: %v", srcPath, err)
                        }
                }
        }

        return nil
}


// moveHTMLFiles moves HTML files from srcDir to destDir
func moveHTMLFiles(srcDir, destDir string) error {
	// Ensure the destination directory exists
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Walk the source directory to find HTML files
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has .html extension
		if !info.IsDir() && filepath.Ext(info.Name()) == ".html" {
			// Define the destination path for the HTML file
			destPath := filepath.Join(destDir, info.Name())

			// Move the file to the destination directory
			err := moveFile(path, destPath)
			if err != nil {
				return fmt.Errorf("failed to move file %s: %v", path, err)
			}

			// Remove the original file
			err = os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to remove original file %s: %v", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk the source directory: %v", err)
	}

	return nil
}

// moveFile copies a file from src to dst
func moveFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}
