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

	var modern bool
	if stat.Size() >= 30 {
		if _, err = f.Seek(28, 0); err != nil {
			return
		}
		magic := make([]byte, 2)
		if _, err = io.ReadFull(f, magic); err != nil {
			return
		}
		modern = isModern(magic)
		if _, err = f.Seek(0, 0); err != nil {
			return
		}
	}

	return newDecoder(f, modern), nil
}

// New creates a new Decoder that reads and decompresses data from the
// given byte slice.
func New(repData []byte) Decoder {
	var modern bool
	if len(repData) >= 30 {
		modern = isModern(repData[28:30])
	}

	return newDecoder(bytes.NewBuffer(repData), modern)
}

// isModern expects the first bytes of the compressed data block of the Header
// section (which starts at offset 28), and tells if the replay is modern
// based on whether the magic is a valid zlib header.
// At least 2 bytes should be passed.
func isModern(magic []byte) bool {
	if len(magic) < 2 {
		return false
	}
	if magic[0] != 0x78 {
		return false
	}
	// Now only checking first byte.
	// 2nd would be
	//     0x01 no compression
	//     0x5E level 1..5
	//     0x9C level 6 (default compression?)
	//     0xDA level 7..9
	return true
}

// newDecoder creates a new Decoder that reads and decompresses data from the given Reader.
// The source is treated as a modern replay if modern is true, else as a
// legacy replay.
func newDecoder(r io.Reader, modern bool) Decoder {
	dec := decoder{
		r:        r,
		int32Buf: make([]byte, 4),
		buf:      make([]byte, 0x2000), // 8 KB buffer
	}

	switch modern {
	case true:
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
