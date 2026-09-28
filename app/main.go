package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var builtIn = map[string]struct{}{
	"type": {},
	"echo": {},
	"exit": {},
	"pwd":  {},
	"cd":   {},
}

// parseArguments tokenizes a command string into discrete arguments,
// properly handling single quotes, double quotes, and backslash escapes.
func parseArguments(input string) []string {
	var args []string
	var current strings.Builder
	inSingle := false
	inDouble := false
	escaped := false
	inToken := false

	for _, r := range input {
		if escaped {
			if inDouble {
				// Inside double quotes, \ only escapes $, `, ", \, and newline.
				switch r {
				case '"', '\\', '$', '`':
					current.WriteRune(r)
				default:
					current.WriteRune('\\')
					current.WriteRune(r)
				}
			} else {
				current.WriteRune(r)
			}
			escaped = false
			inToken = true
			continue
		}

		if inSingle {
			if r == '\'' {
				inSingle = false
			} else {
				current.WriteRune(r)
			}
			continue
		}

		if inDouble {
			switch r {
			case '"':
				inDouble = false
			case '\\':
				escaped = true
			default:
				current.WriteRune(r)
			}
			continue
		}

		// Unquoted context
		switch r {
		case '\'':
			inSingle = true
			inToken = true
		case '"':
			inDouble = true
			inToken = true
		case '\\':
			escaped = true
			inToken = true
		case ' ', '\t':
			if inToken {
				args = append(args, current.String())
				current.Reset()
				inToken = false
			}
		default:
			current.WriteRune(r)
			inToken = true
		}
	}

	// Flushing trailing escaped character if input ended with a dangling backslash
	if escaped {
		current.WriteRune('\\')
	}

	if inToken {
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
			} else if p, err := exec.LookPath(arg); err == nil {
				fmt.Printf("%s is %s\n", arg, p)
			} else {
				fmt.Printf("%s: not found\n", arg)
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
			fmt.Println("cd: too many arguments")
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
		} else if strings.HasPrefix(args[0], "~/") {
			dir, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("Failed to get home directory: %v\n", err)
				return
			}
			targetDir = filepath.Join(dir, args[0][2:])
		} else {
			targetDir = args[0]
		}

		cleanedPath := filepath.Clean(targetDir)
		if !filepath.IsAbs(cleanedPath) {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Failed to get current working directory: %v\n", err)
				return
			}
			cleanedPath = filepath.Join(wd, cleanedPath)
		}

		if err := os.Chdir(cleanedPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("cd: %s: No such file or directory\n", targetDir)
			} else {
				fmt.Printf("cd: %v\n", err)
			}
		}
	default:
		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				// Command executed but returned non-zero exit code
				return
			}
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
			// Handles EOF (Ctrl+D) gracefully
			break
		}

		userInput := strings.TrimRight(readString, "\r\n")
		if len(strings.TrimSpace(userInput)) == 0 {
			continue
		}

		tokens := parseArguments(userInput)
		if len(tokens) == 0 {
			continue
		}

		handle(tokens[0], tokens[1:])
	}
}
