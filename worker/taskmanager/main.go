package taskmanager

import (
	"fmt"
	"iter"

	"gioui.org/layout"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/datasize"
	"github.com/ddkwork/ux"
)

func New() ux.Widget {
	t := ux.NewTreeTable(Task{})
	t.TableContext = ux.TableContext[Task]{
		CustomContextMenuItems: func(gtx layout.Context, n *ux.Node[Task]) iter.Seq[ux.ContextMenuItem] {
			return func(yield func(ux.ContextMenuItem) bool) {
			}
		},
		MarshalRowCells: func(n *ux.Node[Task]) (cells []ux.CellData) {
			if n.Container() {
				n.Data.Name = n.SumChildren()
			}
			// return []ux.CellData{
			//	{Text: n.Data.Name, Tooltip: "", SvgBuffer: "", ImageBuffer: nil, FgColor: 0},
			//	{Text: fmt.Sprintf("%.3g%%", n.Data.CPU), Tooltip: "", SvgBuffer: "", ImageBuffer: nil, FgColor: 0},
			//	{Text: n.Data.RAM.String(), Tooltip: "", SvgBuffer: "", ImageBuffer: nil, FgColor: 0},
			//	{Text: fmt.Sprintf("%.3g%%", n.Data.RAMPct), Tooltip: "", SvgBuffer: "", ImageBuffer: nil, FgColor: 0},
			//	{Text: fmt.Sprint(n.Data.Threads), Tooltip: "", SvgBuffer: "", ImageBuffer: nil, FgColor: 0},
			//	{Text: n.Data.User, Tooltip: "", SvgBuffer: "", ImageBuffer: nil, FgColor: 0},
			//	{Text: fmt.Sprint(n.Data.PID), Tooltip: " todo add hex pid", SvgBuffer: "", ImageBuffer: nil, FgColor: 0},
			// }
			return ux.MarshalRow(n.Data, func(key string, field any) (value string) {
				switch key {
				case "CPU":
					return fmt.Sprintf("%.3g%%", field)
				case "RAMPct":
					return fmt.Sprintf("%.3g%%", field)
				default:
					return ""
				}
			})
		},
		UnmarshalRowCells: func(n *ux.Node[Task], rows []ux.CellData) {
			n.Data = ux.UnmarshalRow[Task](rows, func(key, value string) (field any) {
				switch key {
				case "CPU":
					return mylog.Check2(fmt.Sscanf(value, "%f%%"))
				case "RAM":
					return datasize.Parse(value)
				case "RAMPct":
					return mylog.Check2(fmt.Sscanf(value, "%f%%"))
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
			getTasks(t.Root)
			return

			// m := sync.Mutex{}
			// fnUpdate := func() {
			// 	// m.Lock()
			// 	// defer m.Unlock()
			// 	t.Root.ResetChildren()
			// 	// root.SyncToModel()
			// 	getTasks(t.Root)
			// 	t.SetRootRows(t.Root.Children)
			// }

			// ticker := time.NewTicker(time.Second)

			// go func() {
			// 	for range ticker.C {
			// 		unison.InvokeTaskAfter(func() {
			// 			fnUpdate()
			// 		}, time.Second)
			// 	}
			// }()
		},
		JsonName:   "taskmanager",
		IsDocument: false,
	}
	return t
}
