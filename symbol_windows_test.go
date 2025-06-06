// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package ddk

import (
	"errors"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/saferwall/pe"

	"github.com/ddkwork/golibrary/std/mylog"
)

type TestCOFFSymbol struct {
	errTooManySymbols error
	symbolsCount      int
	symbolIdx         int
	symbol            pe.COFFSymbol
	stringTableOffset uint32
	symbolName        string
	sectionNumberName string
	symbolTypeString  string
}

func getAbsoluteFilePath(testfile string) string {
	_, p, _ := mylog.Check4Bool(runtime.Caller(0))
	join := path.Join(filepath.Dir(p), testfile)
	return strings.ReplaceAll(join, "/", "\\")
}

var symbolTests = []struct {
	in  string
	out TestCOFFSymbol
}{
	{
		getAbsoluteFilePath("test/liblzo2-2.dll"),
		TestCOFFSymbol{
			errTooManySymbols: nil,
			symbolsCount:      50,
			symbolIdx:         0,
			symbol: pe.COFFSymbol{
				Name:               [8]byte{0, 0, 0, 0, 4, 0, 0, 0},
				Value:              0x2ac,
				SectionNumber:      8,
				Type:               0x0,
				StorageClass:       0x2,
				NumberOfAuxSymbols: 0x0,
			},
			stringTableOffset: 0x35184,
			symbolName:        "__imp_abort",
			sectionNumberName: ".idata",
			symbolTypeString:  "Null",
		},
	},

	{
		getAbsoluteFilePath(
			"test/0103daa751660333b7ae5f098795df58f07e3031563e042d2eb415bffa71fe7a.7z",
		),
		TestCOFFSymbol{
			errTooManySymbols: nil,
			symbolsCount:      346,
			symbolIdx:         3,
			symbol: pe.COFFSymbol{
				Name:               [8]byte{0, 0, 0, 0, 4, 0, 0, 0},
				Value:              0x2ac,
				SectionNumber:      8,
				Type:               0x0,
				StorageClass:       0x2,
				NumberOfAuxSymbols: 0x0,
			},
			stringTableOffset: 0x1b054,
			symbolName:        "___mingw_CRTStartup",
			sectionNumberName: ".text",
			symbolTypeString:  "",
		},
	},

	{
		getAbsoluteFilePath(
			"test/0000e876c5b712b6b7b3ce97f757ddd918fb3dbdc5a3938e850716fbd841309f.7z",
		),
		TestCOFFSymbol{
			errTooManySymbols: errCOFFSymbolsTooHigh,
		},
	},
}

var (
	errCOFFTableNotPresent = errors.New(
		"PE image does not contains a COFF symbol table")
	errNoCOFFStringInTable = errors.New(
		"PE image got a PointerToSymbolTable but no string in the COFF string table")
	errCOFFSymbolOutOfBounds = errors.New(
		"COFF symbol offset out of bounds")
	errCOFFSymbolsTooHigh = errors.New(
		"COFF symbols count is absurdly high")
)

func TestParseCOFFSymbolTable(t *testing.T) {
	t.Skip("todo")
	for _, tt := range symbolTests {
		mylog.Check(tt.in)
		t.Run(tt.in, func(t *testing.T) {
			ops := pe.Options{Fast: true}
			file := mylog.Check2(pe.New(tt.in, &ops))
			mylog.Check(file.Parse())
			mylog.Check(file.ParseCOFFSymbolTable())
			// exit early when err is errCOFFSymbolsTooHigh.
			if len(file.COFF.SymbolTable) != tt.out.symbolsCount {
				t.Errorf(
					"symbolsCount assertion failed, want: %d, got: %d",
					tt.out.symbolsCount,
					len(file.COFF.SymbolTable),
				)
			}
			if file.COFF.StringTableOffset != tt.out.stringTableOffset {
				t.Errorf(
					"stringTableOffset assertion failed, want: %d, got: %d",
					tt.out.stringTableOffset,
					file.COFF.StringTableOffset,
				)
			}
			if !stringInSlice(tt.out.symbolName, file.COFF.StringTable) {
				t.Errorf(
					"symbolName assertion failed, want: %s, got: %v",
					tt.out.symbolName,
					file.COFF.StringTable,
				)
			}

			coffSymbol := file.COFF.SymbolTable[tt.out.symbolIdx]
			symbolNameStr := mylog.Check2(coffSymbol.String(file))

			if symbolNameStr != tt.out.symbolName {
				t.Errorf(
					"symbol name to string failed, want: %s, got: %s",
					tt.out.symbolName,
					symbolNameStr,
				)
			}

			secNumName := coffSymbol.SectionNumberName(file)
			if secNumName != tt.out.sectionNumberName {
				t.Errorf(
					"SectionNumberName assertion failed, want: %s, got: %s",
					tt.out.sectionNumberName,
					secNumName,
				)
			}

			typeString := file.PrettyCOFFTypeRepresentation(coffSymbol.Type)
			if typeString != tt.out.symbolTypeString {
				t.Errorf(
					"PrettyCOFFTypeRepresentation assertion failed, want: %s, got: %s",
					tt.out.symbolTypeString,
					typeString,
				)
			}
		})
	}
}

func stringInSlice(a string, list []string) bool {
	return slices.Contains(list, a)
}
