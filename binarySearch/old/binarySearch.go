package xed

import (
	"cmp"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/ddkwork/ddk/xed"
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

// D:\clone\systeminformer-master\SystemInformer
func DecodeTableByBinarySearchFunc() {
	file := xed.ParserPe("C:\\Windows\\System32\\ntoskrnl.exe")
	for _, section := range file.Sections {
		// 140443864确实是在text 140200000区段
		// 偏移大概是243864，之前的二分查找效率和准确度都堪忧，
		if section.String() == ".text" {
			data := section.Data(0, section.Header.SizeOfRawData, file)
			// ZwDeviceIoControlFile
			// KiServiceInternal
			// .text:0000000140443850                            KiSystemServiceStart:                   ; DATA XREF: KiServiceInternal+62↑o
			// .text:0000000140443850                                                                    ; .data:0000000140C02F00↓o
			// .text:0000000140443850 48 89 A3 90 00 00 00                       mov     [rbx+90h], rsp
			// .text:0000000140443857 8B F8                                      mov     edi, eax
			// .text:0000000140443859 C1 EF 07                                   shr     edi, 7
			// .text:000000014044385C 83 E7 20                                   and     edi, 20h
			// .text:000000014044385F 25 FF 0F 00 00                             and     eax, 0FFFh
			// .text:0000000140443864
			// .text:0000000140443864                            KiSystemServiceRepeat:                  ; CODE XREF: KiSystemCall64+C06↓j
			// .text:0000000140443864 4C 8D 15 55 E0 9B 00                       lea     r10, KeServiceDescriptorTable
			// .text:000000014044386B 4C 8D 1D 8E A8 8D 00                       lea     r11, KeServiceDescriptorTableShadow
			// .text:0000000140443872 F7 43 78 80 00 00 00                       test    dword ptr [rbx+78h], 80h
			// .text:0000000140443879 74 13                                      jz      short loc_14044388E
			// .text:000000014044387B F7 43 78 00 00 20 00                       test    dword ptr [rbx+78h], 200000h
			// .text:0000000140443882 74 07                                      jz      short loc_14044388B
			// .text:0000000140443884 4C 8D 1D F5 A9 8D 00                       lea     r11, KeServiceDescriptorTableFilter
			// .text:000000014044388B
			// .text:000000014044388B                            loc_14044388B:                          ; CODE XREF: KiSystemCall64+382↑j
			// .text:000000014044388B 4D 8B D3                                   mov     r10, r11

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
