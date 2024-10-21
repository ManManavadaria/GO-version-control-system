package command

import "fmt"

func LsTreeFunc(hash string, paths []string) ([]TreeDataStruct, error) {
	file, err := CatfileFunc(hash)
	if err != nil {
		return nil, err
	}

	treeHash, err := treeExtractor(file)
	var content string
	if err != nil {
		content = file
	} else {
		content, err = CatfileFunc(treeHash)
		if err != nil {
			return nil, err
		}
	}

	treeContent := parseTreeObject([]byte(content))

	return treeContent, nil
}

func LsTreeFuncAllFilesSearch(hash string) ([]TreeDataStruct, error) {
	file, err := CatfileFunc(hash)
	if err != nil {
		return nil, err
	}

	treeHash, err := treeExtractor(file)
	if err != nil {
		return nil, err
	}

	var allfileshashData []TreeDataStruct

	treeContent := FetchRecursiveFiles(treeHash, allfileshashData, "")

	return treeContent, nil
}

func FetchRecursiveFiles(hash string, allfileshashData []TreeDataStruct, root string) []TreeDataStruct {
	content, _ := CatfileFunc(hash)
	treeContent := parseTreeObject([]byte(content))

	for _, file := range treeContent {
		if file.FileType == "tree" {
			allfileshashData = FetchRecursiveFiles(file.Hex, allfileshashData, fmt.Sprintf("%v%v\\", root, file.Filename))
		} else {
			file.Filename = fmt.Sprintf("%v%v", root, file.Filename)
			allfileshashData = append(allfileshashData, file)
		}
	}

	return allfileshashData
}
