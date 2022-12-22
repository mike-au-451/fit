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

