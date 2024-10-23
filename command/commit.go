package command

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/ManManavadaria/GO-version-control-system/helper"
)

type CommitConfig struct {
	ParentCommitHash  string
	CurrentTreeHash   string
	CurrentCommitHash string
	CommitMsg         string
	AuthorName        string
	AuthorEmail       string
	Timestamp         int64
	TimeZone          string
}

func (c *CommitConfig) CreateCommitObject() error {

	content := c.GenerateCommitObjectContent()

	header := fmt.Sprintf("commit %d\000", len(content))

	fullObject := append([]byte(header), content...)

	commitHash := sha1.Sum(fullObject)

	hashStr := fmt.Sprintf("%x", commitHash)

	var compressObject bytes.Buffer

	w := zlib.NewWriter(&compressObject)

	_, err := w.Write(fullObject)
	if err != nil {
		return err
	}
	w.Close()

	objectDir := fmt.Sprintf(".go-vcs/objects/%s", hashStr[:2])
	os.MkdirAll(objectDir, os.ModePerm)
	objectPath := fmt.Sprintf("%s/%s", objectDir, hashStr[2:])
	err = os.WriteFile(objectPath, compressObject.Bytes(), 0644)
	if err != nil {
		return err
	}
	c.CurrentCommitHash = hashStr
	return nil
}

func (c *CommitConfig) GenerateCommitObjectContent() []byte {

	var commitContent bytes.Buffer

	commitContent.WriteString(fmt.Sprintf("tree %s\n", c.CurrentTreeHash))

	if c.ParentCommitHash != "" {
		commitContent.WriteString(fmt.Sprintf("parent %s\n", c.ParentCommitHash))
	}

	commitContent.WriteString(fmt.Sprintf("author %s <%s> %d +0000 \n", c.AuthorName, c.AuthorEmail, c.Timestamp))
	commitContent.WriteString(fmt.Sprintf("commiter %s <%s> %d +0000 \n", c.AuthorName, c.AuthorEmail, c.Timestamp))

	if c.CommitMsg != "" {
		commitContent.WriteString(c.CommitMsg + "\n")
	}

	return commitContent.Bytes()
}

func (c *CommitConfig) AddCommitStrInHead() {

	commitStr := fmt.Sprintf("%s %s %s <%s> %d %s    %s\n", c.ParentCommitHash, c.CurrentCommitHash, c.AuthorName, c.AuthorEmail, c.Timestamp, c.TimeZone, fmt.Sprintf("commit: %s", c.CommitMsg))

	file, err := os.OpenFile(".go-vcs/logs/HEAD", os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(commitStr)
}

func (c *CommitConfig) AddCommitStr(path string, strType string) {
	commitStr := fmt.Sprintf("%s %s %s <%s> %d %s    %s\n", c.ParentCommitHash, c.CurrentCommitHash, c.AuthorName, c.AuthorEmail, c.Timestamp, c.TimeZone, fmt.Sprintf("%s: %s", strType, c.CommitMsg))

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(commitStr)
}

func (c *CommitConfig) UpdateCommitHash(path string) {

	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Truncate(0)

	file.WriteString(c.CurrentCommitHash)
}

func ReplaceCommitContent(hash string) {
	treeData, err := LsTreeFuncAllFilesSearch(hash)
	if err != nil {
		helper.PrintError(err.Error())
	}

	WriteHeadData(treeData)
}
