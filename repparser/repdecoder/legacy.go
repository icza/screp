/*

This file implements decoding the legacy (pre 1.18) replay format.
It partially implements reading PKWARE Data Compressed data,
tailored to the needs of parsing StarCraft: Brood War replay files (*.rep).


The algorithm comes from JCA's bwreplib.
Rewrite and optimization for Go: Andras Belicza


Information sources:

BWHF replay parser:
https://github.com/icza/bwhf/blob/master/src/hu/belicza/andras/bwhf/control/BinReplayUnpacker.java


Zadislav Zezula:
https://github.com/ladislav-zezula/StormLib/blob/master/src/pklib/explode.c


*/

package repdecoder

import "io"

var off507120 = []byte{ // length = 0x40
	0x02, 0x04, 0x04, 0x05, 0x05, 0x05, 0x05, 0x06,
	0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
	0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x07, 0x07,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
}

var off507160 = []byte{ // length = 0x40, com1
	0x03, 0x0D, 0x05, 0x19, 0x09, 0x11, 0x01, 0x3E,
	0x1E, 0x2E, 0x0E, 0x36, 0x16, 0x26, 0x06, 0x3A,
	0x1A, 0x2A, 0x0A, 0x32, 0x12, 0x22, 0x42, 0x02,
	0x7C, 0x3C, 0x5C, 0x1C, 0x6C, 0x2C, 0x4C, 0x0C,
	0x74, 0x34, 0x54, 0x14, 0x64, 0x24, 0x44, 0x04,
	0x78, 0x38, 0x58, 0x18, 0x68, 0x28, 0x48, 0x08,
	0xF0, 0x70, 0xB0, 0x30, 0xD0, 0x50, 0x90, 0x10,
	0xE0, 0x60, 0xA0, 0x20, 0xC0, 0x40, 0x80, 0x00,
}

var off5071A0 = []byte{ // length = 0x10
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
}

var off5071B0 = []byte{ // length = 0x20
	0x00, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00,
	0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00,
	0x08, 0x00, 0x0A, 0x00, 0x0E, 0x00, 0x16, 0x00,
	0x26, 0x00, 0x46, 0x00, 0x86, 0x00, 0x06, 0x01,
}

var off5071D0 = []byte{ // length = 0x10
	0x03, 0x02, 0x03, 0x03, 0x04, 0x04, 0x04, 0x05,
	0x05, 0x05, 0x05, 0x06, 0x06, 0x06, 0x07, 0x07,
}

var off5071E0 = []byte{ // length = 0x10, com1
	0x05, 0x03, 0x01, 0x06, 0x0A, 0x02, 0x0C, 0x14,
	0x04, 0x18, 0x08, 0x30, 0x10, 0x20, 0x40, 0x00,
}

// legacyDecoder is the Decoder implementation for legacy replays.
type legacyDecoder struct {
	decoder

	// esi struct used in decoding several sections
	esi esi
}

type replayEnc struct {
	src []byte
	m04 int32
	m08 []byte
	m0C int32
	m10 int32
	m14 int32
}

// zeroedEsiData is an array that remains untouched (zeroed) so the esi.data
// slice field can easily and efficiently be zeroed by copying this over
var zeroedEsiData [0x3114 + 0x20]byte // allocates 0x30 extra bytes in the beginning, but we ignore those

type esi struct {
	m00  int32
	m04  int32
	m08  int32
	m0C  int32
	m10  int32
	m14  int32
	m18  int32
	m1C  int32
	m20  int32
	m24  replayEnc
	m28  int32
	m2C  int32
	data []byte
}

func (d *legacyDecoder) Section(size int32) (result []byte, sectionID int32, err error) {
	var count int32
	if count, result, err = d.sectionHeader(size); result != nil || err != nil {
		return
	}

	var n, length2, m1C, m20, resultOffset int32
	result = make([]byte, size)

	d.initEsi()
	rep := &d.esi.m24

	buf := d.buf
	bufLen := int32(len(buf))
	for ; n < count; n, m1C, m20 = n+1, m1C+bufLen, m20+length2 {
		var length int32 // compressed length of the chunk
		if length, err = d.readInt32(); err != nil {
			return nil, sectionID, err
		}
		if length > size-m20 {
			return nil, sectionID, ErrMismatchedSection
		}

		if _, err = io.ReadFull(d.r, result[resultOffset:resultOffset+length]); err != nil {
			return
		}

		if length == min(size-m1C, bufLen) {
			continue
		}

		rep.src = make([]byte, length)
		copy(rep.src, result[resultOffset:])
		rep.m04 = 0
		rep.m08 = buf
		rep.m0C = 0
		rep.m10 = length
		rep.m14 = bufLen

		if d.repSection() == 0 && rep.m0C <= bufLen {
			length2 = rep.m0C
		} else {
			length2 = 0
		}
		if length2 == 0 || length2 > size {
			return nil, sectionID, ErrMismatchedSection
		}

		copy(result[resultOffset:], buf[:length2])
		resultOffset += length2
	}

	return result, sectionID, nil
}

// initEsi initializes (zeroes) the esi struct.
func (d *legacyDecoder) initEsi() {
	if d.esi.data == nil {
		// If this is the first call, we create and slice a new array:
		var data [len(zeroedEsiData)]byte // zeroed
		d.esi.data = data[:]
		// esi.m24 is a struct, its zero value is good.
	} else {
		// Else we copy over the zeroed slice:
		copy(d.esi.data, zeroedEsiData[:])
		// zero esi.m24 by assigning a new, zero-value struct
		d.esi.m24 = replayEnc{}
	}
}

// repSection decodes the esi.m24 (replayEnc) field.
func (d *legacyDecoder) repSection() int32 {
	esi := &d.esi

	esi.m1C = 0x800
	esi.m20 = d.esi28(0x2234, esi.m1C)
	if esi.m20 <= 4 {
		return 3
	}
	rep := &d.esi.m24
	esi.m04 = int32(rep.src[0])
	esi.m0C = int32(rep.src[1])
	esi.m14 = int32(rep.src[2])
	esi.m18 = 0
	esi.m1C = 3
	if esi.m0C < 4 || esi.m0C > 6 {
		return 1
	}
	esi.m10 = 1<<uint32(esi.m0C) - 1 // 2^n -1
	if esi.m04 != 0 {
		return 2
	}

	copy(esi.data[0x30F4:], off5071D0)
	d.com1(int32(len(off5071D0)), 0x30F4, off5071E0, 0x2B34)
	copy(esi.data[0x3104:], off5071A0)
	copy(esi.data[0x3114:], off5071B0)
	copy(esi.data[0x30B4:], off507120)
	d.com1(int32(len(off507160)), 0x30B4, off507160, 0x2A34)
	d.repChunk()

	return 0
}

func (d *legacyDecoder) com1(strlen, srcPos int32, str []byte, dstPos int32) {
	esi := &d.esi

	var x, y int32
	for n := strlen - 1; n >= 0; n-- {
		for x, y = int32(str[n]), 1<<esi.data[srcPos+n]; x < 0x100; x += y {
			esi.data[dstPos+x] = byte(n)
		}
	}
}

func (d *legacyDecoder) repChunk() int32 {
	esi := &d.esi

	esi.m08 = 0x1000
	var length int32
	for {
		length = d.function1()
		if length >= 0x305 {
			break
		}
		if length >= 0x100 { // decode region of size length -0xFE
			length -= 0xFE
			tmp := d.function2(length)
			if tmp == 0 {
				length = 0x306
				break
			}
			for length > 0 {
				esi.data[0x30+esi.m08] = esi.data[0x30+esi.m08-tmp]
				esi.m08++
				length--
			}
		} else {
			// just copy the character
			esi.data[0x30+esi.m08] = byte(length)
			esi.m08++
		}
		if esi.m08 < 0x2000 {
			continue
		}
		d.esi2C(0x1030, 0x1000)
		copy(esi.data[0x30:0x30+esi.m08-0x1000], esi.data[0x1030:])
		esi.m08 -= 0x1000
	}
	d.esi2C(0x1030, esi.m08-0x1000)

	return length
}

func (d *legacyDecoder) function1() int32 {
	esi := &d.esi

	var x, result int32

	// esi.m14 is odd
	if (1 & esi.m14) != 0 {
		if d.common(1) {
			return 0x306
		}
		result = int32(esi.data[0x2B34+(esi.m14&0xff)])
		if d.common(int32(esi.data[0x30F4+result])) {
			return 0x306
		}
		if esi.data[0x3104+result] != 0 {
			x = ((1 << (esi.data[0x3104+result] & 0xff)) - 1) & esi.m14
			if d.common(int32(esi.data[0x3104+result])) && (result+x) != 0x10E {
				return 0x306
			}
			result = (int32(esi.data[0x3114+2*result+1]) << 8) | int32(esi.data[0x3114+2*result]) // memcpy(&result, &myesi->m3114[2*result], 2);
			result += x
		}
		return result + 0x100
	}
	// esi.m14 is even
	if d.common(1) {
		return 0x306
	}
	if esi.m04 == 0 {
		result = esi.m14 & 0xff
		if d.common(8) {
			return 0x306
		}
		return result
	}
	if (esi.m14 & 0xff) == 0 {
		if d.common(8) {
			return 0x306
		}
		result = int32(esi.data[0x2EB4+(esi.m14&0xff)])
	} else {
		result = int32(esi.data[0x2C34+(esi.m14&0xff)])
		if result == 0xFF {
			if (esi.m14 & 0x3F) == 0 {
				if d.common(6) {
					return 0x306
				}
				result = int32(esi.data[0x2C34+(esi.m14&0x7F)])
			} else {
				if d.common(4) {
					return 0x306
				}
				result = int32(esi.data[0x2D34+(esi.m14&0xFF)])
			}
		}
	}
	if d.common(int32(esi.data[0x2FB4+result])) {
		return 0x306
	}
	return result
}

func (d *legacyDecoder) function2(length int32) int32 {
	esi := &d.esi

	tmp := int32(esi.data[0x2A34+esi.m14&0xff])
	if d.common(int32(esi.data[0x30B4+tmp])) {
		return 0
	}
	if length != 2 {
		tmp <<= byte(esi.m0C)
		tmp |= esi.m14 & esi.m10
		if d.common(esi.m0C) {
			return 0
		}
	} else {
		tmp <<= 2
		tmp |= esi.m14 & 3
		if d.common(2) {
			return 0
		}
	} // A38

	return tmp + 1
}

func (d *legacyDecoder) common(count int32) bool {
	esi := &d.esi

	if esi.m18 < count {
		esi.m14 >>= byte(esi.m18)
		if esi.m1C == esi.m20 {
			esi.m20 = d.esi28(0x2234, 0x800)
			if esi.m20 == 0 {
				return true
			}
			esi.m1C = 0
		}
		tmp := int32(esi.data[0x2234+esi.m1C])
		tmp <<= 8
		esi.m1C++
		tmp |= esi.m14
		esi.m14 = tmp
		tmp >>= uint32(count - esi.m18&0xff)
		esi.m14 = tmp
		esi.m18 += 8 - count
	} else {
		esi.m18 -= count
		esi.m14 >>= byte(count)
	}

	return false
}

func (d *legacyDecoder) esi28(dstPos, length int32) int32 {
	rep := &d.esi.m24

	length = min(rep.m10-rep.m04, length)
	copy(d.esi.data[dstPos:], rep.src[rep.m04:rep.m04+length])
	rep.m04 += length
	return length
}

func (d *legacyDecoder) esi2C(srcPos, length int32) {
	rep := &d.esi.m24

	if rep.m0C+length <= rep.m14 {
		copy(rep.m08[rep.m0C:], d.esi.data[srcPos:srcPos+length])
	}
	rep.m0C += length
}

// min returns the smaller of 2 int32 values.
func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
