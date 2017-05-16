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

func (d *modernDecoder) Section(size int32) (result []byte, err error) {
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
			return nil, err
		}
		if length > 4 { // Is it compressed?
			if resetter, ok := zr.(zlib.Resetter); ok {
				err = resetter.Reset(bytes.NewBuffer(compressed), nil)
			} else {
				zr, err = zlib.NewReader(bytes.NewBuffer(compressed))
				defer zr.Close()
			}
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(resBuf, zr); err != nil {
				return nil, err
			}
		} else {
			// it's not compressed
			if _, err = resBuf.Write(compressed); err != nil {
				return
			}
		}
	}

	return resBuf.Bytes(), nil
}
