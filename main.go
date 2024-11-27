package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const scaffoldRepo = "https://github.com/pojol/braid-scaffold.git"
const defaultModuleName = "braid-scaffold"
const defaultVersion = "master"

func replaceInFile(filePath, oldStr, newStr string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	newContent := strings.ReplaceAll(string(content), oldStr, newStr)

	return os.WriteFile(filePath, []byte(newContent), 0644)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: braid-cli new <project-name> [version]")
		fmt.Println("Example: braid-cli new myserver v0.0.1")
		os.Exit(1)
	}

	projectName := os.Args[2]
	version := defaultVersion
	if len(os.Args) > 3 {
		version = os.Args[3]
	}

	// Clone the scaffold repository
	cmd := exec.Command("git", "clone", "-b", version, scaffoldRepo, projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to clone scaffold: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Successfully cloned scaffold template [version: %s] to %s\n", version, projectName)
	}

	// Change to the project directory
	os.Chdir(projectName)

	// Remove the .git directory
	os.RemoveAll(".git")

	// Remove the existing go.mod and go.sum files
	os.Remove("go.mod")
	os.Remove("go.sum")

	// Initialize a new Go module with the project name
	cmd = exec.Command("go", "mod", "init", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to initialize Go module: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("go", "get", fmt.Sprintf("github.com/pojol/braid@%s", version))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to initialize Go module: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("go", "get", "github.com/pojol/gobot@0.4.5")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to initialize gobot Go module: %v\n", err)
		os.Exit(1)
	}

	// Replace module name in all files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".go" {
			replaceInFile(path, defaultModuleName, projectName)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Failed to update module name in files: %v\n", err)
		os.Exit(1)
	}

	// Run go mod tidy to ensure all dependencies are properly managed
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to tidy Go module: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Braid project '%s' created and initialized successfully!\n", projectName)
}
