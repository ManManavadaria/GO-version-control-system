package main

import (
	"fmt"
	"os"
	"strings"

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
		statusData := command.StatusFunc(ActiveFiles)

		for _, filestatus := range statusData {
			helper.PrintOutput(fmt.Sprintf("%s: %s", filestatus.Status, filestatus.Filename))
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

	for _, arg := range args[2:] {
		if strings.HasPrefix(arg, "-") {
			cmd.Options = append(cmd.Options, arg)
		} else {
			cmd.Arguments = append(cmd.Arguments, arg)
		}
	}
	return cmd, nil
}
