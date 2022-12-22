/*
yes...I know about encoding/binary...
*/

package util

import "fmt"

func Int8(src []byte, arch int) interface{} {
	if len(src) == 1 {
		return int(src[0])
	}

	// trivial case
	aa := []int{}
	for ii := 0; ii < len(src); ii += 2 {
		aa = append(aa, int(src[ii]))
	}
	return aa
}

func Int16(src []byte, arch int) interface{} {
	if len(src) == 2 {
		return Int(src[:], arch)
	}

	aa := []int{}
	for ii := 0; ii < len(src); ii += 2 {
		aa = append(aa, Int(src[ii:ii+2], arch))
	}
	return aa
}

func Int32(src []byte, arch int) interface{} {
	if len(src) == 4 {
		return Int(src[:], arch)
	}

	aa := []int{}
	for ii := 0; ii < len(src); ii += 4 {
		aa = append(aa, Int(src[ii:ii+4], arch))
	}
	return aa
}

func Int(src []byte, arch int) int {

	// fmt.Printf(">>>Int(%v, %d)\n", src, arch)

	const (
		ARCH_BIG_ENDIAN = 1
		ARCH_LITTLE_ENDIAN = 0
	)

	val := 0
	switch arch {
	case ARCH_BIG_ENDIAN:
		for ii := 0; ii < len(src); ii++ {
			val = val << 8 + int(src[ii])
		}
	case ARCH_LITTLE_ENDIAN:
		for ii := len(src) - 1; ii >= 0; ii-- {
			val = val << 8 + int(src[ii])
		}
	default:
		panic(fmt.Sprintf("bug: arch: %d\n", arch))
	}

	switch len(src) {
	case 1:
		if val < 0 || val > 0xFF {
			panic(fmt.Sprintf("bug: val: %d", val))
		}
	case 2:
		if val < 0 || val > 0xFFFF {
			panic(fmt.Sprintf("bug: val: %d", val))
		}
	case 4:
		if val < 0 || val > 0xFFFFFFFF {
			panic(fmt.Sprintf("bug: val: %d", val))
		}
	default:
		panic(fmt.Sprintf("bug: len: %d", len(src)))
	}
	return val
}

/*
Published by Garmin
FIT_UINT16 FitCRC_Get16(FIT_UINT16 crc, FIT_UINT8 byte)
{
   static const FIT_UINT16 crc_table[16] =
   {
      0x0000, 0xCC01, 0xD801, 0x1400, 0xF001, 0x3C00, 0x2800, 0xE401,
      0xA001, 0x6C00, 0x7800, 0xB401, 0x5000, 0x9C01, 0x8801, 0x4400
   };
   FIT_UINT16 tmp;

   // compute checksum of lower four bits of byte
   tmp = crc_table[crc & 0xF];
   crc = (crc >> 4) & 0x0FFF;
   crc = crc ^ tmp ^ crc_table[byte & 0xF];

   // now compute checksum of upper four bits of byte
   tmp = crc_table[crc & 0xF];
   crc = (crc >> 4) & 0x0FFF;
   crc = crc ^ tmp ^ crc_table[(byte >> 4) & 0xF];

   return crc;
}
*/

type FIT_UINT8 uint8
type FIT_UINT16 uint16

func FitCRC_Get16(crc FIT_UINT16, i8 FIT_UINT8) FIT_UINT16 {

   crc_table := [16]FIT_UINT16 {
      0x0000, 0xCC01, 0xD801, 0x1400, 0xF001, 0x3C00, 0x2800, 0xE401,
      0xA001, 0x6C00, 0x7800, 0xB401, 0x5000, 0x9C01, 0x8801, 0x4400,
   }

   tmp := FIT_UINT16(0)

   // compute checksum of lower four bits of byte
   tmp = crc_table[crc & 0xF]
   crc = (crc >> 4) & 0x0FFF
   crc = crc ^ tmp ^ crc_table[i8 & 0xF]

   // now compute checksum of upper four bits of byte
   tmp = crc_table[crc & 0xF]
   crc = (crc >> 4) & 0x0FFF
   crc = crc ^ tmp ^ crc_table[(i8 >> 4) & 0xF]

   return crc
}
