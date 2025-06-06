package environment

import (
	"testing"

	"github.com/ddkwork/golibrary/std/assert"

	"github.com/ddkwork/golibrary/std/mylog"
)

func Test_isValidPath(t *testing.T) {
	mylog.Call(func() {
		assert.False(t, isValidPath("D:\\sdk"))
		assert.True(t, isValidPath("%SystemRoot%\\TEMP"))
		assert.True(t, isValidPath("%SystemRoot%\\system32\\cmd.exe"))
		assert.True(t, isValidPath("18 items"))
	})
}

var columnData = []string{
	"18 items",
	"%SystemRoot%\\system32\\cmd.exe",
	"C:\\Windows\\System32\\Drivers\\DriverData",
	"Windows_NT",
	"15 items",
	"C:\\Windows\\System32\\HWAudioDriverLibs",
	"C:\\Windows\\system32",
	"C:\\Windows",
	"C:\\Windows\\System32\\Wbem",
	"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\",
	"C:\\Windows\\System32\\OpenSSH\\",
	"C:\\Program Files\\Git\\cmd",
	"C:\\Users\\Admin\\AppData\\Local\\Microsoft\\WindowsApps",
	"",
	"C:\\Program Files\\JetBrains\\GoLand 2024.1\\bin",
	"",
	"C:\\Users\\Admin\\go\\bin",
	"C:\\Program Files\\dotnet\\",
	"C:\\TDM-GCC-64\\bin",
	"C:\\Program Files\\Go\\bin",
	"11 items",
	".COM",
	".EXE",
	".BAT",
	".CMD",
	".VBS",
	".VBE",
	".JS",
	".JSE",
	".WSF",
	".WSH",
	".MSC",
	"AMD64",
	"2 items",
	"%ProgramFiles%\\WindowsPowerShell\\Modules",
	"%SystemRoot%\\system32\\WindowsPowerShell\\v1.0\\Modules",
	"%SystemRoot%\\TEMP",
	"%SystemRoot%\\TEMP",
	"SYSTEM",
	"%SystemRoot%",
	"8",
	"6",
	"Intel64 Family 6 Model 140 Stepping 1, GenuineIntel",
	"8c01",
	"D:\\sdk",
	"D:\\sdk",
	"1",
}
