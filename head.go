package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func FetchLatestCommitHash() string {

	f, err := os.Open(".go-vcs/logs/HEAD")
	if err != nil {
		fmt.Println("err ", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("err ", err)
	}

	lines := strings.Split(strings.TrimSpace(string(b)), "\n")

	words := strings.Split(lines[len(lines)-1], " ")

	return words[1]
}
