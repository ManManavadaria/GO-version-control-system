package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ManManavadaria/GO-version-control-system/helper"
)

func FetchLatestCommitHash() string {

	headFile, err := os.Open(".go-vcs/HEAD")
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}
	if err != nil {
		helper.PrintError(err.Error())
	}
	defer headFile.Close()

	b, err := io.ReadAll(headFile)
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}
	if err != nil {
		helper.PrintError(err.Error())
	}

	ref := strings.Split(strings.TrimSpace(string(b)), " ")[1]

	refBranch, err := os.Open(fmt.Sprintf(".go-vcs/%s", ref))
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}
	if err != nil {
		helper.PrintError(err.Error())
	}
	defer refBranch.Close()

	hashByte, err := io.ReadAll(refBranch)
	if errors.Is(err, os.ErrNotExist) {
		return ""
	}
	if err != nil {
		helper.PrintError(err.Error())
	}

	return string(strings.TrimSpace(string(hashByte)))
}

func CurrentBranchName() string {
	b, err := os.ReadFile(".go-vcs/HEAD")
	if err != nil {
		helper.PrintError(err.Error())
	}

	s := strings.Split(strings.TrimSpace(string(b)), "/")

	return s[len(s)-1]
}
func FetchBranchHeadFileAddr() string {
	b, err := os.ReadFile(".go-vcs/HEAD")
	if err != nil {
		helper.PrintError(err.Error())
	}

	path := strings.TrimPrefix(strings.TrimSpace(string(b)), "ref: ")

	return fmt.Sprintf(".go-vcs/%s", path)
}

func FetchBranchLogFileAddr() string {
	b, err := os.ReadFile(".go-vcs/HEAD")
	if err != nil {
		helper.PrintError(err.Error())
	}

	path := strings.TrimPrefix(strings.TrimSpace(string(b)), "ref: ")

	return fmt.Sprintf(".go-vcs/logs/%s", path)
}

func UpdateRefPath(path string) {
	f, err := os.OpenFile(".go-vcs/HEAD", os.O_WRONLY, 0666)
	if err != nil {
		helper.PrintError(err.Error())
	}

	defer f.Close()

	f.Truncate(0)

	f.WriteString(fmt.Sprintf("ref: %s", path))
}
