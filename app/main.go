package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	for {
		fmt.Print("$ ")
		readString, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
			os.Exit(1)
		}

		command := readString[:len(readString)-1]

		switch command {
		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
