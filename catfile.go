package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func CatfileFunc() (string, error) {
	if len(os.Args) < 3 {
		return "", fmt.Errorf("invalid arguments, SHA missing.")
	}
	sha := os.Args[2]

	filepath := fmt.Sprintf(".go-vcs/objects/%v/%v", sha[0:2], sha[2:])

	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	r, err := zlib.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("error creating zlib reader: %v", err)
	}
	defer r.Close()

	s, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("error reading file content: %v", err)
	}

	return string(s), nil
}
