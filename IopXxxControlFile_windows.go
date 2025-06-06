package ddk

import (
	"github.com/ddkwork/ddk/xed"
	"github.com/ddkwork/golibrary/std/mylog"
	"golang.org/x/arch/x86/x86asm"
)

func NtDeviceIoControlFile() {
	file := xed.ParserPe("C:\\Windows\\System32\\ntoskrnl.exe")
	for _, entry := range file.Export.Functions {
		// mylog.Check(entry.Name)
		if entry.Name == "NtDeviceIoControlFile" {
			data := mylog.Check2(file.GetData(entry.FunctionRVA, xed.OpcodeDataSize))

			info := xed.NewNewFilterInfo(data, entry.FunctionRVA)
			IopXxxControlFile(info)
			if !info.Ok {
				mylog.Check("offset IopXxxControlFile not found")
				return
			}
			info.PrintDstFunctionRVA()
			// 计算公式
			// 当前区段内的导出函数rva + 函数头到定位停止那条指令的rva + 定位停止那条指令取出的目标函数rva = 目标函数相对于内核基址的rva
			// 注意多条指令联合定位的时候需要合计停止遍历前的所有指令长度，包括停止的那条指令长度，因为它参与目标函数rva的计算，要取正确才行。
			// 所以这个rva传递给驱动，加上驱动取出内核基址就是目标函数的物理地址了
			// 最后强制目标函数类型指针转换这个物理地址就可以全局调用该函数了
			// 注意区分有参数的情况，比如iopxxxcontrolfile之列的
			data = mylog.Check2(file.GetData(info.FunctionRVA(), xed.OpcodeDataSize))

			println(xed.New(data).Decode64().IntelSyntaxAsm.String())
			return
		}
	}
}

func IopXxxControlFile(info *xed.FilterInfo) {
	x := xed.New(info.OpcodeData).Decode64()
	for i, instruction := range x.Instructions {
		info.InstructionsLen += instruction.Len
		nextInst := x.Instructions[i+1]
		nextInst2 := x.Instructions[i+2]
		if instruction.Op == x86asm.CALL &&
			nextInst.Op == x86asm.ADD &&
			nextInst2.Op == x86asm.RET {
			for _, arg := range instruction.Args {
				for rel := range xed.Is[x86asm.Rel](arg) {
					info.DstFunctionRVA = uint64(rel)
					info.Ok = true
					return
				}
			}
		}
	}
	return
}

//                             NtDeviceIoControlFile                                        14014d9c4(*), 14016dda7(*),
//                                                                                          140aa09f2(c)
//       1406b5cb0 48 83 ec 68     SUB        RSP,0x68
//       1406b5cb4 8b 84 24        MOV        EAX,dword ptr [RSP + param_10]
//                 b8 00 00 00
//       1406b5cbb c6 44 24        MOV        byte ptr [RSP + local_18],0x1
//                 50 01
//       1406b5cc0 89 44 24 48     MOV        dword ptr [RSP + local_20],EAX
//       1406b5cc4 48 8b 84        MOV        RAX,qword ptr [RSP + param_9]
//                 24 b0 00
//                 00 00
//       1406b5ccc 48 89 44        MOV        qword ptr [RSP + local_28],RAX
//                 24 40
//       1406b5cd1 8b 84 24        MOV        EAX,dword ptr [RSP + param_8]
//                 a8 00 00 00
//       1406b5cd8 89 44 24 38     MOV        dword ptr [RSP + local_30],EAX
//       1406b5cdc 48 8b 84        MOV        RAX,qword ptr [RSP + param_7]
//                 24 a0 00
//                 00 00
//       1406b5ce4 48 89 44        MOV        qword ptr [RSP + local_38],RAX
//                 24 30
//       1406b5ce9 8b 84 24        MOV        EAX,dword ptr [RSP + param_6]
//                 98 00 00 00
//       1406b5cf0 89 44 24 28     MOV        dword ptr [RSP + local_40],EAX
//       1406b5cf4 48 8b 84        MOV        RAX,qword ptr [RSP + param_5]
//                 24 90 00
//                 00 00
//       1406b5cfc 48 89 44        MOV        qword ptr [RSP + local_48],RAX
//                 24 20
//       1406b5d01 e8 8a ea        CALL       IopXxxControlFile                                undefined IopXxxControlFile(unde
//                 ff ff
//       1406b5d06 48 83 c4 68     ADD        RSP,0x68
//       1406b5d0a c3              RET
//       1406b5d0b cc              ??         CCh
//                             LAB_1406b5d0c                                   XREF[1]:     1401120fc(*)
//       1406b5d0c cc              INT3
//       1406b5d0d cc              ??         CCh
//       1406b5d0e cc              ??         CCh
//       1406b5d0f cc              ??         CCh
//       1406b5d10 cc              ??         CCh

//                             undefined IopXxxControlFile(undefined param_1, undefined
//       1406b4790 4c 89 4c        MOV        qword ptr [RSP + local_res20],param_4
//                 24 20
//       1406b4795 4c 89 44        MOV        qword ptr [RSP + local_res18],param_3
//                 24 18
//       1406b479a 53              PUSH       RBX
//       1406b479b 56              PUSH       RSI
//       1406b479c 57              PUSH       RDI
//       1406b479d 41 54           PUSH       R12
//       1406b479f 41 55           PUSH       R13
//       1406b47a1 41 56           PUSH       R14
//       1406b47a3 41 57           PUSH       R15
//       1406b47a5 48 81 ec        SUB        RSP,0x1c0
//                 c0 01 00 00
//       1406b47ac 48 8b 05        MOV        RAX,qword ptr [__security_cookie]                = 00002B992DDFA232h
//                 bd 86 55 00
//       1406b47b3 48 33 c4        XOR        RAX,RSP
//       1406b47b6 48 89 84        MOV        qword ptr [RSP + local_48],RAX
//                 24 b0 01
//                 00 00
