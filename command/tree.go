package command

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

func treeExtractor(data string) (string, error) {
	var hash string
	var IsTreeObject bool = false

	str, ok := strings.CutPrefix(data, "\x00")
	if !ok {
		return "", fmt.Errorf("fatal: not a tree object")
	}
	words := strings.Split(strings.Split(strings.TrimSpace(str), "\n")[0], " ")

	for i, word := range words {
		if strings.Compare(word, "tree") == 0 {
			hash = words[i+1]
			IsTreeObject = true
			break
		}
	}

	if !IsTreeObject {
		return "", fmt.Errorf("fatal: not a tree object")
	}
	return hash, nil
}

type TreeDataStruct struct {
	Mode     string
	FileType string
	Hex      string
	Filename string
}

func parseTreeObject(data []byte) []TreeDataStruct {

	headerEnd := bytes.IndexByte(data, 0)
	var treeContent []TreeDataStruct

	treeData := data[headerEnd+1:]

	for len(treeData) > 0 {
		modeEnd := bytes.IndexByte(treeData, ' ')
		mode := string(treeData[:modeEnd])

		// Move past the mode and read the filename (null-terminated)
		treeData = treeData[modeEnd+1:]
		nameEnd := bytes.IndexByte(treeData, 0)
		filename := string(treeData[:nameEnd])

		// Move past the filename and read the next 20 bytes (SHA-1 hash)
		treeData = treeData[nameEnd+1:]
		sha1Hash := treeData[:20]

		// Move past the hash for the next entry
		treeData = treeData[20:]

		var fileType string
		if mode == "100644" {
			fileType = "blob"
		} else {
			fileType = "tree"
		}

		treeContent = append(treeContent, TreeDataStruct{Mode: mode, Filename: filename, FileType: fileType, Hex: hex.EncodeToString(sha1Hash)})
	}

	return treeContent
}
