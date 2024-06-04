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

func (vm *VirtualMachine) PrintHeader() string {
	if !vm.initialized {
		return "VM not initialized"
	}

	return string(fmt.Sprintf("%v", vm.Memory[0:32]))
}
