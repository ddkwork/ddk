package xed

import (
	"cmp"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

// D:\clone\systeminformer-master\SystemInformer
func DecodeTableByBinarySearchFunc() {
	file := ParserPe("C:\\Windows\\System32\\ntoskrnl.exe")
	for _, section := range file.Sections {
		// 140443864确实是在text 140200000区段
		// 偏移大概是243864，之前的二分查找效率和准确度都堪忧，
		if section.String() == ".text" {
			data := section.Data(0, section.Header.SizeOfRawData, file)
			Pattern := "/x4C/x8D/x15/x00/x00/x00/x00/x4C/x8D/x1D/x00/x00/x00/x00/xF7"
			index, ok := BinarySearchFunc(data, Pattern)
			if !ok {
				// bytes := section.MetaData(uint32(index), xed.OpcodeDataSize, file)
				toString := hex.Dump(data)
				stream.WriteTruncate("ssdtData.txt", toString)
				// Zydis(data)
				// mylog.HexDump("", data)
				return
			}
			fmt.Printf("Label found at index: %d\n", index)

			// todo text section to ssdt
			// file.GetData(section.MetaData())
		}
	}
}

func BinarySearchFunc(data []byte, Pattern string) (index int, ok bool) {
	all := strings.ReplaceAll(Pattern, `/x`, "")
	decodeString := mylog.Check2(hex.DecodeString(all))
	return slices.BinarySearchFunc(data, decodeString, func(b byte, bytes []byte) int {
		for _, b2 := range bytes {
			if b2 == 0 {
				continue
			}
			return cmp.Compare(b, b2)
		}
		return 0
	})
}
