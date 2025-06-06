package ddk

import (
	"fmt"

	"github.com/ddkwork/ddk/xed"
	"github.com/ddkwork/golibrary/std/mylog"
	"golang.org/x/arch/x86/x86asm"
)

func DecodeTableByDisassembly() {
	file := xed.ParserPe("C:\\Windows\\System32\\ntoskrnl.exe")
	for _, entry := range file.Export.Functions {
		// mylog.Check(entry.Name)
		if entry.Name == "ZwDeviceIoControlFile" {
			data := mylog.Check2(file.GetData(entry.FunctionRVA, xed.OpcodeDataSize))

			info := xed.NewNewFilterInfo(data, entry.FunctionRVA)
			ZwDeviceIoControlFile(info)
			if !info.Ok {
				return
			}
			info.PrintSysCallNumber()
			info.PrintDstFunctionRVA()

			data = mylog.Check2(file.GetData(info.FunctionRVA(), xed.OpcodeDataSize))

			info = xed.NewNewFilterInfo(data, info.FunctionRVA())
			KiServiceInternal(info)
			if !info.Ok {
				return
			}
			data = mylog.Check2(file.GetData(info.FunctionRVA(), xed.OpcodeDataSize))
			info = xed.NewNewFilterInfo(data, info.FunctionRVA())
			// todo add win32k.sys and return ntapis as ssdt
			syscall, o := KiSystemServiceStart(info)
			if !o {
				return
			}
			mylog.Info("offsetKeServiceDescriptorTable", fmt.Sprintf("%#X | ", syscall.offsetKeServiceDescriptorTable)+fmt.Sprint(syscall.offsetKeServiceDescriptorTable))
			mylog.Info("offsetKeServiceDescriptorTableShadow", fmt.Sprintf("%#X | ", syscall.offsetKeServiceDescriptorTableShadow)+fmt.Sprint(syscall.offsetKeServiceDescriptorTableShadow))
			mylog.Info("offsetKeServiceDescriptorTableFilter", fmt.Sprintf("%#X | ", syscall.offsetKeServiceDescriptorTableFilter)+fmt.Sprint(syscall.offsetKeServiceDescriptorTableFilter))

			// data, err = file.GetData(syscall.offsetKeServiceDescriptorTable, xed.OpcodeDataSize)
			// if ! {
			//	return
			// }
			// mylog.HexDump("ssdt rva address", data) //难道参数表是动态填充的，静态解析全部是0
			//
			// data, err = file.GetData(syscall.offsetKeServiceDescriptorTableShadow, xed.OpcodeDataSize)
			// if ! {
			//	return
			// }
			// mylog.HexDump("shadow ssdt rva address", data)
			//
			// data, err = file.GetData(syscall.offsetKeServiceDescriptorTableFilter, xed.OpcodeDataSize)
			// if ! {
			//	return
			// }
			// mylog.HexDump("Filter rva address", data) //参数表,这里有算法 todo
			return
		}
	}
}

func ZwDeviceIoControlFile(info *xed.FilterInfo) {
	x := xed.New(info.OpcodeData).Decode64()
	for i, instruction := range x.Instructions {
		info.InstructionsLen += instruction.Len
		// b807000000              mov eax, 0x7
		// e902840xed.OpcodeDataSize              jmp .+0x18402  KiServiceInternal
		next := x.Instructions[i+1]
		if instruction.Op == x86asm.MOV && next.Op == x86asm.JMP {
			imm := x.MovEaxImm(instruction)
			for _, arg := range next.Args {
				for rel := range xed.Is[x86asm.Rel](arg) {
					info.InstructionsLen += next.Len
					info.DstFunctionRVA = uint64(rel)
					info.SysCallNumber = uint32(imm)
					info.Ok = true
					return
				}
			}
		}
	}
	return
}

func KiServiceInternal(info *xed.FilterInfo) {
	//       1404434e2 4c 8d 1d        LEA        R11,[KiSystemServiceStart]
	//                 67 03 00 00
	//       1404434e9 41 ff e3        JMP        R11=>KiSystemServiceStart                        undefined KiSystemServiceStart()
	x := xed.New(info.OpcodeData).Decode64()
	//	mylog.Json("", x.IntelSyntaxAsm.String())
	for i, instruction := range x.Instructions {
		info.InstructionsLen += instruction.Len
		next := x.Instructions[i+1]
		if instruction.Op == x86asm.LEA && next.Op == x86asm.JMP {
			for i, arg := range instruction.Args {
				// lea r11, ptr [rip+0x367]  KiSystemServiceStart
				//  jmp r11
				for reg := range xed.Is[x86asm.Reg](arg) {
					if reg == x86asm.R11 {
						a := instruction.Args[i+1]
						for mem := range xed.Is[x86asm.Mem](a) {
							info.DstFunctionRVA = uint64(mem.Disp)
							info.Ok = true
							return
						}
					}
				}
			}
		}
	}
	return
}

func KiSystemServiceStart(info *xed.FilterInfo) (syscall *SysCall, ok bool) {
	// todo 计算参数表
	//  *(BADSPACEBASE **)(unaff_RBX + 0x90) = register0x00000020;
	//  uVar10 = in_EAX >> 7 & 0x20;
	//  uVar13 = (ulonglong)(in_EAX & 0xfff);
	//  while( true ) {
	//    puVar12 = &KeServiceDescriptorTable;
	//    if (((*(uint *)(unaff_RBX + 0x78) & 0x80) != 0) &&
	//       (puVar12 = &KeServiceDescriptorTableShadow, (*(uint *)(unaff_RBX + 0x78) & 0x200000) != 0)) {
	//      puVar12 = (undefined8 *)&KeServiceDescriptorTableFilter;
	//    }

	//                             KiSystemServiceStart                            XREF[3]:     14016b453(*),
	//                                                                                          KiServiceInternal:1404434e2(*),
	//                                                                                          KiServiceInternal:1404434e9(c)
	//       140443850 48 89 a3        MOV        qword ptr [RBX + 0x90],RSP
	//                 90 00 00 00
	//       140443857 8b f8           MOV        EDI,EAX
	//       140443859 c1 ef 07        SHR        EDI,0x7
	//       14044385c 83 e7 20        AND        EDI,0x20
	//       14044385f 25 ff 0f        AND        EAX,0xfff
	//                 00 00
	//                             KiSystemServiceRepeat                           XREF[1]:     140444106(j)
	//       140443864 4c 8d 15        LEA        R10,[KeServiceDescriptorTable]
	//                 55 e0 9b 00
	//       14044386b 4c 8d 1d        LEA        R11,[KeServiceDescriptorTableShadow]             = ??
	//                 8e a8 8d 00
	//       140443872 f7 43 78        TEST       dword ptr [RBX + 0x78],0x80
	//                 80 00 00 00
	//       140443879 74 13           JZ         LAB_14044388e
	//       14044387b f7 43 78        TEST       dword ptr [RBX + 0x78],0x200000
	//                 00 00 20 00
	//       140443882 74 07           JZ         LAB_14044388b
	//       140443884 4c 8d 1d        LEA        R11,[KeServiceDescriptorTableFilter]             = ??
	//                 f5 a9 8d 00
	//                             LAB_14044388b                                   XREF[1]:     140443882(j)
	//       14044388b 4d 8b d3        MOV        R10,R11
	//                             LAB_14044388e                                   XREF[1]:     140443879(j)
	//       14044388e 41 3b 44        CMP        EAX,dword ptr [R10 + RDI*offset DAT_140d1e290    = ??
	//                 3a 10

	// 还是用反汇编引擎靠谱，然而，下载12mb的pdb太慢
	// 所以我们尝试无pdb提取三张水表，其实就三个全局未导出的数组

	// KeServiceDescriptorTable       ntdll.dll
	// KeServiceDescriptorTableShadow win32u.dll win32k.sys
	// KeServiceDescriptorTableFilter 参数个数表？
	// 如何计算hook目标函数物理地址？
	// kernelBase + 下标和参数移位？  todo 翻译cpp代码，在Ghidra里面验证

	// 为了承上启下的过滤特征，我们需要
	// 先保存解码过的指令对象到切片
	syscall = NewSysCall(0)
	x := xed.New(info.OpcodeData).Decode64()
	// mylog.Json("", x.IntelSyntaxAsm.String())
	size := 0
	for i, instruction := range x.Instructions {
		size += instruction.Len      // lea r10, ptr [rip+0x9be055]
		next := x.Instructions[i+1]  // lea r11, ptr [rip+0x8da88e]
		next1 := x.Instructions[i+2] // test dword ptr [rbx+0x78], 0x80
		next2 := x.Instructions[i+3] // jz .+0x13
		next3 := x.Instructions[i+4] // test dword ptr [rbx+0x78], 0x200000
		next4 := x.Instructions[i+5] // jz .+0x7
		next5 := x.Instructions[i+6] // lea r11, ptr [rip+0x8da9f5]

		// 不知道为什么改名
		// C:\Users\Admin\go\pkg\mod\golang.org\x\arch@v0.6.0\x86\x86asm\intel.go
		// var intelOp = map[Op]string{
		//	JAE:       "jnb",
		//	JA:        "jnbe",
		//	JGE:       "jnl",
		//	JNE:       "jnz",
		//	JG:        "jnle",
		//	JE:        "jz",
		if instruction.Op == x86asm.LEA &&
			next.Op == x86asm.LEA &&
			next1.Op == x86asm.TEST &&
			next2.Op == x86asm.JE &&
			next3.Op == x86asm.TEST &&
			next4.Op == x86asm.JE &&
			next5.Op == x86asm.LEA {
			for i2, arg := range instruction.Args {
				for reg := range xed.Is[x86asm.Reg](arg) {
					if reg == x86asm.R10 {
						for mem := range xed.Is[x86asm.Mem](instruction.Args[i2+1]) {
							syscall.offsetKeServiceDescriptorTable = info.BaseFunctionRVA + uint32(mem.Disp) + uint32(size)
						}
					}
				}
			}
			for i2, arg := range next.Args {
				for reg := range xed.Is[x86asm.Reg](arg) {
					if reg == x86asm.R11 {
						for mem := range xed.Is[x86asm.Mem](instruction.Args[i2+1]) {
							size += next.Len
							syscall.offsetKeServiceDescriptorTableShadow = info.BaseFunctionRVA + uint32(mem.Disp) + uint32(size)
						}
					}
				}
			}

			for i2, arg := range next5.Args {
				for reg := range xed.Is[x86asm.Reg](arg) {
					if reg == x86asm.R11 {
						for mem := range xed.Is[x86asm.Mem](instruction.Args[i2+1]) {
							size += next1.Len
							size += next2.Len
							size += next3.Len
							size += next4.Len
							size += next5.Len
							syscall.offsetKeServiceDescriptorTableFilter = info.BaseFunctionRVA + uint32(mem.Disp) + uint32(size)
							return syscall, true
						}
					}
				}
			}
		}
	}
	return
}

// typedef struct _SSDT
// {
//	LONG* ServiceTable;
//	PVOID CounterTable;
//	ULONG64 SyscallsNumber;
//	PVOID ArgumentTable;
// }_SSDT, *_PSSDT;
//
// _PSSDT NtTable;
// _PSSDT Win32kTable;

type (
	SysCall struct {
		kernelBase                           int64
		offsetKeServiceDescriptorTable       uint32  // 反汇编得到偏移，加上内核基址就是物理地址了
		offsetKeServiceDescriptorTableShadow uint32  // 反汇编得到偏移，加上内核基址就是物理地址了
		offsetKeServiceDescriptorTableFilter uint32  // 反汇编得到偏移，加上内核基址就是物理地址了
		KeServiceDescriptorTable             []NtApi // ghidra 显示内核基址+反汇编得到的偏移应该指向当前表的地址,
		// 当前表存放一个切片，每个成员指向各自的下标，参数值，签名
		KeServiceDescriptorTableShadow []NtApi
		// 总之就是通过参数表的对应值，下标，内核基址，可以计算一张
		// 对应的nt api物理地址表，hook的原理就是替换物理地址表里面的api物理地址
		// 实际上反汇编显示只存参数表，但是为了直观的展示逻辑
		// 我们把多张表合并为如下树形表格数据结构:
		//
		// kernelBase+offset=KeServiceDescriptorTable
		// KeServiceDescriptorTable容器节点{显示名称为基址+偏移},填充孩子节点普通表格
		// KeServiceDescriptorTableShadow,填充孩子节点
	}
	NtApi struct {
		KernelBase int64
		ArgValue   int32 // 计算物理地址，通过传入 kernelBase

		Name  string
		Index uint32
	}
)

func NewSysCall(kernelBase int64) *SysCall {
	return &SysCall{
		kernelBase:                     kernelBase,
		KeServiceDescriptorTable:       make([]NtApi, 0),
		KeServiceDescriptorTableShadow: make([]NtApi, 0),
	}
}

func (n *NtApi) PhysicalAddress() int64 { return n.KernelBase + int64(n.ArgValue) }
