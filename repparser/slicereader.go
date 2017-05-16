// This file contains a slice reader which aids reading data from a byte slice.

package repparser

import "encoding/binary"

// sliceReader aids reading data from a byte slice
type sliceReader struct {
	// b is the byte slice to read from
	b []byte

	// pos is the index of the next byte to read
	pos uint32
}

// getByte returns the next byte.
func (sr *sliceReader) getByte() (r byte) {
	r, sr.pos = sr.b[sr.pos], sr.pos+1
	return
}

// getUint16 returns the next 2 bytes as an uint16 value.
func (sr *sliceReader) getUint16() (r uint16) {
	r, sr.pos = binary.LittleEndian.Uint16(sr.b[sr.pos:]), sr.pos+2
	return
}

// getUint32 returns the next 4 bytes as an uint32 value.
func (sr *sliceReader) getUint32() (r uint32) {
	r, sr.pos = binary.LittleEndian.Uint32(sr.b[sr.pos:]), sr.pos+4
	return
}

// getString returns the next size bytes as a string.
func (sr *sliceReader) getString(size uint32) (r string) {
	r, sr.pos = string(sr.b[sr.pos:sr.pos+size]), sr.pos+size
	return
}

// readSlice returns a the next size bytes as a slice.
func (sr *sliceReader) readSlice(size uint32) (r []byte) {
	r = make([]byte, size)
	sr.pos += uint32(copy(r, sr.b[sr.pos:]))
	return
}
