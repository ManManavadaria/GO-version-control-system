package command

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CatfileFunc(sha string) (string, error) {

	filePath := fmt.Sprintf(".git/objects/%v/%v/", sha[0:2], sha[2:])

	file, err := os.Open(filepath.Dir(filePath))
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

	str := string(s)

	return str[strings.Index(str, "\x00"):], nil
}

func RemoveFile(sha string) error {

	filePath := fmt.Sprintf(".git/objects/%v/%v/", sha[0:2], sha[2:])

	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
