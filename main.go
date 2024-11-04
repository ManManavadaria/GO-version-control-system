package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ManManavadaria/GO-version-control-system/command"
	"github.com/ManManavadaria/GO-version-control-system/helper"
)

func main() {
	cmd, err := ParseCommand(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\033[31m%v\033[0m\n", err)
		return
	}

	switch cmd.Name {
	case "--version":
		helper.PrintOutput("go-vcs version 0.0.1")
		return
	case "init":
		helper.PrintOutput("Initializing the repository...")
		msg, err := command.InitFunc()
		if err != nil {
			helper.PrintError(err.Error())
			return
		}
		helper.PrintOutput(msg)

	case "cat-file":
		if cmd.Length < 3 {
			helper.PrintError("invalid arguments, SHA missing.")
		}

		if len(cmd.Arguments) <= 0 {
			helper.PrintError("invalid arguments, SHA missing.")
		}
		sha := cmd.Arguments[0]

		out, err := command.CatfileFunc(sha)
		if err != nil {
			helper.PrintError(err.Error())
		}
		helper.PrintOutput(out)

	case "hash-object":
		if cmd.Length != 4 {
			helper.PrintError("Invalid Arguments.")
			return
		}

		if len(cmd.Arguments) <= 0 {
			helper.PrintError("Invalid Arguments.")
			return
		}

		fileName := cmd.Arguments[0]
		hash, err := command.HashObjectFunc(fileName)
		if err != nil {
			helper.PrintError(err.Error())
		}
		helper.PrintOutput(hash)
		return

	case "ls-tree":
		if cmd.IsInitialCommit {
			return
		}
		if cmd.Length < 3 {
			helper.PrintError("invalid arguments, SHA missing.")
		}

		var hash string
		var paths []string
		var additional string

		if hash = cmd.Arguments[0]; hash == "HEAD" {
			hash = command.FetchLatestCommitHash()
		}
		if len(cmd.Options) > 0 {
			additional = cmd.Options[0]
		}
		if len(cmd.Arguments) > 1 {
			paths = append(paths, cmd.Arguments[1:]...)
		}

		data, err := command.LsTreeFunc(hash, paths)
		if err != nil {
			helper.PrintError(err.Error())
		}

		var out string
		if additional == "--name-only" {
			for _, filedata := range data {
				out += filedata.Filename + "\n"
			}
		} else {
			for _, filedata := range data {
				out += filedata.Mode + " " + filedata.FileType + " " + filedata.Hex + "    " + filedata.Filename + "\n"
			}
		}
		helper.PrintOutput(out)
		return

	case "status":
		if idx := command.LoadIndex(); cmd.IsInitialCommit && len(idx.Entries) > 0 {
			helper.PrintOutput("Initial Staged files +++++++++\n")
			for _, entry := range idx.Entries {
				helper.PrintOutput(fmt.Sprintf("%s", entry.Path))
			}
			helper.PrintOutput("\n+++++++++\n\n")
		} else if stagedFiles, _ := command.StagedFiles(); !cmd.IsInitialCommit && len(stagedFiles) > 0 {
			if len(stagedFiles) > 0 {
				helper.PrintOutput("Staged files +++++++++\n")
				for _, file := range stagedFiles {
					helper.PrintOutput(fmt.Sprintf("%s", file.Filename))
				}
				helper.PrintOutput("\n+++++++++\n\n")
			}
		} else {
			helper.PrintInfo("Staging is empty...\n")
		}

		statusData := command.StatusFunc(ActiveFiles)
		if len(statusData) == 0 {
			helper.PrintInfo("Working directory is ideal")
		}

		helper.PrintOutput("Changes not staged for commit ---------\n")
		for _, filestatus := range statusData {
			if filestatus.Status == "modified" {
				helper.PrintInfo(fmt.Sprintf("%s: %s", filestatus.Status, filestatus.Filename))
			} else if filestatus.Status == "new file" {
				helper.PrintOutput(fmt.Sprintf("%s: %s", filestatus.Status, filestatus.Filename))
			} else if filestatus.Status == "removed" {
				helper.PrintDeleted(fmt.Sprintf("%s: %s", filestatus.Status, filestatus.Filename))
			}
		}
		helper.PrintOutput("\n---------")
	case "add":
		err := ValidateFileOptionArgument(cmd.Arguments)
		if err != nil {
			helper.PrintError(err.Error())
		}
		command.InitIndex()

		if len(cmd.Arguments) == 1 && cmd.Arguments[0] == "." {
			var files []string
			if cmd.IsInitialCommit {
				files = ActiveFiles
			} else {
				statusData := command.StatusFunc(ActiveFiles)
				for _, data := range statusData {
					files = append(files, data.Filename)
				}
			}
			command.UpdateIndex(files)
		} else if len(cmd.Arguments) >= 1 {
			var files []string
			for _, file := range cmd.Arguments {
				file = filepath.FromSlash(file)
				f, _ := os.Stat(file)
				if f.IsDir() {
					chFiles := GetAllFiles(fmt.Sprintf("./%s", file))
					files = append(files, chFiles...)
				} else {
					files = append(files, file)
				}
			}
			command.UpdateIndex(files)
		}

		if cmd.IsInitialCommit {
			command.WriteTree()
		}
	case "restore":
		err := ValidateFileOptionArgument(cmd.Arguments)
		if err != nil {
			helper.PrintError(err.Error())
		}
		var treeFiles []command.TreeDataStruct
		if len(cmd.Arguments) == 1 && cmd.Arguments[0] == "." {
			treeFiles = command.FetchRestoreFilesHex(ActiveFiles)
		} else {
			treeFiles = command.FetchRestoreFilesHex(cmd.Arguments)
		}

		if cmd.Options[0] == "--staged" {
			command.UnstageFilesFromIndex(treeFiles)
		} else {
			command.WriteHeadData(treeFiles)
		}
	case "write-tree":
		hash := command.WriteTree()
		helper.PrintSuccess(hash)

	case "commit":
		if len(os.Args) > 4 {
			helper.PrintError("Invalid arguments.")
		}

		commit := command.NewCommitConfig()
		if len(cmd.Options) > 0 {
			for _, option := range cmd.Options {
				if option == "-a" {
					statusData := command.StatusFunc(ActiveFiles)
					var files []string
					for _, data := range statusData {
						files = append(files, data.Filename)
					}
					command.UpdateIndex(files)
				} else if option == "-m" {
					commit.CommitMsg = cmd.Arguments[0]
				}
			}
		}
		if !cmd.IsInitialCommit {
			stagedFiles, _ := command.StagedFiles()
			if len(stagedFiles) <= 0 {
				helper.PrintInfo("Staging is ideal, Please stage changes to complete the commit process.")
				return
			}
		}

		commit.CurrentTreeHash = command.WriteTree()
		if cmd.IsInitialCommit {
			command.CreateInitialBranch(command.CurrentBranchName())
			command.CreateLogsHEAD()
		} else {
			commit.ParentCommitHash = command.FetchLatestCommitHash()
		}

		err := commit.CreateCommitObject()
		if err != nil {
			helper.PrintError(err.Error())
		}

		commit.AddCommitStr(".go-vcs/logs/HEAD", "commit")
		if cmd.IsInitialCommit {
			commit.AddCommitStr(command.FetchBranchLogFileAddr(), "commit (initial)")
		} else {
			commit.AddCommitStr(command.FetchBranchLogFileAddr(), "commit")
		}
		commit.UpdateCommitHash(command.FetchBranchHeadFileAddr())

	case "branch":
		if cmd.IsInitialCommit {
			helper.PrintOutput(command.CurrentBranchName())
			return
		}
		files := command.ListAllBranch()
		if len(cmd.Options) == 0 && len(cmd.Arguments) == 0 {
			current := command.CurrentBranchName()
			for _, file := range files {
				if file.Name() == current {
					helper.PrintSuccess("* " + file.Name())
				} else {
					helper.PrintSuccess(file.Name())
				}
			}
		}

		if len(cmd.Options) == 0 && len(cmd.Arguments) == 1 {
			command.CreateNewBranch(cmd.Arguments[0])
		}

		if len(cmd.Options) == 1 && len(cmd.Arguments) == 1 {
			switch cmd.Options[0] {
			case "-m":
				command.RenameCurrentBranch(cmd.Arguments[0])
				// case "-d":
				// case "-D":
			}
		}
		if len(cmd.Options) == 1 && len(cmd.Arguments) == 2 {
			switch cmd.Options[0] {
			case "-m":
				command.RenameBranch(cmd.Arguments[0], cmd.Arguments[1])
				// case "-d":
				// case "-D":
			}
		}
	case "checkout":
		if cmd.IsInitialCommit {
			if len(cmd.Arguments) == 1 {
				command.UpdateRefPath(fmt.Sprintf("refs/heads/%s", cmd.Arguments[0]))
				return
			}
		}
		if len(command.StatusFunc(ActiveFiles)) > 0 {
			helper.PrintInfo("Changes in working directory are not commited yet")
		}

		if len(cmd.Options) == 0 && len(cmd.Arguments) == 1 {
			if ok := command.CheckExistingBranch(cmd.Arguments[0]); !ok {
				helper.PrintInfo("Incorect branch name")
			}
			hash := command.SwitchBranch(cmd.Arguments[0])
			command.ReplaceCommitContent(hash)
		}

		if len(cmd.Options) == 1 && len(cmd.Arguments) == 1 {
			switch cmd.Options[0] {
			case "-d":
				command.CreateNewBranch(cmd.Arguments[0])
				hash := command.SwitchBranch(cmd.Arguments[0])
				command.ReplaceCommitContent(hash)
			}
		}
	case "log":
		if cmd.IsInitialCommit {
			return
		}
		commits := command.FetchCurrentBranchLogs()
		fmt.Println("length ", len(commits))
		if len(cmd.Options) == 1 && cmd.Options[0] == "--oneline" {
			for i, commit := range commits {
				if i == 0 {
					helper.PrintInfo(fmt.Sprintf("commit %s (HEAD -> %s) - ", commit.CurrentCommitHash[0:7], command.CurrentBranchName()))
				} else {
					helper.PrintInfo(fmt.Sprintf("commit %s ", commit.CurrentCommitHash[0:7]))
				}
				helper.Print(fmt.Sprintf("%v", commit.CommitMsg))
			}
		} else {
			for i, commit := range commits {
				t := time.Unix(commit.Timestamp, 0)
				if i == 0 {
					helper.PrintInfo(fmt.Sprintf("commit %s (HEAD -> %s)", commit.CurrentCommitHash, command.CurrentBranchName()))
				} else {
					helper.PrintInfo(fmt.Sprintf("commit %s", commit.CurrentCommitHash))
				}
				helper.Print(fmt.Sprintf("Author: %s <%s>\nDate: %v %v\n\n    %s\n", commit.AuthorName, commit.AuthorEmail, t.Format(time.RFC1123), commit.TimeZone, commit.CommitMsg))
			}
		}
	default:
		helper.PrintError("Invalid command.")
	}
}

type Command struct {
	Name            string
	Length          int
	Options         []string
	Arguments       []string
	IsInitialCommit bool
}

func ParseCommand(args []string) (Command, error) {
	var cmd Command
	cmd.Length = len(args)
	cmd.Options = []string{}
	cmd.Arguments = []string{}
	cmd.IsInitialCommit = false

	if len(args) < 2 {
		return Command{}, fmt.Errorf("No command provided. Use 'govcs <command> [options] [arguments]'.")
	}

	cmd.Name = args[1]

	if cmd.Name != "init" {
		_, err := os.Stat(".go-vcs")
		if errors.Is(err, os.ErrNotExist) {
			helper.PrintError("Failed to find .go-vcs directory, Repository not availables")
		}

		_, err = os.Stat(".go-vcs/logs/HEAD")
		if errors.Is(err, os.ErrNotExist) {
			cmd.IsInitialCommit = true
		}
	}

	for _, arg := range args[2:] {
		if strings.HasPrefix(arg, "-") {
			cmd.Options = append(cmd.Options, arg)
		} else {
			cmd.Arguments = append(cmd.Arguments, arg)
		}
	}
	return cmd, nil
}
