package command

import (
	"fmt"
	"os"
)

func InitFunc() (string, error) {
	for _, dir := range []string{".go-vcs", ".git/refs", ".go-vcs/objects"} {
		info, err := os.Stat(dir)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", fmt.Errorf("Error creating directory: %s", err)
			}
			continue
		}
		if info.IsDir() {
			return "", fmt.Errorf("Repository already exists.")
		}
	}
	return "Repository initialized successfully...", nil
}
