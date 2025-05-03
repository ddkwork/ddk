package ddk_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/ddkwork/ddk"
	"github.com/ddkwork/ddk/winver"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/saferwall/pe"
)

func TestRenameNtdll(t *testing.T) {
	t.Skip("once only")
	filepath.Walk("ntdll", func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			ok := path[:len(path)-len(filepath.Ext(path))] + "_windows.go"
			println(ok)
			mylog.Check(os.Rename(path, ok))
		}
		return err
	})
}

func TestPdb(t *testing.T) {
	t.Skip()
	pdb := "D:\\workspace\\workspace\\SysCall\\pdbfetch\\symbols\\ntkrnlmp.pdb\\AFA0F866CF448CC4D136836F5E5FAFBC1\\ntkrnlmp.pdb"
	file := mylog.Check2(pe.New(pdb, &pe.Options{
		Fast:                       true,
		SectionEntropy:             false,
		MaxCOFFSymbolsCount:        0,
		MaxRelocEntriesCount:       0,
		DisableCertValidation:      false,
		DisableSignatureValidation: false,
		Logger:                     nil,
		OmitExportDirectory:        false,
		OmitImportDirectory:        false,
		OmitExceptionDirectory:     false,
		OmitResourceDirectory:      false,
		OmitSecurityDirectory:      false,
		OmitRelocDirectory:         false,
		OmitDebugDirectory:         false,
		OmitArchitectureDirectory:  false,
		OmitGlobalPtrDirectory:     false,
		OmitTLSDirectory:           false,
		OmitLoadConfigDirectory:    false,
		OmitBoundImportDirectory:   false,
		OmitIATDirectory:           false,
		OmitDelayImportDirectory:   false,
		OmitCLRHeaderDirectory:     false,
	}))

	// if !mylog.Check(file.Parse()) {
	//	return
	// }
	mylog.Check(file.ParseCOFFSymbolTable())
}

// ntoskrnl.exe		0xFFFFF8015FA00000 todo
func Test_main(t *testing.T) {
	// RtlPcToFileHeader
	println(winver.WindowVersion())
	ddk.MiGetPteAddress()
	ddk.DecodeTableByDll()
	ddk.DecodeTableByDisassembly()
	// todo merge

	// D:\workspace\workspace\private\ui\model\branch\gui\plugin
	// call pdbfetch-master and set file.Debugs to use pdb
	// D:\workspace\hv\EWDK_quickstart
	// C:\Users\Admin\Downloads\cmakeconverter-develop

	// 驱动黑名单
	// https://githubfast.com/HotCakeX/Harden-Windows-Security/issues/125
}
