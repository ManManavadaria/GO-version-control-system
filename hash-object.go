package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
)

func hashObjectFunc(fileName string) (string, error) {

	file, _ := os.ReadFile(fileName)

	blobHeader := fmt.Sprintf("blob %d\x00", len(file))

	blobData := append([]byte(blobHeader), file...)

	hex := sha1.Sum(blobData)

	var buffer bytes.Buffer
	z := zlib.NewWriter(&buffer)
	z.Write([]byte(blobData))
	z.Close()

	dirName := filepath.Dir(fmt.Sprintf(".go-vcs/objects/%x/", hex[0:1]))

	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return "", fmt.Errorf("Failed to create dir")
	}
	f, err := os.Create(fmt.Sprintf(".go-vcs/objects/%x/%x", hex[0:1], hex[1:]))
	if err != nil {
		return "", fmt.Errorf("Failed to create blob file")
	}
	defer f.Close()

	f.Write(buffer.Bytes())

	return fmt.Sprintf("%x", hex), nil
}
