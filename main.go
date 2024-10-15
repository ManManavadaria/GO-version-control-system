package main

import (
	"fmt"
	"os"
)

func main() {
	argLen := len(os.Args)

	if argLen < 2 {
		fmt.Fprintf(os.Stderr, "\033[31mInsufficient arguments, please input the arguments.\033[0m\n")
		return
	}

	switch cmd := os.Args[1]; cmd {
	case "--version":
		fmt.Fprintf(os.Stderr, "go-vcs version 0.0.1")
		os.Exit(1)
	case "init":
		msg, err := InitFunc()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31m\n%s\033[0m\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "\033[32m%v\033[0m\n", msg)

	case "cat-file":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", "invalid arguments, SHA missing.")
			os.Exit(1)
		}

		sha := os.Args[2]

		out, err := CatfileFunc(sha)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "\033[32m%s\033[0m\n", out)

	case "hash-object":
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", "Invalid Arguments.")
		}

		fileName := os.Args[3]
		hash, err := hashObjectFunc(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "\033[32m%s\033[0m\n", hash)

		return
	case "ls-tree":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", "invalid arguments, SHA missing.")
			os.Exit(1)
		}

		var hash string
		var paths []string
		var additional string

		if len(os.Args) == 3 {
			if os.Args[2] == "HEAD" {
				hash = FetchLatestCommitHash()
			} else {
				if len(os.Args[2]) == 40 {
					hash = os.Args[2]
				} else {
					fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", "invalid arguments")
					os.Exit(1)
				}
			}
		} else {
		loop:
			for i, arg := range os.Args {
				if i > 1 {
					switch arg {
					case "HEAD":
						hash = FetchLatestCommitHash()
						paths = append(paths, os.Args[i+1:]...)
						break loop
					case "--name-only":
						additional = arg
						continue
					default:
						if len(arg) == 40 {
							hash = arg
							paths = append(paths, os.Args[i+1:]...)
							break loop
						} else {
							fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", "invalid arguments")
							os.Exit(1)
						}
					}
				}
			}
		}

		// fmt.Println(hash)
		// fmt.Println(paths)
		// fmt.Println(additional)
		data, err := LsTreeFunc(hash, additional, paths)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "\033[32m%s\033[0m\n", data)

		return
	case "status":
		// FetchIgnoreFiles()
		// GetAllFiles()
	default:
		fmt.Fprintf(os.Stderr, "\033[31mInvalid command.\033[0m\n")
	}
}
