package explorer

import (
	"iter"
	"os"
	"path/filepath"
	"time"

	"gioui.org/layout"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/datasize"
	"github.com/ddkwork/ux"

	"github.com/ddkwork/golibrary/mylog"
)

func New() ux.Widget {
	t := ux.NewTreeTable(DirTree{})
	t.TableContext = ux.TableContext[DirTree]{
		CustomContextMenuItems: func(gtx layout.Context, n *ux.Node[DirTree]) iter.Seq[ux.ContextMenuItem] {
			return func(yield func(ux.ContextMenuItem) bool) {
			}
		},
		MarshalRowCells: func(n *ux.Node[DirTree]) (cells []ux.CellData) {
			if n.Container() {
				n.Data.Name = n.SumChildren()
				sum := int64(0)
				n.Data.Size = sum
				for _, n := range n.Walk() {
					sum += n.Data.Size
				}
				n.Data.Size = sum
			}
			return ux.MarshalRow(n.Data, func(key string, field any) (value string) {
				switch key {
				case "Size":
					return datasize.Size(n.Data.Size).String()
				case "ModTime":
					return stream.FormatTime(n.Data.ModTime)
				default:
					return ""
				}
			})
		},
		UnmarshalRowCells: func(n *ux.Node[DirTree], rows []ux.CellData) {
			n.Data = ux.UnmarshalRow[DirTree](rows, func(key, value string) (field any) {
				switch key {
				case "Size":
					return int64(datasize.Parse(value))
				case "ModTime":
					return mylog.Check2(time.Parse(time.RFC3339, value))
				default:
					return nil
				}
			})
		},
		RowSelectedCallback: func() {
		},
		RowDoubleClickCallback: func() {
		},
		SetRootRowsCallBack: func() {
			Walk(".", t.Root)
		},
		JsonName:   "explorer",
		IsDocument: false,
	}
	return t
}

// https://github.com/y4v8/filewatcher
// https://github.com/Ronbb/usn
type (
	DirTree struct {
		Name    string
		Size    int64
		Type    string
		ModTime time.Time
		Path    string
	}
)

func Walk(path string, parent *ux.Node[DirTree]) (ok bool) {
	parent.Data.Type = filepath.Base(path) // 设置root的type，在格式化回调中赋值到Name
	dirEntries := mylog.Check2(os.ReadDir(path))
	for _, entry := range dirEntries {
		info := mylog.Check2(entry.Info())
		dirTree := DirTree{
			Name:    entry.Name(),
			Size:    info.Size(),
			Type:    entry.Type().String(),
			ModTime: info.ModTime(),
			Path:    filepath.Join(path, entry.Name()),
		}
		if entry.IsDir() {
			containerNode := ux.NewContainerNode(entry.Name(), dirTree)
			parent.AddChild(containerNode)
			Walk(filepath.Join(path, entry.Name()), containerNode)
			continue
		}
		parent.AddChild(ux.NewNode(dirTree))
	}
	return true
}
