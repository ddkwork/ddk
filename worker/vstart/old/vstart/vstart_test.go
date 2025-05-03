package main

import (
	"testing"

	"github.com/ddkwork/golibrary/safemap"

	"github.com/ddkwork/golibrary/stream"
)

func TestGeneratedFile_P(t *testing.T) {
	g := stream.NewGeneratedFile()
	m := new(safemap.M[string, string])
	m.Set("app", "应用集")
	m.Set("git", "git相关")
	m.Set("env", "环境变量")
	m.Set("ai", "人工智能")
	m.Set("notes", "笔记")
	m.Set("json2go", "json转go结构体")
	m.Set("translate", "注释翻译")
	m.Set("webBookMark", "网页收藏夹")
	m.Set("golang", "go语言设置")

	g.EnumTypes("toolbox", m)
}
