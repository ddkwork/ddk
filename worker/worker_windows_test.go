package main

import (
	"io/fs"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/safemap"
	"github.com/ddkwork/golibrary/std/stream"
)

func TestWalkWorker(t *testing.T) {
	g := stream.NewGeneratedFile()
	g.P("func TestGenWork(t *testing.T) {")
	g.P("m := safemap.NewOrdered[string, string](func(yield func(string, string) bool) {")
	mylog.Check(filepath.Walk("../worker", func(path string, info fs.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			if stream.DirDepth(path) == 2 {
				key := filepath.Base(path)
				g.P("yield(", strconv.Quote(key), ", ", strconv.Quote(key), ")")
			}
		}
		return err
	}))
	g.P(`	})
	stream.NewGeneratedFile().SetPackageName("main").EnumTypes("work", m)`)
	g.P("}")
	println(g.Format())
}

func TestGenWork(t *testing.T) {
	m := safemap.NewOrdered[string, string](func(yield func(string, string) bool) {
		yield("BuyTomatoes", "BuyTomatoes")
		yield("DecodeRaw", "DecodeRaw")
		yield("ManPiecework", "ManPiecework")
		yield("environment", "environment")
		yield("explorer", "explorer")
		yield("godoc", "godoc")
		yield("jsontree", "jsontree")
		yield("keygen", "keygen")
		yield("struct2table", "struct2table")
		yield("taskmanager", "taskmanager")
		yield("vstart", "vstart")
	})
	stream.NewGeneratedFile().SetPackageName("main").EnumTypes("work", m)
}
