package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtIn map[string]struct{} = map[string]struct{}{
	"type": struct{}{},
	"echo": struct{}{},
	"exit": struct{}{},
	"pwd":  struct{}{},
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
