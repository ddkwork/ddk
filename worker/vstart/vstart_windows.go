//go:build windows

package vstart

import (
	_ "embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/desktop"
	"github.com/ddkwork/ux"
	"github.com/ddkwork/ux/resources/images"
)

var (

	//go:embed appicon.png
	VStartPng []byte

	//go:embed  installer.png
	installerJpg []byte
)

func New() ux.Widget {
	// todo walk目录触发重新布局，当前的WalkDir是一次性的，实时WalkDir的话性能不好，我们需要一个高性能的waik模块，发生改变之后重新执行该函数重新渲染
	paths := make([]string, 0)
	mylog.Check(filepath.WalkDir("d:\\app", func(path string, info fs.DirEntry, err error) error {
		switch {
		case strings.Contains(path, "RECYCLE.BIN"):
			return err
		case info.IsDir():
			return err
		}
		ext := filepath.Ext(path)
		switch ext {
		case ".exe": // msi invalid argument , not support
			if stream.IsWindows() {
				path = filepath.ToSlash(path)
				oldPng := path[:len(path)-len(filepath.Ext(path))] + ".png"
				if stream.IsFilePath(oldPng) {
					mylog.Check(os.Remove(oldPng))
				}
				paths = append(paths, path)
			}
		}
		return err
	}))

	elems := new(safemap.M[*ux.ButtonAnimation, *ux.ContextMenu])
	for _, path := range paths {
		key := ux.NewButton(new(widget.Clickable), ExtractIcon2Png(path), stream.AlignString(stream.BaseName(path), 5), func(gtx layout.Context) {
			stream.RunCommandArgs("start", path)
		})
		value := ux.NewContextMenu()
		value.Once.Do(func() {
			value.AddItem(ux.ContextMenuItem{
				Title: "打开所在目录",
				Icon:  images.FileFolderOpenIcon,
				Can:   func() bool { return true },
				Do: func() {
					go desktop.Open(filepath.Dir(path))
					// go stream.RunCommandArgs("cd ", filepath.Dir(button.Tooltip.String()), "start", button.Tooltip.String())
				},
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
			value.AddItem(ux.ContextMenuItem{
				Title: "删除",
				Icon:  images.FileFolderOpenIcon, // todo 换图标
				Can:   func() bool { return true },
				Do: func() {
					// go desktop.Open(filepath.Dir(path))
					// go stream.RunCommandArgs("cd ", filepath.Dir(button.Tooltip.String()), "start", button.Tooltip.String())
				},
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
			value.AddItem(ux.ContextMenuItem{
				Title: "剪切",
				Icon:  images.FileFolderOpenIcon,
				Can:   func() bool { return true },
				Do: func() {
					// go desktop.Open(filepath.Dir(path))
					// go stream.RunCommandArgs("cd ", filepath.Dir(button.Tooltip.String()), "start", button.Tooltip.String())
				},
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
			value.AddItem(ux.ContextMenuItem{
				Title: "粘贴",
				Icon:  images.FileFolderOpenIcon,
				Can:   func() bool { return true },
				Do: func() {
					// go desktop.Open(filepath.Dir(path))
					// go stream.RunCommandArgs("cd ", filepath.Dir(button.Tooltip.String()), "start", button.Tooltip.String())
				},
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
			value.AddItem(ux.ContextMenuItem{
				Title: "todo",
				Icon:  images.FileFolderOpenIcon,
				Can:   func() bool { return true },
				Do: func() {
					// go desktop.Open(filepath.Dir(path))
					// go stream.RunCommandArgs("cd ", filepath.Dir(button.Tooltip.String()), "start", button.Tooltip.String())
				},
				AppendDivider: false,
				Clickable:     widget.Clickable{},
			})
		})
		elems.Set(key, value)
	}
	flow := ux.NewFlow(8, elems)
	return flow
}
