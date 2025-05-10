package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"github.com/ddkwork/ddk/worker/BuyTomatoes"
	"github.com/ddkwork/ddk/worker/DecodeRaw"
	"github.com/ddkwork/ddk/worker/jsontree"
	"github.com/ddkwork/ddk/worker/struct2table"

	// "github.com/ddkwork/ddk/worker/DecodeRaw"
	"github.com/ddkwork/ddk/worker/environment"
	"github.com/ddkwork/ddk/worker/explorer"
	"github.com/ddkwork/ddk/worker/godoc"
	// "github.com/ddkwork/ddk/worker/jsontree"
	"github.com/ddkwork/ddk/worker/keygen"
	"github.com/ddkwork/ddk/worker/manpiecework"
	// "github.com/ddkwork/ddk/worker/struct2table"
	"github.com/ddkwork/ddk/worker/taskmanager"
	"github.com/ddkwork/ddk/worker/vstart"
	"github.com/ddkwork/golibrary/safemap"
	"github.com/ddkwork/ux"
)

func main() {
	m := new(safemap.M[WorkType, ux.Widget])
	for _, workType := range ManPieceworkType.EnumTypes() {
		switch workType {
		case BuyTomatoesType:
			m.Set(workType, BuyTomatoes.New())
		case DecodeRawType:
			m.Set(workType, DecodeRaw.New())
		case ManPieceworkType:
			m.Set(workType, manpiecework.New())
		case EnvironmentType:
			m.Set(workType, environment.New())

		// https://github.com/aams-eam/gocleasy todo
		case ExplorerType:
			m.Set(workType, explorer.New())
		case GodocType:
			m.Set(workType, godoc.New())
		case JsontreeType:
			m.Set(workType, jsontree.New())
		case KeygenType:
			m.Set(workType, keygen.New())
		case Struct2tableType:
			m.Set(workType, struct2table.New())
		case TaskmanagerType:
			m.Set(workType, taskmanager.New())
		case VstartType:
			m.Set(workType, vstart.New())
		}
	}
	vtab := ux.NewTabView(layout.Vertical)
	for k, v := range m.Range() {
		tab := ux.NewTabItem(k.String(), v)
		vtab.AddTab(tab)
	}
	panel := ux.NewPanel()
	panel.AddChild(vtab)
	ux.Run("Worker", panel)
}
