package command

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// func FetchLatestCommitHash() string {

// 	f, err := os.Open(".go-vcs/logs/HEAD")
// 	if err != nil {
// 		fmt.Println("err ", err)
// 	}
// 	defer f.Close()

// 	b, err := io.ReadAll(f)
// 	if err != nil {
// 		fmt.Println("err ", err)
// 	}

// 	lines := strings.Split(strings.TrimSpace(string(b)), "\n")

// 	words := strings.Split(lines[len(lines)-1], " ")

// 	return words[1]
// }

func FetchLatestCommitHash() string {

	headFile, err := os.Open(".go-vcs/HEAD")
	if err != nil {
		fmt.Println("err ", err)
	}
	defer headFile.Close()

	b, err := io.ReadAll(headFile)
	if err != nil {
		fmt.Println("err ", err)
	}

	ref := strings.Split(strings.TrimSpace(string(b)), " ")[1]

	refBranch, err := os.Open(fmt.Sprintf(".go-vcs/%s", ref))
	if err != nil {
		fmt.Println("err ", err)
	}
	defer refBranch.Close()

	hashByte, err := io.ReadAll(refBranch)
	if err != nil {
		fmt.Println("err ", err)
	}

	return string(strings.TrimSpace(string(hashByte)))
}
