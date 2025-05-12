package ddk

import (
	"github.com/ddkwork/ddk/xed"
	"golang.org/x/arch/x86/x86asm"

	"github.com/ddkwork/golibrary/mylog"
)

// MiGetPteAddress 手工操作:
// Ghidra载入符号搜索MiGetPteAddress
// 查找交叉引用的函数，发现以下几个导出函数调用调用了它
// 没比对一个都要手动比对go解析的ntosken导出函数
// 那么如果写Ghidra脚本可以一劳永逸，流程如下:
// 1 定义一个切片，保存Ghidra解析的导出函数
// 2 输入目标函数执行交叉引用，遍历应用的函数，比较该函数是否在都出表切片中
// 3 人性化一点，检查交叉位置的汇编指令是否易于定义，最后输出最优的一个
// 4 执行非导出函数，宏的通过go反汇编引擎定位代码生成
// 如果Ghidra是go语言开发的话一定要实现这个逻辑，太方便了
// 注意尽量不要使用模糊暴力搜索，可读性和性能，排错都不容易
// 但是有一点，如果是创造调试器的话，是调sdk api实现接口方法方便还是类似od一样输入命令序列方便呢？有待观察
func MiGetPteAddress() {
	//       14022c2ec 48 c1 e9 09     SHR        RCX,0x9
	//       14022c2f0 48 b8 f8        MOV        RAX,0x7ffffffff8
	//                 ff ff ff
	//                 7f 00 00 00
	//       14022c2fa 48 23 c8        AND        RCX,RAX
	//       14022c2fd 48 b8 00        MOV        RAX,-0x98000000000
	//                 00 00 00
	//                 80 f6 ff ff
	//       14022c307 48 03 c1        ADD        RAX,RCX
	//       14022c30a c3              RET
	file := xed.ParserPe("C:\\Windows\\System32\\ntoskrnl.exe")
	for _, entry := range file.Export.Functions {
		// mylog.Check(entry.Name)
		if entry.Name == "MmFreeNonCachedMemory" {
			data := mylog.Check2(file.GetData(entry.FunctionRVA, xed.OpcodeDataSize))

			info := xed.NewNewFilterInfo(data, entry.FunctionRVA)
			MmFreeNonCachedMemory(info)
			if !info.Ok {
				mylog.Check("offset MiGetPteAddress not found")
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

func MmFreeNonCachedMemory(info *xed.FilterInfo) {
	x := xed.New(info.OpcodeData).Decode64()
	//       140a32a0b 48 83 ec 20     SUB        RSP,0x20
	//       140a32a0f 48 8b f2        MOV        RSI,RDX
	//       140a32a12 e8 d5 98        CALL       MiGetPteAddress                                  undefined MiGetPteAddress()
	//                 7f ff
	for i, instruction := range x.Instructions {
		info.InstructionsLen += instruction.Len
		nextInst := x.Instructions[i+1]
		nextInst2 := x.Instructions[i+2]
		if instruction.Op == x86asm.SUB &&
			nextInst.Op == x86asm.MOV &&
			nextInst2.Op == x86asm.CALL {
			for _, arg := range nextInst2.Args {
				for rel := range xed.Is[x86asm.Rel](arg) {
					info.InstructionsLen += nextInst.Len
					info.InstructionsLen += nextInst2.Len
					info.DstFunctionRVA = uint64(rel)
					info.Ok = true
					return
				}
			}
		}
	}
	return
}

//                             undefined MmFreeNonCachedMemory()
//       140a32a00 48 89 5c        MOV        qword ptr [RSP + local_res8],RBX
//                 24 08
//       140a32a05 48 89 74        MOV        qword ptr [RSP + local_res18],RSI
//                 24 18
//       140a32a0a 57              PUSH       RDI
//       140a32a0b 48 83 ec 20     SUB        RSP,0x20
//       140a32a0f 48 8b f2        MOV        RSI,RDX
//       140a32a12 e8 d5 98        CALL       MiGetPteAddress                                  undefined MiGetPteAddress()
//                 7f ff
//       140a32a17 48 8b c8        MOV        RCX,RAX
//       140a32a1a 48 8b f8        MOV        RDI,RAX
//       140a32a1d e8 de 4f        CALL       MI_READ_PTE_LOCK_FREE                            undefined MI_READ_PTE_LOCK_FREE()
//                 83 ff
//       140a32a22 48 8d 4c        LEA        RCX=>local_res10,[RSP + 0x38]
//                 24 38
//       140a32a27 48 89 44        MOV        qword ptr [RSP + local_res10],RAX
//                 24 38
//       140a32a2c e8 cf 4f        CALL       MI_READ_PTE_LOCK_FREE                            undefined MI_READ_PTE_LOCK_FREE()
//                 83 ff

//                             **************************************************************
//                             *                          FUNCTION                          *
//                             **************************************************************
//                             undefined MiGetPteAddress()

//       14022c2ec 48 c1 e9 09     SHR        RCX,0x9
//       14022c2f0 48 b8 f8        MOV        RAX,0x7ffffffff8
//                 ff ff ff
//                 7f 00 00 00
//       14022c2fa 48 23 c8        AND        RCX,RAX
//       14022c2fd 48 b8 00        MOV        RAX,-0x98000000000
//                 00 00 00
//                 80 f6 ff ff
//       14022c307 48 03 c1        ADD        RAX,RCX
//       14022c30a c3              RET

// "Location","Preview","Namespace","Mem Block"
// "1400d6d94","ibo32 MiGetPteAddress (_IMAGE_RUNTIME_FUNCTION_ENTRY.BeginAddress)","",".pdata"
// "14084a491","CALL MiGetPteAddress","FUN_14084a48e","PAGE"
// "14084a648","CALL MiGetPteAddress","FUN_14084a645","PAGE"
// "140a48e4c","CALL MiGetPteAddress","MiAllocateKernelCfgBitmapPageTables","PAGE"
// "140a48e57","CALL MiGetPteAddress","MiAllocateKernelCfgBitmapPageTables","PAGE"
// "140a39ec4","CALL MiGetPteAddress","MiApplyHotPatchToLoadedDriver","PAGE"
// "140a41f9e","CALL MiGetPteAddress","MiCopyPagesIntoEnclave","PAGE"
// "140a42b9a","CALL MiGetPteAddress","MiDecommitEnclavePages","PAGE"
// "14068bce5","CALL MiGetPteAddress","MiDecommitRegion","PAGE"
// "14068bd16","CALL MiGetPteAddress","MiDecommitRegion","PAGE"
// "140a47f13","CALL MiGetPteAddress","MiDeleteSparseRange","PAGE"
// "140a499a8","CALL MiGetPteAddress","MiExpandPartitionIds","PAGE"
// "1406f8509","CALL MiGetPteAddress","MiFindDriverNonPagedSections","PAGE"
// "1406fabce","CALL MiGetPteAddress","MiFreeInitializationCode","PAGE"
// "140a36f0f","CALL MiGetPteAddress","MiGetLargePagesForSystemMapping","PAGE"
// "140a36f28","CALL MiGetPteAddress","MiGetLargePagesForSystemMapping","PAGE"
// "1406fd0a8","CALL MiGetPteAddress","MiGetSystemAddressForImage","PAGE"
// "1406fd137","CALL MiGetPteAddress","MiGetSystemAddressForImage","PAGE"
// "140a3c4d3","CALL MiGetPteAddress","MiIdentifyImageDiscardablePages","PAGE"
// "14084a1f7","CALL MiGetPteAddress","MiInitializeDynamicBitmap","PAGE"
// "14084a367","CALL MiGetPteAddress","MiInitializeDynamicBitmap","PAGE"
// "14084bc16","CALL MiGetPteAddress","MiInitializeShadowPageTable","PAGE"
// "14084bd14","CALL MiGetPteAddress","MiInitializeShadowPageTable","PAGE"
// "140a31846","CALL MiGetPteAddress","MiMapNewPfns","PAGE"
// "140a31866","CALL MiGetPteAddress","MiMapNewPfns","PAGE"
// "140a3f079","CALL MiGetPteAddress","MiMapPatchTable","PAGE"
// "1406fc1ab","CALL MiGetPteAddress","MiMapSystemImage","PAGE"
// "1406fc25f","CALL MiGetPteAddress","MiMapSystemImage","PAGE"
// "1406fc26a","CALL MiGetPteAddress","MiMapSystemImage","PAGE"
// "1408d4a23","CALL MiGetPteAddress","MiMapSystemImage","PAGE"
// "140a4e947","CALL MiGetPteAddress","MiMapSystemImageWithLargePage","PAGE"
// "140a4eae1","CALL MiGetPteAddress","MiMapSystemImageWithLargePage","PAGE"
// "140695b0c","CALL MiGetPteAddress","MiPfPrepareReadList","PAGE"
// "140695b17","CALL MiGetPteAddress","MiPfPrepareReadList","PAGE"
// "14068c59d","CALL MiGetPteAddress","MiPfPrepareSequentialReadList","PAGE"
// "14068c5a8","CALL MiGetPteAddress","MiPfPrepareSequentialReadList","PAGE"
// "14084a81f","CALL MiGetPteAddress","MiProtectSystemImage","PAGE"
// "14084a865","CALL MiGetPteAddress","MiProtectSystemImage","PAGE"
// "14084a8b5","CALL MiGetPteAddress","MiProtectSystemImage","PAGE"
// "14084a9a7","CALL MiGetPteAddress","MiProtectSystemImage","PAGE"
// "140a46c5a","CALL MiGetPteAddress","MiReferenceIncomingPhysicalPages","PAGE"
// "140a40ab5","CALL MiGetPteAddress","MiReleaseHotPatchResources","PAGE"
// "1407ad7b7","CALL MiGetPteAddress","MiReserveDriverPtes","PAGE"
// "1407ce8f7","CALL MiGetPteAddress","MiResidentPagesForSpan","PAGE"
// "1407ce902","CALL MiGetPteAddress","MiResidentPagesForSpan","PAGE"
// "1407ce95e","CALL MiGetPteAddress","MiResidentPagesForSpan","PAGE"
// "1407ce969","CALL MiGetPteAddress","MiResidentPagesForSpan","PAGE"
// "1407ad490","CALL MiGetPteAddress","MiReturnSystemImageAddress","PAGE"
// "1406f9135","CALL MiGetPteAddress","MiSnapDriverRange","PAGE"
// "1406f91c5","CALL MiGetPteAddress","MiSnapDriverRange","PAGE"
// "1406f9215","CALL MiGetPteAddress","MiSnapDriverRange","PAGE"
// "140aad1dd","CALL MiGetPteAddress","MiTerminateHardwareEnclave","PAGELK"
// "140aad1f7","CALL MiGetPteAddress","MiTerminateHardwareEnclave","PAGELK"
// "1406fd27e","CALL MiGetPteAddress","MiUnloadSystemImage","PAGE"
// "14081593a","CALL MiGetPteAddress","MiUnlockDriverCode","PAGE"
// "140815946","CALL MiGetPteAddress","MiUnlockDriverCode","PAGE"
// "140935b37","CALL MiGetPteAddress","MiUnlockDriverPages","PAGE"
// "140a39184","CALL MiGetPteAddress","MiUnlockEntireDriver","PAGE"
// "140a4ebce","CALL MiGetPteAddress","MiUnmapLargeDriver","PAGE"
// "1407dabfe","CALL MiGetPteAddress","MiUnmapLockedPagesInUserSpace","PAGE"
// "14091a980","CALL MiGetPteAddress","MiUnmapLockedPagesInUserSpace","PAGE"
// "14084accf","CALL MiGetPteAddress","MmAllocateIsrStack","PAGE"
// "14077d082","CALL MiGetPteAddress","MmCommitSessionMappedView","PAGE"
// "14077d08d","CALL MiGetPteAddress","MmCommitSessionMappedView","PAGE"
// "14084b937","CALL MiGetPteAddress","MmCreateShadowMapping","PAGE"
// "14084b942","CALL MiGetPteAddress","MmCreateShadowMapping","PAGE"
// "140a41606","CALL MiGetPteAddress","MmDeleteShadowMapping","PAGE"
// "140a41615","CALL MiGetPteAddress","MmDeleteShadowMapping","PAGE"
// "14085ac40","CALL MiGetPteAddress","MmFreeBootRegistry","PAGE"
// "140a3926c","CALL MiGetPteAddress","MmFreeIndependentPages","PAGE"
// "140a41586","CALL MiGetPteAddress","MmFreeIsrStack","PAGE"
// "1407cde3e","CALL MiGetPteAddress","MmFreeMappingAddress","PAGE"
// "1407cdec7","CALL MiGetPteAddress","MmFreeMappingAddress","PAGE"
// "140a32a12","CALL MiGetPteAddress","MmFreeNonCachedMemory","PAGE"
// "140a355f7","CALL MiGetPteAddress","MmFreeSystemCacheReserveView","PAGE"
// "140a90bef","CALL MiGetPteAddress","MmInvalidateDumpAddresses","PAGELK"
// "1407eed27","CALL MiGetPteAddress","MmLockPreChargedPagedPool","PAGE"
// "140a490cb","CALL MiGetPteAddress","MmMapProtectedKernelPage","PAGE"
// "1407ef196","CALL MiGetPteAddress","MmMarkHiberRange","PAGE"
// "1407ef1a2","CALL MiGetPteAddress","MmMarkHiberRange","PAGE"
// "1406f8e47","CALL MiGetPteAddress","MmPageEntireDriver","PAGE"
// "1407edb1c","CALL MiGetPteAddress","MmReleaseDumpHibernateResources","PAGE"
// "1406f9341","CALL MiGetPteAddress","MmResetDriverPaging","PAGE"
// "1406f934d","CALL MiGetPteAddress","MmResetDriverPaging","PAGE"
// "140a356b5","CALL MiGetPteAddress","MmReturnChargesToLockPagedPool","PAGE"
// "1407f28bb","CALL MiGetPteAddress","MmStoreAllocateVirtualMemory","PAGE"
// "1407edac7","CALL MiGetPteAddress","MmUnlockPreChargedPagedPool","PAGE"
// "140a303fb","CALL MiGetPteAddress","MmUnmapLockedRestartPages","PAGE"
// "140a491aa","CALL MiGetPteAddress","MmUnmapProtectedKernelPageRange","PAGE"
// "140a477ed","CALL MiGetPteAddress","NtMapUserPhysicalPages","PAGE"
// "140a477fa","CALL MiGetPteAddress","NtMapUserPhysicalPages","PAGE"

func MmReturnChargesToLockPagedPool() {
}

//                             **************************************************************
//                             *                          FUNCTION                          *
//                             **************************************************************
//                             undefined MmReturnChargesToLockPagedPool()
//                               assume GS_OFFSET = 0xff00000000
//             undefined         AL:1           <RETURN>
//             undefined8        Stack[-0x18]:8 local_18                                XREF[1]:     140a3568f(W)
//             undefined1[16]    Stack[-0x28]   local_28                                XREF[1]:     140a356b0(W)
//             undefined1[16]    Stack[-0x38]   local_38                                XREF[1]:     140a356ab(W)
//             undefined1[16]    Stack[-0x48]   local_48                                XREF[2]:     140a356a6(W),
//                                                                                                   140a356cd(*)
//             undefined8        Stack[-0x50]:8 local_50                                XREF[1]:     140a356d4(W)
//             undefined4        Stack[-0x58]:4 local_58                                XREF[1]:     140a356d9(W)
//                             0xa35680  1584  MmReturnChargesToLockPagedPool
//                             Ordinal_1584                                    XREF[8]:     Entry Point(*), 14013af40(*),
//                             MmReturnChargesToLockPagedPool                               14014d8e4(*), 140170a02(*),
//                                                                                          MiDeleteSubsectionLargePages:140
//                                                                                          PopEnableHiberFile:14080ee5c(c),
//                                                                                          FUN_1408f601a:1408f6237(c),
//                                                                                          HalpUnloadMicrocode:14094bc64(c)
//       140a35680 40 53           PUSH       RBX
//       140a35682 48 83 ec 70     SUB        RSP,0x70
//       140a35686 33 c0           XOR        EAX,EAX
//       140a35688 48 8d 9a        LEA        RBX,[RDX + 0xfff]
//                 ff 0f 00 00
//       140a3568f 48 89 44        MOV        qword ptr [RSP + local_18],RAX
//                 24 60
//       140a35694 0f 57 c0        XORPS      XMM0,XMM0
//       140a35697 48 8b c1        MOV        RAX,RCX
//       140a3569a 25 ff 0f        AND        EAX,0xfff
//                 00 00
//       140a3569f 48 03 d8        ADD        RBX,RAX
//       140a356a2 48 c1 eb 0c     SHR        RBX,0xc
//       140a356a6 0f 11 44        MOVUPS     xmmword ptr [RSP + local_48[0]],XMM0
//                 24 30
//       140a356ab 0f 11 44        MOVUPS     xmmword ptr [RSP + local_38[0]],XMM0
//                 24 40
//       140a356b0 0f 11 44        MOVUPS     xmmword ptr [RSP + local_28[0]],XMM0
//                 24 50
//       140a356b5 e8 32 6c        CALL       MiGetPteAddress                                  undefined MiGetPteAddress()
//                 7f ff
//       140a356ba b9 02 00        MOV        ECX,0x2
//                 00 00
//       140a356bf 4c 8b c0        MOV        R8,RAX
//       140a356c2 e8 15 f6        CALL       MiGetAnyMultiplexedVm                            undefined MiGetAnyMultiplexedVm()
//                 7e ff
//       140a356c7 48 8b c8        MOV        RCX,RAX
//       140a356ca 4c 8b cb        MOV        R9,RBX
//       140a356cd 48 8d 44        LEA        RAX=>local_48,[RSP + 0x30]
//                 24 30
//       140a356d2 33 d2           XOR        EDX,EDX
//       140a356d4 48 89 44        MOV        qword ptr [RSP + local_50],RAX
//                 24 28
//       140a356d9 c7 44 24        MOV        dword ptr [RSP + local_58],0x8
//                 20 08 00
//                 00 00
//       140a356e1 e8 8a 2d        CALL       MiDeleteSystemPagableVm                          undefined MiDeleteSystemPagableV
//                 8b ff
//       140a356e6 48 8b d3        MOV        RDX,RBX
//       140a356e9 48 8d 0d        LEA        RCX,[MiSystemPartition]                          = ??
//                 d0 5e 23 00
//       140a356f0 e8 7f c8        CALL       MiReturnResident                                 undefined MiReturnResident()
//                 88 ff
//       140a356f5 48 83 c4 70     ADD        RSP,0x70
//       140a356f9 5b              POP        RBX
//       140a356fa c3              RET

func MmFreeMappingAddress() {
}

//                             **************************************************************
//                             *                          FUNCTION                          *
//                             **************************************************************
//                             undefined MmFreeMappingAddress()
//                               assume GS_OFFSET = 0xff00000000
//             undefined         AL:1           <RETURN>
//             undefined8        Stack[0x20]:8  local_res20                             XREF[2]:     1407cde0f(W),
//                                                                                                   1407cdeb3(R)
//             undefined8        Stack[0x18]:8  local_res18                             XREF[2]:     1407cde0b(W),
//                                                                                                   1407cdeae(R)
//             undefined8        Stack[0x10]:8  local_res10                             XREF[2]:     1407cde07(W),
//                                                                                                   1407cdea9(R)
//             undefined8        Stack[0x8]:8   local_res8                              XREF[2]:     1407cde03(W),
//                                                                                                   1407cdea4(R)
//             undefined8        Stack[-0x28]:8 local_28                                XREF[2]:     14091697c(RW),
//                                                                                                   14091699b(W)
//                             0x7cde00  1519  MmFreeMappingAddress
//                             Ordinal_1519                                    XREF[11]:    Entry Point(*), 140099b90(*),
//                             MmFreeMappingAddress                                         14011bf2c(*), 14014d7e0(*),
//                                                                                          14016eaea(*),
//                                                                                          SmFpCleanup:1404d1c56(c),
//                                                                                          PnprFreeMappingReserve:140977ac9
//                                                                                          PnprInitializeMappingReserve:140
//                                                                                          EtwpSavePersistedLogger:1409f2a3
//                                                                                          HalpDmaAllocateEmergencyResource
//                                                                                          HalpDmaAllocateMappingResources:
//       1407cde00 48 8b c4        MOV        RAX,RSP
//       1407cde03 48 89 58 08     MOV        qword ptr [RAX + local_res8],RBX
//       1407cde07 48 89 68 10     MOV        qword ptr [RAX + local_res10],RBP
//       1407cde0b 48 89 70 18     MOV        qword ptr [RAX + local_res18],RSI
//       1407cde0f 48 89 78 20     MOV        qword ptr [RAX + local_res20],RDI
//       1407cde13 41 54           PUSH       R12
//       1407cde15 41 56           PUSH       R14
//       1407cde17 41 57           PUSH       R15
//       1407cde19 48 83 ec 30     SUB        RSP,0x30
//       1407cde1d 8b fa           MOV        EDI,EDX
//       1407cde1f 48 8b d9        MOV        RBX,RCX
//       1407cde22 e8 a5 ba        CALL       MiRemoveMappingNode                              undefined MiRemoveMappingNode()
//                 b8 ff
//       1407cde27 4c 8b f8        MOV        R15,RAX
//       1407cde2a 39 78 28        CMP        dword ptr [RAX + 0x28],EDI
//       1407cde2d 0f 85 49        JNZ        LAB_14091697c
//                 8b 14 00
//       1407cde33 48 8b 68 18     MOV        RBP,qword ptr [RAX + 0x18]
//       1407cde37 48 8b 70 20     MOV        RSI,qword ptr [RAX + 0x20]
//       1407cde3b 48 8b cd        MOV        RCX,RBP
//       1407cde3e e8 a9 e4        CALL       MiGetPteAddress                                  undefined MiGetPteAddress()
//                 a5 ff
//       1407cde43 4c 8b f0        MOV        R14,RAX
//       1407cde46 48 8b d8        MOV        RBX,RAX
//       1407cde49 4c 8d 24 f0     LEA        R12,[RAX + RSI*0x8]
//       1407cde4d 49 3b c4        CMP        RAX,R12
//       1407cde50 73 28           JNC        LAB_1407cde7a
//                             LAB_1407cde52                                   XREF[1]:     1407cde78(j)
//       1407cde52 49 3b de        CMP        RBX,R14
//       1407cde55 74 6d           JZ         LAB_1407cdec4
//       1407cde57 48 f7 c3        TEST       RBX,0xfff
//                 ff 0f 00 00
//       1407cde5e 74 64           JZ         LAB_1407cdec4
