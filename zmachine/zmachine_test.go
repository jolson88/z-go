package zmachine

import (
	"testing"
)

func TestZmCharConversions(t *testing.T) {
	if zmCharToZscii(0, 0x6) != byte('a') {
		t.Errorf("Expected %d, got %d", byte('a'), zmCharToZscii(0, 6))
	}
	if zmCharToZscii(2, 0x7) != byte('\n') {
		t.Errorf("Expected %d, got %d", byte('\n'), zmCharToZscii(2, 0x7))
	}
}
