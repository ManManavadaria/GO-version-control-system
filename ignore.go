package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	FilesToIgnore []string
	ActiveFiles   []string
)

// NOTE: Add error handling
func GetAllFiles(root string) []string {
	var files []string
	excludeMap := make(map[string]bool)
	for _, name := range FilesToIgnore {
		excludeMap[name] = true
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)

		if excludeMap[base] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// if path != root {
		if !info.IsDir() && path != root {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", root, err)
		return nil
	}

	return files
}

func init() {
	f, err := os.Open(".gitignore")
	if os.IsNotExist(err) {
		FilesToIgnore = append(FilesToIgnore, ".git", ".go-vcs")
		ActiveFiles = GetAllFiles(".")
		return
	}
	defer f.Close()

	b, _ := io.ReadAll(f)

	fileNames := strings.Split(strings.TrimSpace(string(b)), "\n")

	var validFileNames []string

	for _, filename := range fileNames {
		filename = strings.TrimSpace(filename)

		if len(filename) == 0 || strings.HasPrefix(filename, "#") {
			continue
		}

		validFileNames = append(validFileNames, filename)
	}

	FilesToIgnore = append(FilesToIgnore, validFileNames...)

	FilesToIgnore = append(FilesToIgnore, ".git", ".go-vcs")

	ActiveFiles = GetAllFiles(".")
	return
}

func ValidateFileOptionArgument(files []string) error {
	activeFiles := ActiveFiles

	var correctFiles map[string]bool = map[string]bool{}

	for _, filePath := range activeFiles {
		correctFiles[filePath] = true
	}

	for _, file := range files {
		if file == "." {
			return nil
		}
		_, ok := correctFiles[file]
		if !ok {
			return fmt.Errorf("Error searching file : %s", file)
		}
	}

	return nil
}
