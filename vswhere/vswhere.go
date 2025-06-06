package vswhere

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
)

// todo move into library
type (
	WdkConfig struct {
		Root              string
		Version           string
		VsInstallDir      string
		VsVersion         string // 2015, 2017, 2019, 2022
		VsType            string // community professional Enterprise
		Includes          []string
		Libs              []string // todo
		CC                string
		Mt                string
		Ar                string
		CXX               string
		Linker            string
		DevEnvDir         string
		RcCompiler        string
		AsmCompiler       string
		VCInstallDir      string
		VS170ComnTools    string
		VCIDEInstallDir   string
		VCToolsRedistDir  string
		VCToolsInstallDir string
	}
)

// CMAKE_CXX_STANDARD_LIBRARIES
// kernel32.lib user32.lib gdi32.lib winspool.lib shell32.lib ole32.lib oleaut32.lib uuid.lib comdlg32.lib advapi32.lib
//
// LIB
// C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC\Tools\MSVC\14.42.34433\ATLMFC\lib\x64;
// C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC\Tools\MSVC\14.42.34433\lib\x64;
// C:\Program Files (x86)\Windows Kits\NETFXSDK\4.8\lib\um\x64;
// C:\Program Files (x86)\Windows Kits\10\lib\10.0.26100.0\ucrt\x64;
// C:\Program Files (x86)\Windows Kits\10\\lib\10.0.26100.0\\um\x64
//
// LIBPATH
// C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC\Tools\MSVC\14.42.34433\ATLMFC\lib\x64;
// C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC\Tools\MSVC\14.42.34433\lib\x64;
// C:\Program Files\Microsoft Visual Studio\2022\Enterprise\VC\Tools\MSVC\14.42.34433\lib\x86\store\references;
// C:\Program Files (x86)\Windows Kits\10\UnionMetadata\10.0.26100.0;
// C:\Program Files (x86)\Windows Kits\10\References\10.0.26100.0;
// C:\Windows\Microsoft.NET\Framework64\v4.0.30319

func New() (v *WdkConfig) {
	return &WdkConfig{
		Root:              "",
		Version:           "",
		VsInstallDir:      "",
		VsVersion:         "",
		VsType:            "",
		Includes:          make([]string, 0),
		Libs:              make([]string, 0),
		Mt:                "",
		Ar:                "",
		Linker:            "",
		DevEnvDir:         "",
		CC:                "",
		CXX:               "",
		RcCompiler:        "",
		AsmCompiler:       "",
		VCInstallDir:      "",
		VS170ComnTools:    "",
		VCIDEInstallDir:   "",
		VCToolsRedistDir:  "",
		VCToolsInstallDir: "",
	}
}

func (c *WdkConfig) VisualStudio() *WdkConfig {
	// C:\Program Files\Microsoft Visual Studio\2022\Enterprise
	var (
		versions = []string{
			"2015",
			"2017",
			"2019",
			"2022",
		}
		Types = []string{
			"community",
			"professional",
			"Enterprise",
		}
	)

	for _, version := range versions {
		p := filepath.Join("C:\\Program Files\\Microsoft Visual Studio", version)
		if stream.IsDir(p) {
			c.VsVersion = version
			for _, l := range Types {
				d := filepath.Join(p, l)
				if stream.IsDir(d) {
					c.VsInstallDir = filepath.ToSlash(d)
					c.VsType = l
					break
				}
			}
			break
		}
	}
	mylog.Check(filepath.Walk(c.VsInstallDir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Base(path) == "cl.exe" && strings.Contains(path, "Hostx64\\x64") {
			path = filepath.ToSlash(path)
			c.CC = path
			c.CXX = path
			return err
		}
		return err
	}))
	c.findWdk()
	c.Mt = filepath.ToSlash(filepath.Join(c.Root, "bin", c.Version, "x64", "mt.exe"))
	c.RcCompiler = filepath.ToSlash(filepath.Join(c.Root, "bin", c.Version, "x64", "rc.exe"))

	dir := filepath.ToSlash(filepath.Dir(c.CC))
	c.AsmCompiler = filepath.ToSlash(filepath.Join(dir, "ml64.exe"))
	c.Ar = filepath.ToSlash(filepath.Join(dir, "lib.exe"))
	c.Linker = filepath.ToSlash(filepath.Join(dir, "link.exe"))
	c.DevEnvDir = filepath.ToSlash(filepath.Join(c.VsInstallDir, "Common7", "IDE"))
	VCToolsVersion := findVersion(dir, "MSVC", "bin")
	c.VCIDEInstallDir = filepath.ToSlash(filepath.Join(c.VsInstallDir, "Common7", "IDE", "VC"))
	c.VCInstallDir = filepath.ToSlash(filepath.Join(c.VsInstallDir, "VC"))
	c.VCToolsInstallDir = filepath.ToSlash(filepath.Join(c.VsInstallDir, "VC", "Tools", "MSVC", VCToolsVersion))
	c.VCToolsRedistDir = filepath.ToSlash(filepath.Join(c.VsInstallDir, "VC", "Redist", "MSVC", VCToolsVersion))
	c.VS170ComnTools = filepath.ToSlash(filepath.Join(c.VsInstallDir, "Common7", "Tools"))

	kmdfVersions := make([]string, 0)
	root := filepath.Join(c.Root, "Include", "wdf", "kmdf")
	mylog.Check(filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			if filepath.Base(path) == "kmdf" {
				return err
			}
			kmdfVersions = append(kmdfVersions, filepath.Base(path))
		}
		return err
	}))

	for _, s := range []string{"shared", "km", "km/crt"} {
		c.Includes = append(c.Includes, filepath.Join(c.Root, "Include", c.Version, s))
	}
	kmdfVersionValues := make([]float64, 0)
	for _, version := range kmdfVersions {
		split := strings.Split(version, ".")
		kmdfVersionValues = append(kmdfVersionValues, mylog.Check2(strconv.ParseFloat(split[1], 64)))
	}
	kmdfVersionMax := 1.0 + slices.Max(kmdfVersionValues)/100.0
	//	mylog.Success("max wdf kmdf version", kmdfVersionMax)
	c.Includes = append(c.Includes, filepath.Join(root, strconv.FormatFloat(kmdfVersionMax, 'f', 2, 64)))
	c.Includes = append(c.Includes, filepath.Join(c.VsInstallDir, "VC", "Tools", "MSVC", VCToolsVersion, "include"))
	c.Includes = append(c.Includes, filepath.Join(c.VsInstallDir, "VC", "Tools", "MSVC", VCToolsVersion, "ATLMFC", "include"))
	c.Includes = append(c.Includes, filepath.Join(c.VsInstallDir, "VC", "Auxiliary", "VS", "include"))
	for _, s := range []string{"ucrt", "um", "shared", "winrt", "cppwinrt"} {
		c.Includes = append(c.Includes, filepath.Join(c.Root, "Include", c.Version, s))
	}
	for i, include := range c.Includes {
		include = filepath.Clean(include)
		include = filepath.ToSlash(include)
		include = strings.ReplaceAll(include, " ", "^")
		include = strconv.Quote(include)
		c.Includes[i] = include
	}
	//mylog.Struct(c)
	return c
}

func (c *WdkConfig) findWdk() (ntddk string) {
	for driver := range stream.GetWindowsLogicalDrives() {
		driver = filepath.Join(driver, "Program Files (x86)", "Windows Kits")
		if stream.IsDir(driver) {
			c.Root = filepath.Join(driver, "10") // todo
			c.Root = filepath.ToSlash(c.Root)
			//mylog.Success("wdk root", c.Root)
			break
		}
	}
	if c.Root == "" {
		panic("wdk root not found")
	}

	mylog.Check(filepath.WalkDir(c.Root, func(path string, d fs.DirEntry, err error) error {
		if filepath.Base(path) == "ntddk.h" {
			ntddk = filepath.ToSlash(path)
			//mylog.Success("ntddk.h", ntddk)
			return err
		}
		return err
	}))

	// C:/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/km/ntddk.h
	cut := findVersion(ntddk, "Include", "km")
	c.Version = cut
	//	mylog.Success("wdk version", c.Version)
	return
}

// findVersion "C:/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/km/ntddk.h", --->10.0.26100.0
func findVersion(path, left, right string) string {
	for s := range strings.SplitSeq(path, "/") {
		if strings.Index(s, ".") > 1 && !strings.HasSuffix(s, ".h") {
			return s
		}
	}
	panic("version not found")
}
