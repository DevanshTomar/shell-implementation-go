package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var builtIn map[string]struct{} = map[string]struct{}{
	"type": struct{}{},
	"echo": struct{}{},
	"exit": struct{}{},
	"pwd":  struct{}{},
	"cd":   struct{}{},
}

func handle(command, arguments string) {
	switch command {
	case "echo":
		fmt.Println(arguments)
	case "type":
		if len(arguments) == 0 {
			fmt.Println("type: command requires at least one argument")
			return
		}

		args := strings.Split(arguments, " ")
		for _, arg := range args {
			if _, ok := builtIn[arg]; ok {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else {
				if path, err := exec.LookPath(arg); err == nil {
					fmt.Printf("%s is %s\n", arg, path)
				} else {
					fmt.Printf("%s: not found\n", arg)
				}

			}
		}
	case "exit":
		os.Exit(0)
	case "pwd":
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get PWD: %v\n", err)
			return
		}
		fmt.Println(wd)
	case "cd":
		args := strings.Split(arguments, " ")
		if len(args) > 1 {
			fmt.Printf("cd: too many arguments\n")
			return
		}

		if len(args) == 0 || arguments == "~" {
			dir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("Failed to get home directory: %v\n", err)
				return
			}
			args[0] = dir
		}

		cleanedPath := path.Clean(args[0])
		if !path.IsAbs(cleanedPath) {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Failed to get current working directory: %v\n", err)
				return
			}
			cleanedPath = path.Join(wd, cleanedPath)
		}

		if err := os.Chdir(cleanedPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("cd: %s: No such file or directory\n", args[0])
			} else {
				fmt.Printf("Failed to change directory: %v\n", err)
			}
		}
	default:
		if _, err := exec.LookPath(command); err == nil {
			argList := strings.Split(arguments, " ")
			cmd := exec.Command(command, argList...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("%s: command ran with error: %s\n", command, err)
				return
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}

	}
}

func main() {
	for {
		fmt.Print("$ ")
		readString, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
			os.Exit(1)
		}

		userInput := strings.TrimSpace(readString[:len(readString)-1])
		command, arguments, _ := strings.Cut(userInput, " ") // extracting the command and arguments

		handle(command, arguments)
	}
}
