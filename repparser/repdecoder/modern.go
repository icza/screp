/*

This file implements decoding the modern (starting from 1.18) replay format.

*/

package repdecoder

import (
	"bytes"
	"compress/zlib"
	"io"
)

// modernDecoder is the Decoder implementation for modern replays.
type modernDecoder struct {
	decoder
}

var knownModernSectionIDSizeHints = map[int32]int32{
	1313426259: 0x15e0, // "SKIN"
	1398033740: 0x1c,   // "LMTS"
	1481197122: 0x08,   // "BFIX"
	1380729667: 0xc0,   // "CCLR"
	1195787079: 0x19,   // "GCFG"
}

func (d *modernDecoder) Section(size int32) (result []byte, sectionID int32, err error) {
	if d.sectionsCounter > 5 {
		// These are the sections added in modern replays.
		if sectionID, err = d.readInt32(); err != nil { // This is the StrID of the section, not checking it
			return
		}
		var rawSize int32
		if rawSize, err = d.readInt32(); err != nil { // raw, remaining section size
			return
		}

		sizeHint := knownModernSectionIDSizeHints[sectionID]
		if sizeHint == 0 {
			// It's not a known, SCR section, but some custom section.
			// Don't assume anything about its format, return the raw data:
			result = make([]byte, rawSize)
			_, err = io.ReadFull(d.r, result)
			return
		}
		size = sizeHint
	}

	var count int32
	if count, result, err = d.sectionHeader(size); result != nil || err != nil {
		return
	}

	resBuf := bytes.NewBuffer(make([]byte, 0, size))

	var zr io.ReadCloser // zlib reader

	for ; count > 0; count-- {
		var length int32 // compressed length of the chunk
		if length, err = d.readInt32(); err != nil {
			return
		}

		if int32(len(d.buf)) < length {
			d.buf = make([]byte, length)
		}
		compressed := d.buf[:length]
		if _, err = io.ReadFull(d.r, compressed); err != nil {
			return nil, sectionID, err
		}
		if length > 4 && compressed[0] == 0x78 { // Is it compressed? (0x78 zlib magic)
			if resetter, ok := zr.(zlib.Resetter); ok {
				err = resetter.Reset(bytes.NewBuffer(compressed), nil)
			} else {
				zr, err = zlib.NewReader(bytes.NewBuffer(compressed))
				defer zr.Close()
			}
			if err != nil {
				return nil, sectionID, err
			}
			if _, err = io.Copy(resBuf, zr); err != nil {
				return nil, sectionID, err
			}
		} else {
			// it's not compressed
			if _, err = resBuf.Write(compressed); err != nil {
				return
			}
		}
	}

	return resBuf.Bytes(), sectionID, nil
}
