package xed

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/saferwall/pe"
)

func ParserPe(filename string) (file *pe.File) {
	file = mylog.Check2(pe.New(filename, &pe.Options{}))
	mylog.Check(file.Parse())
	return file
}
