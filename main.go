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
	case "init":
		msg, err := InitFunc()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31m\n%s\033[0m\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "\033[32m%v\033[0m\n", msg)

	case "cat-file":
		out, err := CatfileFunc()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "\033[32m%s\033[0m\n", out)

	case "hash-object":
		hash, err := hashObjectFuc()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "\033[32m%s\033[0m\n", hash)

		return
	default:
		fmt.Fprintf(os.Stderr, "\033[31mInvalid command.\033[0m\n")
	}
}
