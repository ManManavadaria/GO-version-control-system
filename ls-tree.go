package main

func LsTreeFunc(hash string, additional string, paths []string) (string, error) {
	file, err := CatfileFunc(hash)
	if err != nil {
		return "", err
	}

	treeHash, err := treeExtractor(file)
	if err != nil {
		return "", err
	}

	out, err := CatfileFunc(treeHash)
	if err != nil {
		return "", err
	}

	data := parseTreeObject([]byte(out), additional)

	return data, nil
}
