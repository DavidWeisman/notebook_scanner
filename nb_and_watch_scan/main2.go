package main

import (
	"fmt"
	"os"
	"os/exec"
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

	repoName := getRepoName(githubURL)

	destDir := filepath.Join(desktopPath, repoName)
	err = os.Mkdir(destDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory on desktop: %v\n", err)
		return
	}
	
	tempDir := filepath.Join(os.TempDir(), repoName)
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: githubURL,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Printf("Failed to clone the repository: %v\n", err)
		os.RemoveAll(tempDir)
		return
	}

	defer os.RemoveAll(tempDir) // Clean up the temporary directory after the program ends

	srcPath := filepath.Join(tempDir, pathInRepo)
	destPath := filepath.Join(destDir, filepath.Base(pathInRepo))
	err = copy.Copy(srcPath, destPath)
	if err != nil {
		fmt.Printf("Failed to copy folder or file: %v\n", err)
		return
	}

	fmt.Println("Download complete! File copied to desktop directory.")

	err = runNbDefense(destPath)
	if err != nil {
		fmt.Printf("Failed to run nbdefense: %v\n", err)
		return
	}

	fmt.Println("nbdefense scan complete!")
	
	err = runWatchtower(destPath)
        if err != nil {
                fmt.Printf("Failed to run nbdefense: %v\n", err)
                return
        }

        fmt.Println("watchtower scan complete!")
	
	goFile := "/Users/david/desktop/notebook_scaner/nb_and_watch_scan/set_report_location.go"
	
	err = runGoFile(goFile)
	if err != nil {
		fmt.Printf("failed to run first Go file: %v", err)
		return
	}
}

func getRepoName(githubURL string) string {
	parts := strings.Split(githubURL, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}

func runNbDefense(dir string) error {
	cmd := exec.Command("nbdefense", "scan", "-o","json", dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runWatchtower(dir string) error {
        cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd /Users/david/desktop/watchtower/src && python watchtower.py --repo_type=folder --path=%s", dir))
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        return cmd.Run()
}

func runGoFile(goFile string) error {
	// Compile the Go file
	executable := goFile[:len(goFile)-3] // remove ".go" extension
	cmd := exec.Command("go", "build", "-o", executable, goFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to compile %s: %v", goFile, err)
	}

	// Run the compiled executable
	cmd = exec.Command(executable)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run %s: %v", executable, err)
	}

	// Clean up the compiled executable
	err = os.Remove(executable)
	if err != nil {
		return fmt.Errorf("failed to remove executable %s: %v", executable, err)
	}

	return nil
}
