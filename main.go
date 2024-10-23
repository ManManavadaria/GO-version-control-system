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
		command.InitIndex()
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
		stagedFiles, ok := command.StagedFiles()

		if !ok {
			for _, file := range ActiveFiles {
				helper.PrintInfo(file)
			}
			return
		}

		if len(stagedFiles) > 0 {

			helper.PrintOutput("Staged files +++++++++\n")
			for _, file := range stagedFiles {
				helper.PrintOutput(fmt.Sprintf("%s", file.Filename))
			}
			helper.PrintOutput("\n+++++++++\n\n")
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
		statusData := command.StatusFunc(ActiveFiles)

		if len(cmd.Arguments) == 1 && cmd.Arguments[0] == "." {
			var files []string

			for _, data := range statusData {
				files = append(files, data.Filename)
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

		var commit command.CommitConfig
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

		stagedFiles, _ := command.StagedFiles()
		if len(stagedFiles) <= 0 {
			helper.PrintInfo("Staging is ideal, Please stage changes to complete the commit process.")
			return
		}

		//NOTE: create a initcommit function to generate a base config struct
		commit.CurrentTreeHash = command.WriteTree()
		commit.ParentCommitHash = command.FetchLatestCommitHash()

		commit.AuthorName = "ManPatel"
		commit.AuthorEmail = "mam@gmail.com"
		commit.Timestamp = time.Now().Unix()
		_, tzOffset := time.Now().Zone()
		tzHours := tzOffset / 3600
		tzMinutes := (tzOffset % 3600) / 60

		commit.TimeZone = fmt.Sprintf("%+03d%02d", tzHours, tzMinutes)

		err := commit.CreateCommitObject()
		if err != nil {
			helper.PrintError(err.Error())
		}

		commit.AddCommitStr(".go-vcs/logs/HEAD", "commit")
		commit.AddCommitStr(command.FetchBranchLogFileAddr(), "commit")
		commit.UpdateCommitHash(command.FetchBranchHeadFileAddr())

	case "branch":
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
	case "checkout":
		if len(command.StatusFunc(ActiveFiles)) > 0 {
			helper.PrintInfo("Changes in working directory are not commited yet")
		}
		if len(cmd.Options) == 0 && len(cmd.Arguments) == 1 {
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

	default:
		helper.PrintError("Invalid command.")
	}
}

type Command struct {
	Name      string
	Length    int
	Options   []string
	Arguments []string
}

func ParseCommand(args []string) (Command, error) {
	var cmd Command
	cmd.Length = len(args)
	cmd.Options = []string{}
	cmd.Arguments = []string{}

	if len(args) < 2 {
		return Command{}, fmt.Errorf("No command provided. Use 'govcs <command> [options] [arguments]'.")
	}

	cmd.Name = args[1]

	if cmd.Name != "init" {
		_, err := os.Stat(".go-vcs")
		if errors.Is(err, os.ErrNotExist) {
			helper.PrintError("Failed to find .go-vcs directory, Repository not availables")
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
