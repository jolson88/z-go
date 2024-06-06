package zmachine

import (
	"testing"
)

func TestZmCharToZscii(t *testing.T) {
	if zmCharToZscii(0, 0x6) != byte('a') {
		t.Errorf("Expected %d, got %d", byte('a'), zmCharToZscii(0, 6))
	}
	if zmCharToZscii(1, 0xE) != byte('I') {
		t.Errorf("Expected %d, got %d", byte('I'), zmCharToZscii(1, 0xE))
	}
	if zmCharToZscii(2, 0x7) != byte('\n') {
		t.Errorf("Expected %d, got %d", byte('\n'), zmCharToZscii(2, 0x7))
	}
	if zmCharToZscii(2, 0x1F) != byte(')') {
		t.Errorf("Expected %d, got %d", byte(')'), zmCharToZscii(2, 0x1F))
	}
}

func TestZmTextToString(t *testing.T) {
	result := zmTextToString([]*ZmText{
		{
			Chars:      [3]ZmChar{0x18, 0x6, 0xE},
			IsLastWord: false,
		},
		{
			Chars:      [3]ZmChar{0x11, 0x14, 0x17},
			IsLastWord: true,
		},
	})
	if result != "sailor" {
		t.Errorf("Expected 'sailor', got '%s'", result)
	}
}

func TestZmTextToStringShifting(t *testing.T) {
	result := zmTextToString([]*ZmText{
		{
			Chars:      [3]ZmChar{0x4, 0x6, 0x5},
			IsLastWord: false,
		},
		{
			Chars:      [3]ZmChar{0x11, 0x5, 0x5},
			IsLastWord: true,
		},
	})
	if result != "A9" {
		t.Errorf("Expected 'A9', got '%s'", result)
	}
}
