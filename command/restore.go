package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/ManManavadaria/GO-version-control-system/helper"
)

func RestoreAll() {

}

func FetchRestoreFilesHex(files []string) []TreeDataStruct {
	hex := FetchLatestCommitHash()

	tree, err := LsTreeFuncAllFilesSearch(hex)
	if err != nil {
		helper.PrintError(err.Error())
	}

	var treeFiles []TreeDataStruct

	//NOTE: implement map converson approach
	for _, file := range files {
		for _, f := range tree {
			if file == f.Filename {
				treeFiles = append(treeFiles, f)
			}
		}
	}
	return treeFiles
}

func WriteHeadData(treeFiles []TreeDataStruct) {
	for _, file := range treeFiles {
		content, err := CatfileFunc(file.Hex)
		if err != nil {
			fmt.Println("err ", err)
		}

		if strings.HasPrefix(content, "\x00") {
			content = strings.TrimPrefix(content, "\x00")
		}

		f, err := os.OpenFile(file.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}

		f.Truncate(0)

		if _, err := f.WriteString(content); err != nil {
			fmt.Println("err ", err)
		}

		f.Close()
	}
}
