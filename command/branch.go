package command

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/ManManavadaria/GO-version-control-system/helper"
)

func ListAllBranch() []fs.DirEntry {

	files, err := os.ReadDir(".go-vcs/refs/heads")
	if err != nil {
		fmt.Println("Error : ", err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println("Error : folder available in the heads directory", file.Name())
		}
	}

	return files
}

func RenameCurrentBranch(newName string) {
	branchHead := FetchBranchHeadFileAddr()
	branchLog := FetchBranchLogFileAddr()

	newHeadPath := fmt.Sprintf(".go-vcs/refs/heads/%s", newName)
	newLogPath := fmt.Sprintf(".go-vcs/logs/refs/heads/%s", newName)

	if err := os.Rename(branchHead, newHeadPath); err != nil {
		helper.PrintError(err.Error())
	}

	if err := os.Rename(branchLog, newLogPath); err != nil {
		helper.PrintError(err.Error())
	}

	UpdateRefPath(fmt.Sprintf("refs/heads/%s", newName))
}

func CreateNewBranch(name string) {

	newHeadPath := fmt.Sprintf(".go-vcs/refs/heads/%s", name)
	newLogPath := fmt.Sprintf(".go-vcs/logs/refs/heads/%s", name)

	logFile, err := os.Create(newLogPath)
	if err != nil {
		helper.PrintError(err.Error())
	}
	logFile.Close()
	writeFirstCommitToNewBranch(newLogPath)

	headFile, err := os.Create(newHeadPath)
	if err != nil {
		helper.PrintError(err.Error())
	}
	defer headFile.Close()

	headFile.WriteString(FetchLatestCommitHash())
}

func writeFirstCommitToNewBranch(filePath string) {
	commit := CommitConfig{
		ParentCommitHash:  "0000000000000000000000000000000000000000",
		CurrentTreeHash:   "",
		CurrentCommitHash: FetchLatestCommitHash(),
		CommitMsg:         fmt.Sprintf("Created from %s", CurrentBranchName()),
	}

	commit.AuthorName = "ManPatel"
	commit.AuthorEmail = "mam@gmail.com"
	commit.Timestamp = time.Now().Unix()
	_, tzOffset := time.Now().Zone()
	tzHours := tzOffset / 3600
	tzMinutes := (tzOffset % 3600) / 60

	commit.TimeZone = fmt.Sprintf("%+03d%02d", tzHours, tzMinutes)

	commit.AddCommitStr(filePath, "branch")
}

func SwitchBranch(branch string) string {
	UpdateRefPath(fmt.Sprintf("refs/heads/%s", branch))

	return FetchLatestCommitOfBranch(fmt.Sprintf(".go-vcs/refs/heads/%s", branch))
}

func FetchLatestCommitOfBranch(path string) string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("err ", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("err ", err)
	}

	return strings.TrimSpace(string(b))
}
