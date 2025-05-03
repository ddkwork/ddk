package keygen

import (
	"github.com/ddkwork/ux"
)

type Object struct {
	MachineID string
	RegCode   string
	Version   string
	Website   string
	SimpleMid string
}

func New() ux.Widget {
	//keygen.Api()//todo

	data := Object{
		MachineID: "1111-2222-3333-4444",
		RegCode:   "aaa-bbb-ccc-ddd",
		Version:   "1.1.1",
		Website:   "https://www.baidu.com",
		SimpleMid: "2222-3333-4444-5555",
	}

	form := ux.NewStructView("edit node meta data", data, func(key string, field any) (value string) {
		return "" //todo test custom format field
	},
	)
	form.SetOnApply(func() {
		form.Data = ux.UnmarshalRow[Object](form.Rows, func(key, value string) (field any) {
			return nil //todo test custom format field
		})
	})

	//todo add select app type
	//p := unison.NewPopupMenu[string]()
	//appKind := keygen.InvalidAppKind
	//p.SelectionChangedCallback = func(popup *unison.PopupMenu[string]) {
	//	if title, ok := popup.Selected(); ok {
	//		appKind = keygen.InvalidAppKind.AssertKind(title)
	//		mylog.Trace("selected appKind", title)
	//	}
	//}
	//popupMenu := widget.CreatePopupMenu(s.RowPanel, p, 0, "choose a app", keygen.InvalidAppKind.Keys()...)
	//s.RowPanel.AddChildAtIndex(widget.NewLabelRightAlign(widget.KeyValueToolTip{
	//	Key:     "app name",
	//	Value:   "",
	//	Tooltip: "",
	//}), 0)
	//s.RowPanel.AddChildAtIndex(popupMenu, 1)
	//panel := widget.NewButtonsPanel(
	//	[]string{"apply", "cancel"},
	//	func() {
	//		mylog.Warning(s.View.MetaData.MachineID)
	//		var g keygen.Api
	//		s.View.MetaData.MachineID = s.View.Editors[0].Label.String()
	//		switch appKind {
	//		case keygen.SuperRecovery4Kind:
	//			g = keygen.NewSuperRecovery4()
	//		default:
	//			panic("unhandled default case")
	//		}
	//		mylog.Check(g.GenRegCode(s.View.MetaData.MachineID)) // todo get from filed or unmarshal method
	//		s.View.MetaData.RegCode = g.RegCode()
	//		s.View.UpdateField(1, s.View.MetaData.RegCode)
	//	},
	//	func() {},
	//)

	// userName := ux.NewInput("please input username")
	// password := ux.NewInput("please input password")
	// email := ux.NewInput("please input email")

	// form.Add("username", userName.Layout)
	// form.Add("password", password.Layout)
	// form.Add("email", email.Layout)
	// dropDown := ux.NewDropDown(SuperRecovery2Type.Names()...)
	//dropDown := ux.NewDropDown()
	//for _, s := range SuperRecovery2Type.Names() {
	//	dropDown.SetOptions(ux.NewDropDownOption(s))
	//}

	// form.InsertAt(0, "choose a app", dropDown.Layout)
	// form.Add("", ux.BlueButton(&clickable, "submit", unit.Dp(100)).Layout)
	return form
}
