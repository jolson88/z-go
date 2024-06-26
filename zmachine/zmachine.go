package zmachine

import (
	"fmt"
	"log"
	"os"
)

const KILOBYTES = 1024

type ZmWord uint16
type ZmChar byte
type ZmText struct {
	IsLastWord bool
	Chars      [3]ZmChar
}

type ZmVm struct {
	initialized bool
	Memory      [512 * KILOBYTES]byte // The maximum size of any z-machine file

	// Dictionary
	dictEntryLength   byte
	dictEntryCount    ZmWord
	dictEntryStart    uint16
	encodedTextLength byte
	wordSeparators    []byte
}

func NewVirtualMachine(pathToStoryFile string) *ZmVm {
	vm := &ZmVm{}

	data, err := os.ReadFile(pathToStoryFile)
	if err != nil {
		log.Fatalf("Error reading story file: %v", err)
	}

	copy(vm.Memory[:], data)

	dictionaryAddress := vm.dictionaryLocation()
	wordSeparatorCount := vm.Memory[dictionaryAddress]
	for i := 0; i < int(wordSeparatorCount); i++ {
		vm.wordSeparators = append(vm.wordSeparators, vm.Memory[dictionaryAddress+ZmWord(i)+1])
	}
	entryLengthLocation := dictionaryAddress + ZmWord(wordSeparatorCount) + 1
	vm.dictEntryLength = vm.Memory[entryLengthLocation]
	vm.dictEntryCount = vm.readWord(uint16(entryLengthLocation) + 1)
	vm.dictEntryStart = uint16(entryLengthLocation) + 3
	vm.encodedTextLength = 4

	vm.initialized = true
	return vm
}

/*
Memory Read/Write
*/
func (vm *ZmVm) PrintMemory(startAddress uint16, lines uint16) {
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

func (vm *ZmVm) readWord(address uint16) ZmWord {
	// Z-Machine is Big Endian. Need to swap around the bytes on Little Endian systems.
	return ZmWord(uint16(vm.Memory[address])<<8 | uint16(vm.Memory[address+1]))
}

/*
Header Information
*/
func (vm *ZmVm) StoryChecksum() ZmWord {
	return vm.readWord(0x1C)
}

func (vm *ZmVm) StoryLength() uint32 {
	// Up to v3, the story length this value multiplied by 2.
	// See "packed addresses" in the specification for more information.
	return uint32(vm.readWord(0x1A)) * 2
}

func (vm *ZmVm) StoryVersion() byte {
	return vm.Memory[0x00]
}

func (vm *ZmVm) highMemoryBase() ZmWord {
	return vm.readWord(0x4)
}

func (vm *ZmVm) staticMemoryBase() ZmWord {
	return vm.readWord(0xE)
}

func (vm *ZmVm) initialProgramCounter() ZmWord {
	return vm.readWord(0x6)
}

func (vm *ZmVm) dictionaryLocation() ZmWord {
	return vm.readWord(0x8)
}

func (vm *ZmVm) objectsLocation() ZmWord {
	return vm.readWord(0xA)
}

func (vm *ZmVm) globalsLocation() ZmWord {
	return vm.readWord(0xC)
}

func (vm *ZmVm) abbreviationsLocation() ZmWord {
	return vm.readWord(0x18)
}

func (vm *ZmVm) PrintHeader() {
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
	fmt.Printf("Abbreviation Table: 0x%04X\n", vm.abbreviationsLocation())
}

/*
Text
*/
var defaultAlphabets = [3]string{
	"abcdefghijklmnopqrstuvwxyz",    // A0
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",    // A1
	" \n0123456789.,!?_#'\"/\\-:()", // A2
}

func parseText(word uint16) *ZmText {
	return &ZmText{
		IsLastWord: word&0x8000 != 0,
		Chars: [3]ZmChar{
			ZmChar(word >> 10 & 0x1F),
			ZmChar(word >> 5 & 0x1F),
			ZmChar(word & 0x1F),
		},
	}
}

func zmCharToZscii(alphabet byte, char ZmChar) byte {
	if char < 6 || char > 31 || alphabet >= 3 {
		return 0
	}

	return defaultAlphabets[alphabet][char-6]
}

func zmTextToString(text []*ZmText) string {
	result := ""
	var nextAlphabet byte = 0
	for _, word := range text {
		for i := 0; i < 3; i++ {
			char := word.Chars[i]
			if char == 0x4 {
				nextAlphabet = 1
				continue
			}
			if char == 0x5 {
				nextAlphabet = 2
				continue
			}

			result += string(zmCharToZscii(nextAlphabet, word.Chars[i]))
			nextAlphabet = 0
		}
		if word.IsLastWord {
			break
		}
	}
	return result
}

/*
Dictionary
*/
type DictionaryEntry struct {
	Text []ZmWord
	Data []byte
}

func (vm *ZmVm) PrintDictionary() {
	if !vm.initialized {
		fmt.Println("VM not initialized")
		return
	}

	fmt.Print("Word Separators: ")
	for _, separator := range vm.wordSeparators {
		fmt.Printf("%c ", separator)
	}
	fmt.Printf("\nEntry Length: %d\n", vm.dictEntryLength)
	fmt.Printf("Entry Count: %d\n", vm.dictEntryCount)
}

func (vm *ZmVm) PrintDictionaryEntry(entryIndex uint16) {
	if !vm.initialized {
		fmt.Println("VM not initialized")
		return
	}
	if entryIndex >= uint16(vm.dictEntryCount) {
		fmt.Println("Invalid dictionary entry index:", entryIndex)
		return
	}

	dictEntry := vm.dictionaryEntry(entryIndex)

	texts := []*ZmText{}
	for _, text := range dictEntry.Text {
		texts = append(texts, parseText(uint16(text)))
	}
	fmt.Printf("Text: %s\n", zmTextToString(texts))
	fmt.Printf("Data: %v\n", dictEntry.Data)
}

func (vm *ZmVm) dictionaryEntry(idx uint16) *DictionaryEntry {
	dictEntryLocation := vm.dictEntryStart + uint16(idx)*uint16(vm.dictEntryLength)
	dictEntry := vm.Memory[dictEntryLocation : dictEntryLocation+uint16(vm.dictEntryLength)]

	textBytes := dictEntry[0:vm.encodedTextLength]
	entryText := []ZmWord{}
	for i := 0; i < len(textBytes); i += 2 {
		word := uint16(textBytes[i])<<8 | uint16(textBytes[i+1])
		entryText = append(entryText, ZmWord(word))
	}

	return &DictionaryEntry{
		Text: entryText,
		Data: dictEntry[vm.encodedTextLength:],
	}
}
