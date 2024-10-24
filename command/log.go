package command

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ManManavadaria/GO-version-control-system/helper"
)

func FetchCurrentBranchLogs() []CommitConfig {
	addr := FetchBranchLogFileAddr()

	fmt.Println("addr : ", addr)
	commits, err := FetchBranchLogs(addr)
	if err != nil {
		helper.PrintError(err.Error())
	}

	return commits
}

func FetchBranchLogs(logFilePath string) ([]CommitConfig, error) {
	data, err := os.ReadFile(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %v", err)
	}

	// fmt.Println("data =>", string(data))

	lines := strings.Split(string(data), "\n")

	var commitLogs []CommitConfig

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			fmt.Println("continue empty")
			continue
		}

		parts := strings.SplitN(line, " ", 7)
		if len(parts) < 7 {
			fmt.Println("continue 7")
			continue
		}

		parentHash := parts[0]
		currentHash := parts[1]

		authorName := parts[2]
		authorEmail := parts[3]

		timestamp, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %v", err)
		}
		timezone := parts[5]

		commitMsg := strings.TrimSpace(parts[6])

		commitConfig := CommitConfig{
			ParentCommitHash:  parentHash,
			CurrentCommitHash: currentHash,
			AuthorName:        authorName,
			AuthorEmail:       authorEmail,
			Timestamp:         timestamp,
			TimeZone:          timezone,
			CommitMsg:         commitMsg,
		}

		commitLogs = append(commitLogs, commitConfig)
	}
	return commitLogs, nil
}

func parseTimestamp(timestampStr string) (int64, error) {
	timestamp, err := time.Parse("20060102150405", timestampStr)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp format: %s", timestampStr)
	}
	return timestamp.Unix(), nil
}
