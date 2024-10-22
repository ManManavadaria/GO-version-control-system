package command

import (
	"fmt"
	"io/fs"
	"os"
)

func ListAllBranch() []fs.DirEntry {

	files, err := os.ReadDir(".go-vcs/refs/heads")
	if err != nil {
		fmt.Println("Error : ", err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println("Error : folder available in the heads directory", file.Name())
		}
	}

	return files
}
