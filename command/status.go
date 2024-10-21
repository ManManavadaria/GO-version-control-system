package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/ManManavadaria/GO-version-control-system/helper"
)

func StatusFunc(ActiveFiles []string) []FileStatusStruct {
	if len(ActiveFiles) == 0 {
		return nil
	}

	idx := LoadIndex()

	return idx.IndexHashCompare(ActiveFiles)

	// return idx.IndexHashComape(FetchLatestCommitHash(), &ActiveFiles, "")
}

type FileStatusStruct struct {
	Filename string
	Status   string
	BlobHash string
}

func (idx *Index) IndexHashCompare(activeFiles []string) []FileStatusStruct {
	var fileTrack []FileStatusStruct
	indexFileMap := make(map[string]IndexEntry)
	for _, entry := range idx.Entries {
		indexFileMap[entry.Path] = entry
	}

	for _, path := range activeFiles {
		// if info.IsDir() {
		// 	continue
		// }

		// relativePath := strings.TrimPrefix(path, workingDir+"/")
		entry, existsInIndex := indexFileMap[path]

		//NOTE: implement a create hash only function to avaoid using RemoveFile funtion
		currentHash, err := GenerateBlobHash(path)
		if err != nil {
			fmt.Println("err HashObjectFunc", err)
		}

		if existsInIndex {
			if fmt.Sprintf("%x", entry.Sha1) != currentHash {
				fileTrack = append(fileTrack, FileStatusStruct{Filename: path, Status: "modified", BlobHash: fmt.Sprintf("%x", entry.Sha1)})
				// fmt.Printf("Modified: %s\n", path)
			}
			delete(indexFileMap, path)
		} else {
			fileTrack = append(fileTrack, FileStatusStruct{Filename: path, Status: "new file", BlobHash: ""})
			// fmt.Printf("New file: %s\n", path)
		}
	}

	if len(indexFileMap) > 0 {
		for path := range indexFileMap {
			fileTrack = append(fileTrack, FileStatusStruct{Filename: path, Status: "removed", BlobHash: ""})
		}
	}

	return fileTrack
}

func StagedFiles() []TreeDataStruct {
	var stagedFiles []TreeDataStruct
	hash := FetchLatestCommitHash()

	treeData, err := LsTreeFuncAllFilesSearch(hash)
	if err != nil {
		helper.PrintError(err.Error())
	}

	var treeMap map[string]TreeDataStruct = map[string]TreeDataStruct{}

	for _, file := range treeData {
		treeMap[file.Filename] = file
	}

	index := LoadIndex()

	for _, entry := range index.Entries {
		file, ok := treeMap[entry.Path]
		if ok {
			if file.Hex != fmt.Sprintf("%x", entry.Sha1) {
				stagedFiles = append(stagedFiles, file)
			}
		}
	}

	return stagedFiles
}

// func StagedFiles(treehash string, fileTrack *[]string, path string) []FileStatusStruct {
// 	var changedFilesData []FileStatusStruct

// 	treeData, err := LsTreeFunc(treehash, []string{})
// 	if err != nil {
// 		fmt.Println("error ", err)
// 	}

// 	for i := len(treeData) - 1; i >= 0; i-- {
// 		file := treeData[i]
// 		if file.FileType == "blob" {
// 			if ok := IsFileAvailable(file.Filename); !ok {
// 				changedFilesData = append(changedFilesData, FileStatusStruct{
// 					Filename: file.Filename,
// 					Status:   "Removed file",
// 					BlobHash: "",
// 				})
// 				treeData = append(treeData[:i], treeData[i+1:]...)
// 			}
// 		}
// 	}

// 	for _, file := range treeData {
// 		if file.FileType == "blob" {
// 			if path != "" {
// 				file.Filename = path + "/" + file.Filename
// 			}
// 			hash, err := HashObjectFunc(file.Filename)
// 			if err != nil {
// 				fmt.Println("err ", err)
// 			}
// 			if hash != file.Hex {
// 				changedFilesData = append(changedFilesData, FileStatusStruct{
// 					Filename: file.Filename,
// 					Status:   "modified",
// 					BlobHash: hash,
// 				})
// 				go RemoveFile(hash)
// 				removeElement(fileTrack, file.Filename)
// 			} else if hash == file.Hex {
// 				removeElement(fileTrack, file.Filename)
// 			}
// 		} else if file.FileType == "tree" {
// 			changedFilesData = append(changedFilesData, StaggedFiles(file.Hex, fileTrack, file.Filename)...)
// 		}
// 	}

// 	if len(*fileTrack) > 0 {
// 		for _, file := range *fileTrack {
// 			hash, err := HashObjectFunc(file)
// 			if err != nil {
// 				fmt.Println("error : ", err)
// 			}
// 			go RemoveFile(hash)
// 			changedFilesData = append(changedFilesData, FileStatusStruct{Filename: file, Status: "new file", BlobHash: hash})
// 		}
// 	}
// 	//NOTE: Implement Heuristic Analysis on deleted and new added files by comparing its content to determine the renamed files
// 	return changedFilesData
// }

func removeElement(slice *[]string, element string) {
	for i, v := range *slice {
		if v == element {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}

func IsFileAvailable(file string) bool {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}
