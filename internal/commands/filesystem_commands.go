package commands

import (
	"fmt"
	"os"
	"strings"
)

// === Make Directory
func HandleMakeDirectory(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("no dir path specified")
	}

	dirName := args[0]

	err := os.Mkdir(dirName, 0755)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("made directory %s", dirName), nil
}

// === Cat
func HandleCat(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("no file specified")
	}
	return readFile(args[1])
}

func readFile(filename string) (string, error) {
	// Read the file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Return the content as a string
	return string(content), nil
}

// === CD
func HandleChangeDir(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("no directory specified")
	}
	return changeDir(args[1])
}

func changeDir(path string) (string, error) {
	// Use os.Chdir to change directories
	err := os.Chdir(path)
	if err != nil {
		return "", fmt.Errorf("failed to change directory: %v", err)
	}
	return "Directory changed successfully\n", nil
}

// === ls
func HandleLS(args ...string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Read the directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return strings.Join(fileNames, "\n"), nil
}

// === PWD
func HandlePWD(args ...string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir, nil
}
