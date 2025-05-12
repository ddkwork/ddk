package xed

import (
	"fmt"
	"iter"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"golang.org/x/arch/x86/x86asm"
)

func Is[T any](v any) iter.Seq[T] {
	return func(yield func(T) bool) {
		value, ok := v.(T)
		if ok {
			if !yield(value) {
				return
			}
		}
	}
}

// MovEaxImm  MOV  EAX,0x7   用于ntdll ntosker获取nt api编号，
// 存放于下面三张表的前两张表中，第一张表可以是nt和zw开头的api，
// 第二张表只有nt开头，个别api要排除，编号过大也不是，参考 DecodeNtApi 函数的过滤部分
// KeServiceDescriptorTable
// KeServiceDescriptorTableShadow
// KeServiceDescriptorTableFilter
func (o *object[T]) MovEaxImm(instruction x86asm.Inst) (imm int64) {
	if instruction.Op == x86asm.MOV {
		for i, arg := range instruction.Args {
			for reg := range Is[x86asm.Reg](arg) {
				if reg != x86asm.EAX {
					continue
				}
			}
			for mem := range Is[x86asm.Imm](instruction.Args[i+1]) {
				if mem > 10000 {
					return
				}
				return int64(mem)
			}
		}
	}
	return
}

type Disassembly struct {
	Address     uint64
	Opcode      []byte
	Instruction string // todo set color
	Comment     string
}
type (
	object[T stream.Type] struct {
		baseAddress    uint64
		data           []byte
		isFilterModel  bool // 反之调试器模式不应该在一个函数完成反汇编就结束
		Instructions   []x86asm.Inst
		IntelSyntaxAsm *stream.Buffer
		AsmObjects     []Disassembly
	}
	FilterInfo struct { // 用于搜索目标函数的虚拟地址，+内核基址=物理地址
		OpcodeData      []byte // 从导出表中取出的api虚拟地址取出100个字节大小的操作码，反汇编接口会自动识别ret指令停止
		BaseFunctionRVA uint32 // 从导出表中取出的api虚拟地址
		InstructionsLen int    // 找到目标函数停止时执行的指令长度
		DstFunctionRVA  uint64 // 相对于停止搜索时候的目标函数虚拟地址
		SysCallNumber   uint32 // 内核api编号，其实就是那两个表的切片下标，偏移，一个意思
		Ok              bool   // 是否找到目标函数，来源是观察符号或者调试确定位置之类的操作
	}
)

const (
	OpcodeDataSize = 500
)

func NewNewFilterInfo(opcodeData []byte, baseFunctionRVA uint32) *FilterInfo {
	return &FilterInfo{OpcodeData: opcodeData, BaseFunctionRVA: baseFunctionRVA}
}

func (f *FilterInfo) PrintDstFunctionRVA() { mylog.Hex("DstFunctionRVA", f.DstFunctionRVA) }
func (f *FilterInfo) PrintSysCallNumber() {
	mylog.Info("SysCallNumber", fmt.Sprintf("%#X | ", f.SysCallNumber)+fmt.Sprint(f.SysCallNumber))
}

func (f *FilterInfo) FunctionRVA() uint32 {
	return f.BaseFunctionRVA + uint32(f.InstructionsLen) + uint32(f.DstFunctionRVA)
}

func (o *object[T]) SetIsFilterModel(isFilterModel bool) { o.isFilterModel = isFilterModel }

func New[T stream.Type](data T) (x *object[T]) {
	b := stream.NewBuffer(data)
	// mylog.HexDump("data", b.Bytes())
	return &object[T]{
		baseAddress:    0xFFFFF8015FA00000,
		data:           b.Bytes(),
		isFilterModel:  true,
		Instructions:   make([]x86asm.Inst, 0),
		IntelSyntaxAsm: stream.NewBuffer(""),
	}
}

func (o *object[T]) SetBaseAddress(baseAddress uint64) *object[T] {
	o.baseAddress = baseAddress
	return o
}
func (o *object[T]) Decode32() (x *object[T]) { return o.decode(false) }
func (o *object[T]) Decode64() (x *object[T]) { return o.decode(true) }
func (o *object[T]) decode(is64Bit bool) *object[T] {
	mode := 32
	if is64Bit {
		mode = 64
	}
	o.AsmObjects = make([]Disassembly, 0)
	for len(o.data) > 0 {
		// mylog.Todo("inst need make object for set color") // todo D:\workspace\workspace\app\widget\CodeView.go return []label
		inst, e := x86asm.Decode(o.data, mode) // token 就不用解了，直接按nasm关键字着色返回label切片即可，然后在表格控件中填充每一行
		if e != nil {
			mylog.CheckIgnore(e)
			return o
		}
		o.Instructions = append(o.Instructions, inst)
		intelSyntax := x86asm.IntelSyntax(inst, 0, nil)
		o.IntelSyntaxAsm.WriteString(fmt.Sprintf("%#X", o.baseAddress))
		o.IntelSyntaxAsm.Indent(1)
		sprintf := fmt.Sprintf("%-20x\t", o.data[:inst.Len])
		o.IntelSyntaxAsm.WriteString(sprintf)
		o.IntelSyntaxAsm.WriteStringLn(intelSyntax)

		disassembly := Disassembly{
			Address:     o.baseAddress,
			Opcode:      o.data[:inst.Len],
			Instruction: intelSyntax, // todo set color
			Comment:     "",
		}
		o.AsmObjects = append(o.AsmObjects, disassembly)

		o.baseAddress += uint64(inst.Len)
		o.data = o.data[inst.Len:]
		isStop := false
		switch {
		case inst.Op == x86asm.RET:
			isStop = true
		case inst.Op == x86asm.INT: // cc int3 应该停止
			isStop = true
		}
		if isStop {
			// 如果是调试器调用，一个函数结束就停止反汇编是错误的操作
			// 然而对于定位偏移，取立即数，内存操作数等之类的情况，只需要在一个函数内完成
			if o.isFilterModel {
				break
			}
		}
	}
	return o
}

// 不知道为什么改名,基本上找不到的指令都可以翻这个因特尔文件得到
// C:\Users\Admin\go\pkg\mod\golang.org\x\arch@v0.6.0\x86\x86asm\intel.go
// var intelOp = map[Op]string{
//	JAE:       "jnb",
//	JA:        "jnbe",
//	JGE:       "jnl",
//	JNE:       "jnz",
//	JG:        "jnle",
//	JE:        "jz",

//	case INT:
//		if inst.Opcode>>24 == 0xCC {
//			args = nil
//			op = "int3"
//		}
