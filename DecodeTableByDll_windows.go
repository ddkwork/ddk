package ddk

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aquasecurity/table"
	"github.com/ddkwork/ddk/winver"
	"github.com/ddkwork/ddk/xed"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
)

func DecodeTableByDll() {
	sysCall := NewSysCall(0)
	sysCall.KeServiceDescriptorTable = DecodeNtApi("C:\\Windows\\System32\\ntdll.dll")
	sysCall.KeServiceDescriptorTableShadow = DecodeNtApi("C:\\Windows\\System32\\win32u.dll")
	// DecodeNtApi("C:\\Windows\\System32\\win32k.sys")
	// ntos c00082
	// D:\workspace\workspace\private\ui\model\branch\gui\plugin\symbol\AtomicSyscall-main
}

// D:\clone\systeminformer-master\SystemInformer\ksyscall.c
func DecodeNtApi(filename string) (ntApis []NtApi) {
	file := xed.ParserPe(filename)
	m := new(safemap.M[string, NtApi])
	for _, entry := range file.Export.Functions {
		if entry.Name == "NtGetTickCount" {
			continue
		}
		if strings.HasPrefix(entry.Name, "Ntdll") {
			continue
		}
		if strings.HasPrefix(entry.Name, "Nt") {
			data := mylog.Check2(file.GetData(entry.FunctionRVA, xed.OpcodeDataSize))
			x := xed.New(data).Decode64()
			index := uint32(0)
			for _, instruction := range x.Instructions {
				imm := x.MovEaxImm(instruction)
				index = uint32(imm)
				break
			}
			if strings.Contains(filename, "win32") {
				if index == 0 {
					continue
				}
			}
			if index >= 4096 { // win32kSyscall
				index -= 4096
			}
			call := NtApi{Name: entry.Name, Index: index}
			m.Set(call.Name, call)
		}
	}
	ntApis = make([]NtApi, m.Len())
	for i, s := range m.Keys() {
		get := m.GetMust(s)
		ntApis[i] = get
	}
	sort.Slice(ntApis, func(i, j int) bool {
		return ntApis[i].Index < ntApis[j].Index
	})

	if m.Len() == 0 {
		return
	}

	s := stream.NewBuffer("")
	s.WriteStringLn(winver.WindowVersion())
	t := table.New(s)
	t.SetRowLines(false)
	t.SetHeaders("Id", "Name", "Index")

	for i, call := range ntApis {
		t.AddRow(fmt.Sprint(i+1),
			call.Name,
			stream.FormatInteger(call.Index),
		)
	}
	time.Sleep(time.Second)
	t.Render()

	stream.WriteTruncate(stream.BaseName(filename)+".txt", s.String())
	return
}
