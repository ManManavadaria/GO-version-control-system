package command

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type treeFileStruct struct {
	entries []IndexEntry
	fileMap map[string][]TreeDataStruct
	mu      sync.Mutex
}

func WriteTree() string {
	index := LoadIndex()

	tfs := treeFileStruct{
		entries: index.Entries,
		fileMap: map[string][]TreeDataStruct{},
	}

	hash, _ := tfs.TraversIndexEntry(".")

	return hash
}

func (tfs *treeFileStruct) TraversIndexEntry(root string) (string, string) {

	for i := len(tfs.entries) - 1; i >= 0; i-- {
		dirRoot, ok := CheckFileLocation(root, tfs.entries[i].Path)
		if ok && dirRoot == root {
			tfs.mu.Lock()
			var fileName string
			if root != "." {
				fileName = strings.TrimPrefix(tfs.entries[i].Path, root+"\\")
			} else {
				fileName = tfs.entries[i].Path
			}
			tfs.fileMap[root] = append(tfs.fileMap[root],
				TreeDataStruct{
					// Mode:     strconv.Itoa(int(tfs.entries[i].Mode)),
					Mode:     "100644",
					FileType: "blob",
					Hex:      fmt.Sprintf("%x", tfs.entries[i].Sha1),
					Filename: fileName,
				})
			tfs.mu.Unlock()
		} else if ok && dirRoot != root {
			if _, ok := tfs.fileMap[dirRoot]; !ok {
				hash, treeName := tfs.TraversIndexEntry(dirRoot)
				tfs.mu.Lock()
				tfs.fileMap[root] = append(tfs.fileMap[root],
					TreeDataStruct{
						Mode:     "040000",
						FileType: "tree",
						Hex:      hash,
						Filename: treeName,
					})
				tfs.mu.Unlock()
			}
			continue
		} else {
			continue
		}
	}

	return GenerateTreeHash(tfs.fileMap[root]), root
}

func CheckFileLocation(rootDir string, filePath string) (string, bool) {
	relPath, err := filepath.Rel(rootDir, filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return "", false
	}

	if !strings.Contains(relPath, string(filepath.Separator)) {
		return rootDir, true
	}

	dir := filepath.Dir(relPath)

	if strings.HasPrefix(relPath, "..") {
		return "", false
	}
	return dir, true
}

type ByFileName []TreeDataStruct

func (a ByFileName) Len() int           { return len(a) }
func (a ByFileName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFileName) Less(i, j int) bool { return a[i].Filename < a[j].Filename }

func GenerateTreeHash(treeEntries []TreeDataStruct) string {
	sort.Sort(ByFileName(treeEntries))

	treeContent, err := createTreeObject(treeEntries)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	header := fmt.Sprintf("tree %d\000", len(treeContent))
	fullContent := append([]byte(header), treeContent...)

	hash := sha1.Sum(fullContent)

	hashString := fmt.Sprintf("%x", hash)

	dirName := filepath.Join(".go-vcs", "objects", hashString[:2])
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		fmt.Println("Error: ", err)
	}

	filePath := filepath.Join(dirName, hashString[2:])
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("error :", err)
	}
	defer f.Close()

	var buffer bytes.Buffer
	zlibWriter := zlib.NewWriter(&buffer)
	_, err = zlibWriter.Write(fullContent)
	if err != nil {
		fmt.Println("error :", err)
	}
	zlibWriter.Close()

	_, err = f.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("error :", err)
	}

	return hex.EncodeToString(hash[:])
}

func createTreeObject(entries []TreeDataStruct) ([]byte, error) {
	var buffer bytes.Buffer

	for _, entry := range entries {
		buffer.WriteString(entry.Mode + " ")

		buffer.WriteString(entry.Filename)
		buffer.WriteByte(0)

		hashBytes, err := hex.DecodeString(entry.Hex)
		if err != nil {
			return nil, err
		}

		buffer.Write(hashBytes)
	}

	return buffer.Bytes(), nil
}
