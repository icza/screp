package repparser

import (
	"testing"
)

// TestParseRawTrailingChunk is a regression test for a modern replay whose
// command section ends with an uncompressed chunk that happens to start with
// the byte 0x78. The decoder used to treat any chunk starting with 0x78 as
// zlib-compressed, so it fed the raw bytes to zlib and failed with
// "zlib: invalid header".
func TestParseRawTrailingChunk(t *testing.T) {
	const name = "testdata/shieldbattery_raw_trailing_0x78.rep"

	r, err := ParseFile(name)
	if err != nil {
		t.Fatalf("ParseFile(%q) error: %v", name, err)
	}
	if r.Commands == nil || len(r.Commands.Cmds) == 0 {
		t.Fatalf("expected parsed commands, got none")
	}
}
