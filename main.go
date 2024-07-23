package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <github_url> <path_in_repo>")
		return
	}

	githubURL := os.Args[1]
	pathInRepo := os.Args[2]

	// Get the current user to determine the desktop path
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to get current user: %v\n", err)
		return
	}

	desktopPath := filepath.Join(usr.HomeDir, "Desktop")
	if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
		fmt.Println("Desktop path does not exist")
		return
	}

	// Parse the repository URL to extract the repo name
	repoName := getRepoName(githubURL)

	// Clone the repository to a temporary directory
	tempDir := filepath.Join(os.TempDir(), repoName)
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: githubURL,
	})
	if err != nil {
		fmt.Printf("Failed to clone the repository: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir) // Clean up the temporary directory after the program ends

	// Copy the specified folder or file to the desktop path
	srcPath := filepath.Join(tempDir, pathInRepo)
	destPath := filepath.Join(desktopPath, filepath.Base(pathInRepo))
	err = copy.Copy(srcPath, destPath)
	if err != nil {
		fmt.Printf("Failed to copy folder or file: %v\n", err)
		return
	}

	fmt.Println("Download complete! File copied to desktop.")
}

// getRepoName extracts the repository name from a GitHub URL
func getRepoName(githubURL string) string {
	parts := strings.Split(githubURL, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}
