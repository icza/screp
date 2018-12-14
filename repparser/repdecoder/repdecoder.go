/*

This file contains the interface to the replay decoder, and common parts
of the 2 types (legacy and modern).

Information sources:

BWHF replay parser:
https://github.com/icza/bwhf/blob/master/src/hu/belicza/andras/bwhf/control/BinReplayUnpacker.java

*/

package repdecoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	// ErrMismatchedSection is returned if the section size is not the expected one
	ErrMismatchedSection = errors.New("mismatched section")
)

// Decoder wraps a Section method for decoding a section of a given size.
type Decoder interface {
	// Section decodes a section of the given size.
	Section(size int32) (data []byte, err error)

	// Close closes the decoder, releases any associated resources.
	io.Closer
}

// NewFromFile creates a new Decoder that reads and decompresses data form a
// file.
func NewFromFile(name string) (d Decoder, err error) {
	var f *os.File
	f, err = os.Open(name)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	stat, err := f.Stat()
	if err != nil {
		return
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("not a file: %s", name)
	}

	var rf repFormat
	if stat.Size() >= 30 {
		fileHeader := make([]byte, 30)
		if _, err = io.ReadFull(f, fileHeader); err != nil {
			return
		}
		rf = detectRepFormat(fileHeader)
		if _, err = f.Seek(0, io.SeekStart); err != nil {
			return
		}
	}

	return newDecoder(f, rf), nil
}

// New creates a new Decoder that reads and decompresses data from the
// given byte slice.
func New(repData []byte) Decoder {
	rf := repFormatUnknown
	if len(repData) >= 30 {
		rf = detectRepFormat(repData[:30])
	}

	return newDecoder(bytes.NewBuffer(repData), rf)
}

// repFormat identifies the replay format
type repFormat int

// Possible values of repFormat
const (
	repFormatUnknown   repFormat = iota // Unknown replay format
	repFormatLegacy                     // Legacy replay format (pre 1.18)
	repFormatModern                     // Modern replay format (1.18 - 1.20)
	repFormatModern121                  // Modern 1.21 replay format (starting from 1.21)
)

// detectRepFormat detects the replay format based on the file header
// (the initial bytes of the binary replay).
// Information used from the header includes the replay ID section's data
// (which is 4 bytes, starting at offset 12), and the first bytes of the compressed
// data block of the Header section (which starts at offset 28).
// If the compressed data block starts with the magic of the valid zlib header,
// it is modern. If it is modern, the replay ID data decides which version.
func detectRepFormat(fileHeader []byte) repFormat {
	if len(fileHeader) < 30 {
		return repFormatUnknown
	}

	// Now only checking first byte of the compressed data block.
	// 2nd would be
	//     0x01 no compression
	//     0x5E level 1..5
	//     0x9C level 6 (default compression?)
	//     0xDA level 7..9
	if fileHeader[28] != 0x78 {
		return repFormatLegacy
	}

	// It is modern. Which one?
	// legacy and pre 1.21 modern replays have replay ID data "reRS".
	// Starting from 1.21, replay ID data is "seRS".
	if fileHeader[12] == 'r' {
		return repFormatModern
	}
	return repFormatModern121
}

// newDecoder creates a new Decoder that reads and decompresses data from the given Reader.
// The source is treated as a modern replay if modern is true, else as a
// legacy replay.
func newDecoder(r io.Reader, rf repFormat) Decoder {
	dec := decoder{
		r:        r,
		rf:       rf,
		int32Buf: make([]byte, 4),
		buf:      make([]byte, 0x2000), // 8 KB buffer
	}

	switch rf {
	case repFormatModern, repFormatModern121:
		return &modernDecoder{
			decoder: dec,
		}
	default:
		return &legacyDecoder{
			decoder: dec,
		}
	}
}

// decoder is the Decoder base (incomplete) implementation.
// Contains common parts of the 2 replay types.
type decoder struct {
	// r is the source of replay data
	r io.Reader

	// rf identifiers the rep format
	rf repFormat

	// intBuf is a general buffer for reading an int32 value
	int32Buf []byte

	// buf is a general buffer (re)used in decoding several sections
	buf []byte
}

// readInt32 reads an int32 from the underlying Reader.
func (d *decoder) readInt32() (n int32, err error) {
	if _, err = io.ReadFull(d.r, d.int32Buf); err != nil {
		return
	}

	n = int32(binary.LittleEndian.Uint32(d.int32Buf))
	return
}

// sectionHeader reads the section header.
func (d *decoder) sectionHeader(size int32) (count int32, result []byte, err error) {
	if size == 0 {
		result = []byte{}
		return
	}

	// checksum, we're not checking it
	if _, err = d.readInt32(); err != nil {
		return
	}

	// number of chunks the section data is split into
	count, err = d.readInt32()

	return
}

// Close closes the underlying io.Reader if it implements io.Closer.
func (d *decoder) Close() error {
	if closer, ok := d.r.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
