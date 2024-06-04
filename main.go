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

	fmt.Println("\nWelcome to z-go! Type 'exit' to quit.")
	for {
		fmt.Print("z-go> ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				break
			}

			words := strings.SplitN(input, " ", 2)
			command := words[0]
			parameters := ""
			if len(words) > 1 {
				parameters = words[1]
			}

			fmt.Printf("\n")
			switch command {

			case "ascii":
				var charCode byte
				_, err := fmt.Sscanf(parameters, "%X", &charCode)
				if err != nil {
					fmt.Println("Invalid character code:", words[1])
					fmt.Println("Usage: ascii <hex code>")
					break
				}
				fmt.Printf("ASCII Character: '%c'\n", charCode)

			case "header":
				vm.PrintHeader()

			case "mem":
				var address, lineCount uint16
				lineCount = 8
				n, err := fmt.Sscanf(parameters, "%X %d", &address, &lineCount)
				if err != nil && n < 1 {
					fmt.Println("Invalid start address:", words[1])
					fmt.Println("Usage: mem <hex address> <optional line count>")
					break
				}
				vm.PrintMemory(address, lineCount)

			default:
				fmt.Println("Unknown command:", command)
			}
			fmt.Printf("\n")
		} else {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
