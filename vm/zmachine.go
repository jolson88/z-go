package vm

import (
	"fmt"
	"log"
	"os"
)

type VirtualMachine struct {
	initialized bool
	Memory      [512 * 1024]byte
}

func NewVirtualMachine(pathToStoryFile string) *VirtualMachine {
	vm := &VirtualMachine{}

	data, err := os.ReadFile(pathToStoryFile)
	if err != nil {
		log.Fatalf("Error reading story file: %v", err)
	}

	copy(vm.Memory[:], data)

	vm.initialized = true
	return vm
}

func (vm *VirtualMachine) Run() {
	// Start the Z-Machine interpreter
}

func (vm *VirtualMachine) PrintHeader() {
	if !vm.initialized {
		fmt.Println("VM not initialized")
		return
	}

	fmt.Printf("%v\n\n", vm.Memory[0:32])
}

func (vm *VirtualMachine) PrintMemory(startAddress uint16, lines uint16) {
	const BYTES_PER_LINE uint16 = 16

	alignedStart := startAddress & 0xFFF0
	alignedEnd := (startAddress + BYTES_PER_LINE*lines) & 0xFFF0

	fmt.Println("        00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F")
	for i := alignedStart; i < alignedEnd; i += BYTES_PER_LINE {
		var offset uint16 = 0
		fmt.Printf("0x%04X: ", i+offset)
		for offset < BYTES_PER_LINE {
			fmt.Printf("%02X ", vm.Memory[i+offset])
			offset += 1
		}
		fmt.Println("")
	}
	fmt.Println("")
}
