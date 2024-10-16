package main

import (
	"fmt"
	"path/filepath"
)

func StatusFunc() []FileStatusStruct {
	if len(ActiveFiles) == 0 {
		return nil
	}

	return treeFunc(FetchLatestCommitHash(), &ActiveFiles)
}

type FileStatusStruct struct {
	filename string
	status   string
	blobHash string
}

func treeFunc(treehash string, fileTrack *[]string) []FileStatusStruct {
	var changedFilesData []FileStatusStruct

	treeData, err := LsTreeFunc(treehash, []string{})
	if err != nil {
		fmt.Println("error ", err)
	}

	for _, file := range treeData {
		if file.fileType == "blob" {
			hash, err := hashObjectFunc(file.filename)
			if err != nil {
				fmt.Println("err ", err)
			}
			if hash != file.hex {
				changedFilesData = append(changedFilesData, FileStatusStruct{
					filename: file.filename,
					status:   "modified",
					blobHash: hash,
				})
				removeElement(fileTrack, file.filename)
			} else if hash == file.hex {
				removeElement(fileTrack, file.filename)
			}
		} else if file.fileType == "tree" {
			changedFilesData = append(changedFilesData, treeFunc(file.hex, fileTrack)...)
		}
	}

	if len(*fileTrack) > 0 {
		for _, file := range *fileTrack {
			hash, err := hashObjectFunc(file)
			if err != nil {
				fmt.Println("error : ", err)
			}
			changedFilesData = append(changedFilesData, FileStatusStruct{filename: filepath.Base(file), status: "new file", blobHash: hash})
		}
	}

	return changedFilesData
}

func removeElement(slice *[]string, element string) {
	for i, v := range *slice {
		if v == element {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}
