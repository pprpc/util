package common

import "fmt"

// ByteConvertString byte conver to string
func ByteConvertString(in []byte) (str string) {
	for _, v := range in {
		str = str + fmt.Sprintf("%02x", v)
	}
	return
}

// DebugByte conver byte to string
func DebugByte(in []byte) (str string) {
	str = str + fmt.Sprintf("      01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16\n")
	for i, v := range in {
		if i == 0 {
			str = str + fmt.Sprintf("0000:")
		}
		str = str + fmt.Sprintf(" %02x", v)
		if i != 0 && (i%16) == 15 {
			str = str + fmt.Sprintf("\n")
			str = str + fmt.Sprintf("%04x:", i)
		}
	}
	return str
}
