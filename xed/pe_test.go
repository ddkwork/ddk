package xed

import (
	"testing"
)

func TestParserPe(t *testing.T) {
	t.Skip()
	p := "D:\\迅雷下载\\microsoft.windows.wdk.win32metadata.0.9.9-experimental.nupkg"
	file := ParserPe(p)
	for _, exception := range file.Exceptions {
		println(exception.UnwindInfo.ExceptionHandler)
	}
}
