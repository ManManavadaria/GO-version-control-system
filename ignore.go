package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/ManManavadaria/GO-version-control-system/helper"
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

	var wg sync.WaitGroup
	var validFileChan = make(chan string)
	var inputChan = make(chan string)

	worker := runtime.NumCPU()

	for i := 1; i < worker; i++ {
		wg.Add(1)
		go func(outChan chan string, inChan chan string) {
			for file := range inChan {
				file = strings.TrimSpace(file)
				if len(file) == 0 || strings.HasPrefix(file, "#") {
					return
				}
				outChan <- file
			}
		}(validFileChan, inputChan)
	}

	for _, file := range fileNames {
		inputChan <- file
	}

	wg.Wait()
	close(inputChan)

	for file := range validFileChan {
		FilesToIgnore = append(FilesToIgnore, file)
	}

	FilesToIgnore = append(FilesToIgnore, ".git", ".go-vcs")

	ActiveFiles = GetAllFiles(".")
	return
}

func ValidateFileOptionArgument(files []string) error {
	if len(files) == 0 {
		helper.PrintError("Invalid command arguments, filename are missing.")
	}

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
