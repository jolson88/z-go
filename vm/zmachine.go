package vm

import (
	"fmt"
	"log"
	"os"
)

const KILOBYTES = 1024

type VirtualMachine struct {
	initialized bool
	Memory      [512 * KILOBYTES]byte // The maximum size of any z-machine file
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

/*
Memory Read/Write
*/
func (vm *VirtualMachine) PrintMemory(startAddress uint16, lines uint16) {
	if !vm.initialized {
		fmt.Println("VM not initialized")
		return
	}

	const BYTES_PER_LINE uint16 = 16

	alignedStart := startAddress & 0xFFF0
	alignedEnd := (startAddress + BYTES_PER_LINE*lines) & 0xFFF0

	fmt.Println("        00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F")
	for i := alignedStart; i < alignedEnd; i += BYTES_PER_LINE {
		var offset, cellAddress uint16
		fmt.Printf("0x%04X: ", i)
		for offset = 0; offset < BYTES_PER_LINE; offset += 1 {
			cellAddress = i + offset
			if cellAddress == startAddress {
				fmt.Printf("\033[32m%02X\033[0m ", vm.Memory[cellAddress])
			} else {
				fmt.Printf("%02X ", vm.Memory[cellAddress])
			}
		}
		fmt.Printf("\n")
	}
}

func (vm *VirtualMachine) readWord(address uint16) uint16 {
	// Z-Machine is Big Endian. Need to swap around the bytes on Little Endian systems.
	return uint16(vm.Memory[address])<<8 | uint16(vm.Memory[address+1])
}

/*
Header Information
*/
func (vm *VirtualMachine) StoryChecksum() uint16 {
	return vm.readWord(0x1C)
}

func (vm *VirtualMachine) StoryLength() uint32 {
	// Up to v3, the story length this value multiplied by 2.
	// See "packed addresses" in the specification for more information.
	return uint32(vm.readWord(0x1A)) * 2
}

func (vm *VirtualMachine) StoryVersion() byte {
	return vm.Memory[0x00]
}

func (vm *VirtualMachine) highMemoryBase() uint16 {
	return vm.readWord(0x4)
}

func (vm *VirtualMachine) staticMemoryBase() uint16 {
	return vm.readWord(0xE)
}

func (vm *VirtualMachine) initialProgramCounter() uint16 {
	return vm.readWord(0x6)
}

func (vm *VirtualMachine) dictionaryLocation() uint16 {
	return vm.readWord(0x8)
}

func (vm *VirtualMachine) objectsLocation() uint16 {
	return vm.readWord(0xA)
}

func (vm *VirtualMachine) globalsLocation() uint16 {
	return vm.readWord(0xC)
}

func (vm *VirtualMachine) abbreviationsLocation() uint16 {
	return vm.readWord(0x18)
}

func (vm *VirtualMachine) PrintHeader() {
	if !vm.initialized {
		fmt.Println("VM not initialized")
		return
	}

	fmt.Printf("Story: v%d, %dKB (max address: 0x%08X)\n", vm.StoryVersion(), vm.StoryLength()/KILOBYTES, vm.StoryLength()-1)
	fmt.Printf("HighMem: 0x%04X\n", vm.highMemoryBase())
	fmt.Printf("StaticMem: 0x%04X\n", vm.staticMemoryBase())
	fmt.Printf("InitialPC: 0x%04X\n", vm.initialProgramCounter())
	fmt.Printf("Dictionary: 0x%04X\n", vm.dictionaryLocation())
	fmt.Printf("Object Table: 0x%04X\n", vm.objectsLocation())
	fmt.Printf("Globals: 0x%04X\n", vm.globalsLocation())
	fmt.Printf("Abbreviation Table: 0x%04X", vm.abbreviationsLocation())
}
