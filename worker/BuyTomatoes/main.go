package BuyTomatoes

import (
	"iter"
	"time"

	"gioui.org/layout"
	"github.com/ddkwork/golibrary/std/stream"
	"github.com/ddkwork/ux"
)

func New() ux.Widget {
	t := ux.NewTreeTable(BuyTomato{})
	t.TableContext = ux.TableContext[BuyTomato]{
		CustomContextMenuItems: func(gtx layout.Context, n *ux.Node[BuyTomato]) iter.Seq[ux.ContextMenuItem] {
			return func(yield func(ux.ContextMenuItem) bool) {
			}
		},
		MarshalRowCells: func(n *ux.Node[BuyTomato]) (cells []ux.CellData) {
			timeData := stream.FormatTime(n.Data.Date)
			if n.Container() {
				timeData = n.SumChildren()
				sum := int64(0)
				n.Data.NetWeight = sum
				for _, node := range n.Walk() {
					sum += node.Data.NetWeight
				}
				n.Data.NetWeight = sum

				sum2 := .0
				n.Data.Settled = sum2
				for _, node := range n.Walk() {
					sum2 += node.Data.Settled
				}
				n.Data.Settled = sum2

				sum2 = .0
				n.Data.Unsettled = sum2
				for _, node := range n.Walk() {
					sum2 += node.Data.Unsettled
				}
				n.Data.Unsettled = sum2

				sum2 = .0
				n.Data.Payment = sum2
				for _, node := range n.Walk() {
					sum2 += node.Data.Payment
				}
				n.Data.Payment = sum2
			}
			// return []ux.CellData{
			//	{Text: timeData},
			//	{Text: node.Data.Name},
			//	{Text: strconv.FormatInt(node.Data.GrossWeight, 10)},
			//	{Text: strconv.FormatInt(node.Data.Wagon, 10)},
			//	{Text: strconv.FormatInt(node.Data.NetWeight, 10)},
			//	{Text: fmt.Sprint(node.Data.Price)},
			//	{Text: fmt.Sprint(node.Data.Payment)},
			//	{Text: node.Data.Phone},
			//	{Text: stream.FormatTime(node.Data.SettlementDate)},
			//	{Text: node.Data.Settler},
			//	{Text: fmt.Sprint(node.Data.Settled)},
			//	{Text: fmt.Sprint(node.Data.Unsettled)},
			//	{Text: fmt.Sprint(node.Data.Tags)}, // todo
			//	{Text: fmt.Sprint(node.Data.Score)},
			//	{Text: fmt.Sprint(node.Data.Enable)},
			//	{Text: node.Data.Note},
			// }
			return ux.MarshalRow(n.Data, func(key string, field any) (value string) {
				switch key {
				case "Date":
					return timeData
				// case "ModTime":
				//	return stream.FormatTime(n.Data.ModTime)
				default:
					return ""
				}
			})
		},
		UnmarshalRowCells: func(n *ux.Node[BuyTomato], rows []ux.CellData) BuyTomato {
			return ux.UnmarshalRow[BuyTomato](rows, func(key, value string) (field any) {
				switch key {
				// case "Size":
				//	return int64(datasize.Parse(value))
				// case "ModTime":
				//	return mylog.Check2(time.Parse(time.RFC3339, value))
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
			tomatoes := make([]*ux.Node[BuyTomato], 0)
			now := time.Now()
			for range 100 {
				container := ux.NewContainerNode(now.Format("2006-01-02"),
					BuyTomato{
						Date:           time.Time{},
						Name:           "",
						GrossWeight:    10,
						Wagon:          11,
						NetWeight:      12,
						Price:          13,
						Payment:        14,
						Phone:          "01234567890",
						SettlementDate: now,
						Settler:        "asd",
						Settled:        30,
						Unsettled:      40,
						Tags:           make([]string, 0),
						Score:          1.0,
						Enable:         false,
						Note:           "this is a note",
					})
				container.SetParent(t.Root)
				container.SetOpen(true)
				tomatoes = append(tomatoes, container)
			}
			t.Root.SetChildren(tomatoes)

			for i, tomato := range tomatoes {
				tomato.AddChild(ux.NewNode(BuyTomato{
					Date:           now.Add(time.Duration(i) * time.Hour),
					Name:           "kGe",
					GrossWeight:    int64(i),
					Wagon:          11,
					NetWeight:      12,
					Price:          13,
					Payment:        14,
					Phone:          "01234567890",
					SettlementDate: now.Add(time.Duration(i) * time.Hour),
					Settler:        "asd",
					Settled:        30,
					Unsettled:      40,
					Tags:           make([]string, 0),
					Score:          1.0,
					Enable:         false,
					Note:           "this is a note",
				}))
				if i%2 == 0 {
					tomato.AddChild(ux.NewNode(BuyTomato{
						Date:           now.Add(time.Duration(i) * time.Hour),
						Name:           "ddkwork",
						GrossWeight:    int64(i),
						Wagon:          11,
						NetWeight:      12,
						Price:          13,
						Payment:        14,
						Phone:          "18188888888",
						SettlementDate: now.Add(time.Duration(i) * time.Hour),
						Settler:        "asd",
						Settled:        30,
						Unsettled:      40,
						Tags:           make([]string, 0),
						Score:          1.0,
						Enable:         false,
						Note:           "this is a note",
					}))
				}
			}
		},
		JsonName:   "BuyTomato",
		IsDocument: false,
	}
	return t
}
