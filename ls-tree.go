package main

func LsTreeFunc(hash string, paths []string) ([]TreeDataStruct, error) {
	file, err := CatfileFunc(hash)
	if err != nil {
		return nil, err
	}

	treeHash, err := treeExtractor(file)
	if err != nil {
		return nil, err
	}

	out, err := CatfileFunc(treeHash)
	if err != nil {
		return nil, err
	}

	treeContent := parseTreeObject([]byte(out))

	return treeContent, nil
}
