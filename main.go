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
		fmt.Println("Usage: braid-cli new <project-name>")
		os.Exit(1)
	}

	projectName := os.Args[2]

	// Clone the scaffold repository
	cmd := exec.Command("git", "clone", scaffoldRepo, projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to clone scaffold: %v\n", err)
		os.Exit(1)
	}

	// Change to the project directory
	os.Chdir(projectName)

	// Remove the .git directory
	os.RemoveAll(".git")

	// Remove the existing go.mod and go.sum files
	os.Remove("go.mod")
	os.Remove("go.sum")

	// Initialize a new git repository
	cmd = exec.Command("git", "init")
	cmd.Run()

	// Initialize a new Go module with the project name
	cmd = exec.Command("go", "mod", "init", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to initialize Go module: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("go", "get", "github.com/pojol/braid@master")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to initialize Go module: %v\n", err)
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
