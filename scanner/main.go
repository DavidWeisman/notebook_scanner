package main

import (
	"fmt"
	"io"
	"log"

	"notebook_scaner/comper_scanners"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
)

// test
const (
	usage        = "Usage: go run main.go <github_url> <path_in_repo>"
	nbDefenseCmd = "nbdefense"
)

func main() {

	begginingPath, err := comper_scanners.GetParentDirBeforeNotebookScaner()
	if err != nil {
		fmt.Println(err)
		return
	}

	scannerSrcDir := filepath.Join(begginingPath, "notebook_scanner/scanner")

	scannerDstDir := filepath.Join(begginingPath, "notebook_scanner/scand_reports/nbdefense")

	watchtowerSrcDir := filepath.Join(begginingPath, "notebook_scanner/watchtower/src/scanned_reports")
	watchtowerDstDir := filepath.Join(begginingPath, "notebook_scanner/scand_reports/watchtower")

	githubURL := os.Args[1]
	pathInRepo := os.Args[2]

	desktopPath, err := getDesktopPath()
	if err != nil {
		log.Printf("Failed to get desktop path: %v", err)
	}

	repoName := getRepoName(githubURL)
	destDir := filepath.Join(desktopPath, repoName)
	tempDir := filepath.Join(os.TempDir(), repoName)

	if err := os.Mkdir(destDir, 0755); err != nil {
		log.Printf("Failed to create directory on desktop: %v", err)
	}

	if err := cloneRepo(githubURL, tempDir); err != nil {
		log.Printf("Failed to clone the repository: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcPath := filepath.Join(tempDir, pathInRepo)
	destPath := filepath.Join(destDir, filepath.Base(pathInRepo))
	if err := copy.Copy(srcPath, destPath); err != nil {
		log.Printf("Failed to copy folder or file: %v", err)
	}

	fmt.Println("Download complete! File copied to desktop directory.")

	if err := runNbDefense(destPath); err != nil {
		log.Printf("Failed to run nbdefense: %v", err)
	}
	fmt.Println("nbdefense scan complete!")

	if err := runWatchtower(destPath); err != nil {
		log.Printf("Failed to run watchtower: %v", err)
	}
	fmt.Println("watchtower scan complete!")

	if err := moveHTMLFiles(scannerSrcDir, scannerDstDir); err != nil {
		log.Printf("Failed to move HTML files: %v", err)
	}
	fmt.Println("Successfully moved nbdefense report.")

	if err := moveDirectories(watchtowerSrcDir, watchtowerDstDir); err != nil {
		log.Printf("Failed to move directories: %v", err)
	}
	fmt.Println("Successfully moved watchtower report.")

	comper_scanners.Comper()

}

func getDesktopPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %v", err)
	}
	desktopPath := filepath.Join(usr.HomeDir, "Desktop")
	if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
		return "", fmt.Errorf("desktop path does not exist")
	}
	return desktopPath, nil
}

func getRepoName(githubURL string) string {
	parts := strings.Split(githubURL, "/")
	return strings.TrimSuffix(parts[len(parts)-1], ".git")
}

func cloneRepo(githubURL, tempDir string) error {
	_, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      githubURL,
		Progress: os.Stdout,
	})
	return err
}

func runNbDefense(dir string) error {
	return runCommand(nbDefenseCmd, "scan", "-o", "json", dir)
}

func runWatchtower(dir string) error {

	begginingPath, err := comper_scanners.GetParentDirBeforeNotebookScaner()
	if err != nil {
		fmt.Println(err)
	}

	watchtowerSkrip := filepath.Join(begginingPath, "notebook_scanner/watchtower/src/watchtower.py")
	watchtowerSkrip2 := filepath.Join(begginingPath, "notebook_scanner/watchtower/src")
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("cd %s && python %s --repo_type=folder --path=%s", watchtowerSkrip2, watchtowerSkrip, dir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func moveHTMLFiles(srcDir, destDir string) error {
	return moveFiles(srcDir, destDir, ".json")
}

func moveFiles(srcDir, destDir, ext string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ext {
			destPath := filepath.Join(destDir, info.Name())
			if err := moveFile(path, destPath); err != nil {
				return fmt.Errorf("failed to move file %s: %v", path, err)
			}
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove original file %s: %v", path, err)
			}
		}
		return nil
	})
}

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
	return err
}

func moveDirectories(srcDir, destDir string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %v", err)
	}

	for _, entry := range entries {

		if entry.IsDir() {
			srcPath := filepath.Join(srcDir, entry.Name())
			destPath := filepath.Join(destDir, entry.Name())

			if err := os.Rename(srcPath, destPath); err != nil {
				return fmt.Errorf("failed to move directory %s: %v", srcPath, err)
			}
		}
	}

	return nil
}
