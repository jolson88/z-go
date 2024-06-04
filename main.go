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
			fmt.Printf("\n")
			switch command {

			case "header":
				vm.PrintHeader()

			case "mem":
				vm.PrintMemory(0, 8)

			default:
				fmt.Println("Unknown command:", command)
			}
			fmt.Printf("\n\n")
		} else {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
