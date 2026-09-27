package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var builtIn map[string]struct{} = map[string]struct{}{
	"type": struct{}{},
	"echo": struct{}{},
	"exit": struct{}{},
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
				fmt.Printf("%s: not found\n", arg)
			}
		}
	case "exit":
		os.Exit(0)
	default:
		fmt.Printf("%s: command not found\n", command)
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
