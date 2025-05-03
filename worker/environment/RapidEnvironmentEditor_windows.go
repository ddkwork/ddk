package environment

import (
	_ "embed"
	"gioui.org/layout"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/ux"
	"golang.org/x/sys/windows/registry"
	"iter"
	"os"
	"path/filepath"
	"strings"
)

//go:generate go build -x .
//go:generate go install .

// https://faststone-photo-resizer.en.lo4d.com/windows

type EnvironmentEditor struct {
	Key     string
	Value   string
	IsValid bool
	Type    kind
}

func Layout() ux.Widget {
	t := ux.NewTreeTable(EnvironmentEditor{})
	t.TableContext = ux.TableContext[EnvironmentEditor]{
		CustomContextMenuItems: func(gtx layout.Context, n *ux.Node[EnvironmentEditor]) iter.Seq[ux.ContextMenuItem] {
			return func(yield func(ux.ContextMenuItem) bool) {

			}
		},
		MarshalRowCells: func(n *ux.Node[EnvironmentEditor]) (cells []ux.CellData) {
			if n.Container() {
				n.Data.Key = n.SumChildren()
				n.Data.Value = ""
			}
			n.Data.IsValid = isValidPath(n.Data.Value)
			enable := !n.Data.IsValid
			status := "✓"
			if enable {
				status = "✗"
			}
			return ux.MarshalRow(n.Data, func(key string, field any) (value string) {
				if key == "IsValid" {
					return status
				}
				return ""
			})
		},
		UnmarshalRowCells: func(n *ux.Node[EnvironmentEditor], rows []ux.CellData) {
			n.Data = ux.UnmarshalRow[EnvironmentEditor](rows, func(key, value string) (field any) {
				if key == "IsValid" {
					if value == "✓" {
						return true
					}
					return false
				}
				return nil
			})
		},
		RowSelectedCallback: func() {

		},
		RowDoubleClickCallback: func() {

		},
		SetRootRowsCallBack: func() {
			const EnvPath = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
			key := mylog.Check2(registry.OpenKey(registry.LOCAL_MACHINE, EnvPath, registry.ALL_ACCESS))
			valueNames := mylog.Check2(key.ReadValueNames(-1))
			for _, valueName := range valueNames {
				value, valueType := mylog.Check3(key.GetStringValue(valueName))
				layout_(t.Root, EnvironmentEditor{
					Key:     valueName,
					Value:   value,
					IsValid: !isValidPath(value),
					Type:    kind(valueType),
				})
			}
		},
		JsonName:   "RapidEnvironmentEditor",
		IsDocument: false,
	}
	return t
}

func layout_(parent *ux.Node[EnvironmentEditor], data EnvironmentEditor) {
	if strings.Contains(data.Value, ";") {
		v := data.Value
		container := ux.NewContainerNode(data.Key, data)
		parent.AddChild(container)
		for value := range strings.SplitSeq(v, ";") {
			container.AddChildByData(EnvironmentEditor{
				Key:     data.Key,
				Value:   value,
				IsValid: !isValidPath(value),
				Type:    data.Type,
			})
		}
		return
	}
	parent.AddChildByData(data)
}

func isValidPath(path string) bool {
	switch {
	case strings.Contains(path, "items"): // contains items
		return true
	case path == "": // path or not path,this is wrong value,need remove it,todo one click remove
		return false
	case !strings.Contains(path, "\\") && !strings.Contains(path, "/"):
		return true // skip not path
	case strings.HasPrefix(path, `%`): // decode env path
		split := strings.Split(path, `%`)
		env := os.Getenv(split[1])
		path = filepath.Join(env, split[2])
		ext := filepath.Ext(path)
		if ext == "" {
			return stream.IsDirEx(path)
		}
		return stream.IsFilePathEx(path)
	default: // must have filepath.Separator
		return stream.IsDirEx(path)
	}
}

type kind byte

const (
	// Registry value types.
	NONE kind = iota
	SZ
	EXPAND_SZ
	BINARY
	DWORD
	DWORD_BIG_ENDIAN
	LINK
	MULTI_SZ
	RESOURCE_LIST
	FULL_RESOURCE_DESCRIPTOR
	RESOURCE_REQUIREMENTS_LIST
	QWORD
)

func (k kind) String() string {
	switch k {
	case NONE:
		return "NONE"
	case SZ:
		return "SZ"
	case EXPAND_SZ:
		return "EXPAND_SZ"
	case BINARY:
		return "BINARY"
	case DWORD:
		return "DWORD"
	case DWORD_BIG_ENDIAN:
		return "DWORD_BIG_ENDIAN"
	case LINK:
		return "LINK"
	case MULTI_SZ:
		return "MULTI_SZ"
	case RESOURCE_LIST:
		return "RESOURCE_LIST"
	case FULL_RESOURCE_DESCRIPTOR:
		return "FULL_RESOURCE_DESCRIPTOR"
	case RESOURCE_REQUIREMENTS_LIST:
		return "RESOURCE_REQUIREMENTS_LIST"
	case QWORD:
		return "QWORD"
	default:
		return "unknown"
	}
}
