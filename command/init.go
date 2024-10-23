package command

import (
	"fmt"
	"os"
)

func InitFunc() (string, error) {
	for _, dir := range []string{".go-vcs", ".go-vcs/refs/heads", ".go-vcs/logs/refs/heads", ".go-vcs/objects"} {
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

	headFile, err := os.Create(".go-vcs/HEAD")
	if err != nil {
		return "", err
	}
	defer headFile.Close()

	headFile.WriteString("ref: refs/heads/main")

	return "Repository initialized successfully...", nil
}
