package command

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
)

func HashObjectFunc(fileName string) (string, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("err :", err)
		return "", fmt.Errorf("Failed to read file: %v", err)
	}

	normalizedContent := bytes.ReplaceAll(fileContent, []byte("\r\n"), []byte("\n"))

	blobHeader := fmt.Sprintf("blob %d\x00", len(normalizedContent))

	blobData := append([]byte(blobHeader), normalizedContent...)

	hash := sha1.New()
	hash.Write(blobData)
	hex := hash.Sum(nil)

	hashString := fmt.Sprintf("%x", hex)

	dirName := filepath.Join(".go-vcs", "objects", hashString[:2])
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return "", fmt.Errorf("Failed to create directory: %v", err)
	}

	filePath := filepath.Join(dirName, hashString[2:])
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("Failed to create blob file: %v", err)
	}
	defer f.Close()

	var buffer bytes.Buffer
	zlibWriter := zlib.NewWriter(&buffer)
	_, err = zlibWriter.Write(blobData)
	if err != nil {
		return "", fmt.Errorf("Failed to write compressed data: %v", err)
	}
	zlibWriter.Close()

	_, err = f.Write(buffer.Bytes())
	if err != nil {
		return "", fmt.Errorf("Failed to write to blob file: %v", err)
	}

	return hashString, nil
}
func GenerateBlobHash(fileName string) (string, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("err :", err)
		return "", fmt.Errorf("Failed to read file: %v", err)
	}

	normalizedContent := bytes.ReplaceAll(fileContent, []byte("\r\n"), []byte("\n"))

	blobHeader := fmt.Sprintf("blob %d\x00", len(normalizedContent))

	blobData := append([]byte(blobHeader), normalizedContent...)

	hash := sha1.New()
	hash.Write(blobData)
	hex := hash.Sum(nil)

	hashString := fmt.Sprintf("%x", hex)

	return hashString, nil
}
