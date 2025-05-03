package drivertool

import (
	_ "embed"
)

//go:embed icon.png
var icon []byte

// func main() {
// 	app.Run("driver tool", func(w *unison.Window) {
// 		img := mylog.Check2(unison.NewImageFromBytes(icon, 0.5))
// 		w.SetTitleIcons([]*unison.Image{img})
// 		w.Content().AddChild(New().Layout())
// 	})
// }

type DriverLoad struct {
	ReloadPath string
	Link       string
	IoCode     string
}

type StructView struct{}

// func New() ux.Widget {
// 	return &StructView{}
// }
//
// func (s *StructView) Layout() ux.Widget {
// 	newPanel := widget.NewPanel()
//
// 	view := DriverLoad{
// 		ReloadPath: "",
// 		Link:       "",
// 		IoCode:     "",
// 	}
// 	structView, rowPanel := widget.NewStructView(view, func(data DriverLoad) (values []widget.CellData) {
// 		return []widget.CellData{
// 			{ImageBuffer: nil, Text: data.ReloadPath, Tooltip: "", FgColor: 0},
// 			{ImageBuffer: nil, Text: data.Link, Tooltip: "", FgColor: 0},
// 			{ImageBuffer: nil, Text: data.IoCode, Tooltip: "", FgColor: 0},
// 		}
// 	})
//
// 	p := unison.NewPopupMenu[string]()
//
// 	p.SelectionChangedCallback = func(popup *unison.PopupMenu[string]) {
// 		if title, ok := popup.Selected(); ok {
// 			structView.MetaData.ReloadPath = title
// 			structView.MetaData.Link = stream.BaseName(title)
// 			structView.UpdateField(0, title)
// 			structView.UpdateField(1, stream.BaseName(title))
// 		}
// 	}
//
// 	root := "../"
// 	abs := mylog.Check2(filepath.Abs(root))
// 	names := make([]string, 0)
// 	if abs == "C:\\Users\\Admin" {
// 		names = WalkAllDriverPath(".")
// 	} else {
// 		names = WalkAllDriverPath(root)
// 	}
//
// 	popupMenu := widget.CreatePopupMenu(newPanel.AsPanel(), p, 0, "choose a driver", names...)
//
// 	key := widget.NewLabelRightAlign(widget.KeyValueToolTip{
// 		Key:     "sys path",
// 		Value:   "",
// 		Tooltip: "",
// 	})
// 	rowPanel.AddChildAtIndex(key, 0)
// 	rowPanel.AddChildAtIndex(popupMenu, 1)
// 	newPanel.AddChild(structView)
//
// 	log := unison.NewMultiLineField() // todo log out format is not good
// 	log.MinimumTextWidth = 800
// 	log.SetText(`log view
//
//
//
//
//
//
// `)
// 	d := New()
// 	panel := widget.NewButtonsPanel(
// 		[]string{"load", "unload"},
// 		func() {
// 			d.Load(structView.MetaData.ReloadPath)
// 			log.SetText(mylog.Row())
// 		},
// 		func() {
// 			d.Unload()
// 			log.SetText(mylog.Row())
// 		},
// 	)
// 	rowPanel.AddChild(widget.NewVSpacer())
// 	rowPanel.AddChild(panel)
// 	newPanel.AddChild(log)
// 	return newPanel
// }
//
// func WalkAllDriverPath(root string) (drivers []string) {
// 	drivers = make([]string, 0)
// 	mylog.Check(filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
// 		if filepath.Ext(path) == ".sys" {
// 			drivers = append(drivers, path)
// 		}
// 		return nil
// 	}))
// 	return
// }
