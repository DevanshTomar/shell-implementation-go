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
	"type": {},
	"echo": {},
	"exit": {},
	"pwd":  {},
	"cd":   {},
}

// parseArguments tokenizes a command string into discrete arguments,
// properly honoring single and double quotes while preserving spaces inside them.
func parseArguments(input string) []string {
	var args []string
	var current strings.Builder
	inSingleQuotes := false
	inDoubleQuotes := false
	hasToken := false

	for _, r := range input {
		if inSingleQuotes {
			if r == '\'' {
				inSingleQuotes = false
			} else {
				current.WriteRune(r)
				hasToken = true
			}
		} else if inDoubleQuotes {
			if r == '"' {
				inDoubleQuotes = false
			} else {
				current.WriteRune(r)
				hasToken = true
			}
		} else {
			switch r {
			case '\'':
				inSingleQuotes = true
				hasToken = true
			case '"':
				inDoubleQuotes = true
				hasToken = true
			case ' ', '\t':
				if hasToken {
					args = append(args, current.String())
					current.Reset()
					hasToken = false
				}
			default:
				current.WriteRune(r)
				hasToken = true
			}
		}
	}
	if hasToken {
		args = append(args, current.String())
	}
	return args
}

func handle(command string, args []string) {
	switch command {
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "type":
		if len(args) == 0 {
			fmt.Println("type: command requires at least one argument")
			return
		}
		for _, arg := range args {
			if _, ok := builtIn[arg]; ok {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else {
				if p, err := exec.LookPath(arg); err == nil {
					fmt.Printf("%s is %s\n", arg, p)
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
		if len(args) > 1 {
			fmt.Printf("cd: too many arguments\n")
			return
		}

		targetDir := ""
		if len(args) == 0 || args[0] == "~" {
			dir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("Failed to get home directory: %v\n", err)
				return
			}
			targetDir = dir
		} else {
			targetDir = args[0]
		}

		cleanedPath := path.Clean(targetDir)
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
				fmt.Printf("cd: %s: No such file or directory\n", targetDir)
			} else {
				fmt.Printf("Failed to change directory: %v\n", err)
			}
		}
	default:
		if _, err := exec.LookPath(command); err == nil {
			cmd := exec.Command(command, args...)
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
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		readString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
			os.Exit(1)
		}

		userInput := strings.TrimSpace(readString)
		if len(userInput) == 0 {
			continue
		}

		tokens := parseArguments(userInput)
		if len(tokens) == 0 {
			continue
		}

		command := tokens[0]
		arguments := tokens[1:]

		handle(command, arguments)
	}
}
