package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type WdkConfig struct {
	Root         string
	Version      string
	WinVer       string
	NtddkVersion string
	Platform     string
	CompileFlags []string
	CompileDefs  []string
	LinkFlags    []string
	LibVersion   string
	IncVersion   string
	Architecture string
	IncludeDirs  []string
	Libraries    []string
}

func findWdk() (*WdkConfig, error) {
	var wdkNtddkFiles []string

	wdkContentRoot := os.Getenv("WDKContentRoot")
	if wdkContentRoot != "" {
		err := filepath.Walk(wdkContentRoot, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Base(path) == "ntddk.h" {
				wdkNtddkFiles = append(wdkNtddkFiles, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		paths := []string{
			"C:/Program Files*/Windows Kits/*/Include/*/km/ntddk.h",
			"C:/Program Files*/Windows Kits/*/Include/km/ntddk.h",
		}
		for _, path := range paths {
			matches, err := filepath.Glob(path)
			if err != nil {
				return nil, err
			}
			wdkNtddkFiles = append(wdkNtddkFiles, matches...)
		}
	}

	if len(wdkNtddkFiles) > 0 {
		sort.Strings(wdkNtddkFiles)
		wdkLatestNtddkFile := wdkNtddkFiles[len(wdkNtddkFiles)-1]

		wdkRoot := filepath.Dir(filepath.Dir(wdkLatestNtddkFile))
		wdkVersion := filepath.Base(wdkRoot)
		wdkRoot = filepath.Dir(wdkRoot)

		wdkConfig := &WdkConfig{
			Root:         wdkRoot,
			Version:      wdkVersion,
			WinVer:       "0x0601",
			NtddkVersion: "",
			LibVersion:   wdkVersion,
			IncVersion:   wdkVersion,
			Architecture: "",
			CompileFlags: []string{"/Zp8", "/GF", "/GR-", "/Gz", "/kernel", "/FIwarning.h", "/FIwdkflags.h", "/Oi"},
			CompileDefs:  []string{"WINNT=1"},
			LinkFlags:    []string{"/MANIFEST:NO", "/DRIVER", "/OPT:REF", "/INCREMENTAL:NO", "/OPT:ICF", "/SUBSYSTEM:NATIVE", "/MERGE:_TEXT=.text;_PAGE=PAGE", "/NODEFAULTLIB", "/SECTION:INIT,d", "/VERSION:10.0"},
		}

		if !strings.Contains(wdkRoot, "/[0-9][0-9.]*") {
			wdkRoot = filepath.Dir(wdkRoot)
			wdkConfig.Root = wdkRoot
			wdkConfig.LibVersion = wdkVersion
			wdkConfig.IncVersion = wdkVersion
		} else {
			wdkConfig.IncVersion = ""
			versions := []string{"winv6.3", "win8", "win7"}
			for _, version := range versions {
				if _, err := os.Stat(filepath.Join(wdkRoot, "Lib", version)); err == nil {
					wdkConfig.LibVersion = version
					wdkConfig.Version = version
					break
				}
			}
		}

		// Determine the architecture
		architecture := os.Getenv("CMAKE_CXX_COMPILER_ARCHITECTURE_ID")
		architecture = runtime.GOARCH
		switch architecture {
		case "x86":
			wdkConfig.Architecture = "x86"
			wdkConfig.CompileDefs = append(wdkConfig.CompileDefs, "_X86_=1", "i386=1", "STD_CALL")
		case "ARM64":
			wdkConfig.Architecture = "arm64"
			wdkConfig.CompileDefs = append(wdkConfig.CompileDefs, "_ARM64_", "ARM64", "_USE_DECLSPECS_FOR_SAL=1", "STD_CALL")
		case "AMD64":
			wdkConfig.Architecture = "x64"
			wdkConfig.CompileDefs = append(wdkConfig.CompileDefs, "_AMD64_", "AMD64")
		default:
			return nil, fmt.Errorf("Unsupported architecture")
		}

		// Set include directories
		wdkConfig.IncludeDirs = []string{
			filepath.Join(wdkConfig.Root, "Include", wdkConfig.IncVersion, "shared"),
			filepath.Join(wdkConfig.Root, "Include", wdkConfig.IncVersion, "km"),
			filepath.Join(wdkConfig.Root, "Include", wdkConfig.IncVersion, "km/crt"),
		}

		// Set libraries
		libPath := filepath.Join(wdkConfig.Root, "Lib", wdkConfig.LibVersion, "km", wdkConfig.Platform, "*.lib")
		matches, err := filepath.Glob(libPath)
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			wdkConfig.Libraries = append(wdkConfig.Libraries, match)
		}

		return wdkConfig, nil
	}

	return nil, fmt.Errorf("WDK not found")
}

func main() {
	wdkConfig, err := findWdk()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("WDK_ROOT: %s\n", wdkConfig.Root)
	fmt.Printf("WDK_VERSION: %s\n", wdkConfig.Version)
	fmt.Printf("WDK_COMPILE_FLAGS: %v\n", wdkConfig.CompileFlags)
	fmt.Printf("WDK_COMPILE_DEFINITIONS: %v\n", wdkConfig.CompileDefs)
	fmt.Printf("WDK_LINK_FLAGS: %v\n", wdkConfig.LinkFlags)
	fmt.Printf("WDK_INCLUDE_DIRS: %v\n", wdkConfig.IncludeDirs)
	fmt.Printf("WDK_LIBRARIES: %v\n", wdkConfig.Libraries)

	// Simulate adding a driver and library
	if wdkConfig.Architecture == "x86" {
		wdkConfig.CompileFlags = append(wdkConfig.CompileFlags, "/ENTRY:FxDriverEntry@8")
	} else if wdkConfig.Architecture == "arm64" || wdkConfig.Architecture == "x64" {
		wdkConfig.CompileFlags = append(wdkConfig.CompileFlags, "/ENTRY:FxDriverEntry")
	}

	// Simulate including additional flags file
	additionalFlagsFile := filepath.Join(os.Getenv("CMAKE_CURRENT_BINARY_DIR"), os.Getenv("CMAKE_FILES_DIRECTORY"), "wdkflags.h")
	if err := os.WriteFile(additionalFlagsFile, []byte("#pragma runtime_checks(\"suc\", off)"), 0o644); err != nil {
		fmt.Println(err)
		return
	}

	// Simulate compile definitions for debug mode
	if os.Getenv("CMAKE_BUILD_TYPE") == "Debug" {
		wdkConfig.CompileDefs = append(wdkConfig.CompileDefs, "MSC_NOOPT", "DEPRECATE_DDK_FUNCTIONS=1", "DBG=1")
	}

	// Simulate setting NTDDI_VERSION if specified
	if wdkConfig.NtddkVersion != "" {
		wdkConfig.CompileDefs = append(wdkConfig.CompileDefs, fmt.Sprintf("NTDDI_VERSION=%s", wdkConfig.NtddkVersion))
	}

	// Simulate setting KMDF if specified
	if kmdfVersion := os.Getenv("WDK_KMDF"); kmdfVersion != "" {
		wdkConfig.IncludeDirs = append(wdkConfig.IncludeDirs, filepath.Join(wdkConfig.Root, "Include", "wdf", "kmdf", kmdfVersion))
		wdkConfig.Libraries = append(wdkConfig.Libraries,
			filepath.Join(wdkConfig.Root, "Lib", "wdf", "kmdf", wdkConfig.Platform, kmdfVersion, "WdfDriverEntry.lib"),
			filepath.Join(wdkConfig.Root, "Lib", "wdf", "kmdf", wdkConfig.Platform, kmdfVersion, "WdfLdr.lib"),
		)
	}

	// Simulate setting platform specific libraries
	if wdkConfig.Architecture == "arm64" {
		wdkConfig.Libraries = append(wdkConfig.Libraries, "arm64rt.lib")
	}
	if wdkConfig.Architecture == "x86" {
		wdkConfig.Libraries = append(wdkConfig.Libraries, "wdk::MEMCMP")
	}

	fmt.Printf("WDK_COMPILE_FLAGS (final): %v\n", wdkConfig.CompileFlags)
	fmt.Printf("WDK_COMPILE_DEFINITIONS (final): %v\n", wdkConfig.CompileDefs)
	fmt.Printf("WDK_LINK_FLAGS (final): %v\n", wdkConfig.LinkFlags)
	fmt.Printf("WDK_INCLUDE_DIRS (final): %v\n", wdkConfig.IncludeDirs)
	fmt.Printf("WDK_LIBRARIES (final): %v\n", wdkConfig.Libraries)
}
