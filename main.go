package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jolson88/z-go/vm"
)

func main() {
	vm := vm.NewVirtualMachine("rom/zork2-r63-s860811.z3")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("z-go> ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				break
			}

			words := strings.SplitN(input, " ", 2)
			command := words[0]
			commandInput := ""
			if len(words) > 1 {
				commandInput = words[1]
			}

			switch command {

			case "header":
				fmt.Printf("Story header: %s\n", vm.PrintHeader())

			default:
				fmt.Println("Unknown command:", command, commandInput)
			}
		} else {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
