package manpiecework

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/ux"
	"iter"
	"strconv"
	"time"
)

func New() ux.Widget {
	t := ux.NewTreeTable(ManPieceworkEditData{})
	t.TableContext = ux.TableContext[ManPieceworkEditData]{
		CustomContextMenuItems: func(gtx layout.Context, n *ux.Node[ManPieceworkEditData]) iter.Seq[ux.ContextMenuItem] {
			return func(yield func(ux.ContextMenuItem) bool) {
				yield(ux.ContextMenuItem{
					Title: "xxx",
					Icon:  nil,
					Can: func() bool {
						return true
					},
					Do: func() {

					},
					AppendDivider: false,
					Clickable:     widget.Clickable{},
				})
			}
		},
		MarshalRowCells: func(n *ux.Node[ManPieceworkEditData]) (cells []ux.CellData) {
			timeFmt := n.Data.Time.Format("2006-01-02")
			NumberOfCompleteVehiclesSum := n.Data.NumberOfCompleteVehiclesSum()
			FoamBoxSum := n.Data.FoamBoxSum()
			PlasticFrameSum := n.Data.PlasticFrameSum()
			Mark1 := n.Data.Mark1
			Mark1Bar := n.Data.Mark1Bar
			Mark2 := n.Data.Mark2
			Mark3 := n.Data.Mark3
			Level3 := n.Data.Level3
			Supermarket3Labels := n.Data.Supermarket3Labels
			SmallFruit := n.Data.SmallFruit
			BigFruit := n.Data.BigFruit
			FlowersFruits := n.Data.FlowersFruits
			SmallFruitMark1 := n.Data.SmallFruitMark1

			Mark1_J := n.Data.Mark1_J
			Mark2_J := n.Data.Mark2_J
			Mark3_J := n.Data.Mark3_J
			SmallBunchFruits := n.Data.SmallBunchFruits
			BunchFruits := n.Data.BunchFruits
			FineBunchFruits := n.Data.FineBunchFruits
			SweetBunchFruits := n.Data.SweetBunchFruits
			BunchFruitsAA := n.Data.BunchFruitsAA
			AAA := n.Data.AAA
			TurnoverBasket := n.Data.TurnoverBasket
			ElectronicCommerce := n.Data.ElectronicCommerce
			Carpooling := n.Data.Carpooling

			if n.Container() {
				timeFmt = n.SumChildren()
				for _, n := range n.Walk() {
					FoamBoxSum += n.Data.FoamBoxSum()
					PlasticFrameSum += n.Data.PlasticFrameSum()
					NumberOfCompleteVehiclesSum += n.Data.NumberOfCompleteVehiclesSum()
					Mark1 += n.Data.Mark1
					Mark1Bar += n.Data.Mark1Bar
					Mark2 += n.Data.Mark2
					Mark3 += n.Data.Mark3
					Level3 += n.Data.Level3
					Supermarket3Labels += n.Data.Supermarket3Labels
					SmallFruit += n.Data.SmallFruit
					BigFruit += n.Data.BigFruit
					FlowersFruits += n.Data.FlowersFruits
					SmallFruitMark1 += n.Data.SmallFruitMark1

					Mark1_J += n.Data.Mark1_J
					Mark2_J += n.Data.Mark2_J
					Mark3_J += n.Data.Mark3_J
					SmallBunchFruits += n.Data.SmallBunchFruits
					BunchFruits += n.Data.BunchFruits
					FineBunchFruits += n.Data.FineBunchFruits
					SweetBunchFruits += n.Data.SweetBunchFruits
					BunchFruitsAA += n.Data.BunchFruitsAA
					AAA += n.Data.AAA
					TurnoverBasket += n.Data.TurnoverBasket
					ElectronicCommerce += n.Data.ElectronicCommerce
					Carpooling += n.Data.Carpooling
				}
			}

			//todo check it
			//return []ux.CellData{
			//	{Text: timeFmt},
			//	{Text: n.Data.Destination},
			//	{Text: fmt.Sprint(NumberOfCompleteVehiclesSum), FgColor: unison.Red},
			//	{Text: fmt.Sprint(FoamBoxSum), FgColor: unison.Pink},
			//	{Text: fmt.Sprint(PlasticFrameSum), FgColor: unison.Green},
			//	{Text: fmt.Sprint(Mark1)},
			//	{Text: fmt.Sprint(Mark1Bar)},
			//	{Text: fmt.Sprint(Mark2)},
			//	{Text: fmt.Sprint(Mark3)},
			//	{Text: fmt.Sprint(Level3)},
			//	{Text: fmt.Sprint(Supermarket3Labels)},
			//	{Text: fmt.Sprint(SmallFruit)},
			//	{Text: fmt.Sprint(BigFruit)},
			//	{Text: fmt.Sprint(FlowersFruits)},
			//
			//	{Text: fmt.Sprint(SmallFruitMark1), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(Mark1_J), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(Mark2_J), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(Mark3_J), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(SmallBunchFruits), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(BunchFruits), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(FineBunchFruits), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(SweetBunchFruits), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(BunchFruitsAA), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(AAA), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(TurnoverBasket), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(ElectronicCommerce), FgColor: unison.YellowGreen},
			//	{Text: fmt.Sprint(Carpooling), FgColor: unison.Tomato},
			//	{Text: n.Data.Note},
			//}

			return ux.MarshalRow(n.Data, func(key string, field any) (value string) {
				switch key {
				case "日期":
					return timeFmt
				case "总件数":
					return strconv.Itoa(NumberOfCompleteVehiclesSum) //todo Red
				case "泡沫箱":
					return strconv.Itoa(FoamBoxSum) //todo Pink
				case "胶框":
					return strconv.Itoa(PlasticFrameSum) //todo Green
				case "1标":
					return strconv.Itoa(Mark1)
				case "2<UNK>":
					return "2标"
				case "3标":
					return strconv.Itoa(Mark3)
				case "3层":
					return strconv.Itoa(Level3)
				case "超市3标":
					return strconv.Itoa(Supermarket3Labels)
				case "小果":
					return strconv.Itoa(SmallFruit)
				case "大果":
					return strconv.Itoa(BigFruit)
				case "花果":
					return strconv.Itoa(FlowersFruits)
				case "<UNK>":
					return "小果一标(胶)"
				case "小串":
					return strconv.Itoa(SmallBunchFruits)
				case "串果":
					return strconv.Itoa(BunchFruits)
				case "精串":
					return strconv.Itoa(FineBunchFruits)
				case "甜串":
					return strconv.Itoa(SweetBunchFruits)
				case "串AA":
					return strconv.Itoa(BunchFruitsAA)
				case "AAA":
					return strconv.Itoa(AAA)
				case "周转筐":
					return strconv.Itoa(TurnoverBasket)
				case "电商":
					return strconv.Itoa(ElectronicCommerce)
				case "拼车":
					return strconv.Itoa(Carpooling)
				case "备注":
					return n.Data.Note
				default:
					return ""
				}
			})
		},
		UnmarshalRowCells: func(n *ux.Node[ManPieceworkEditData], rows []ux.CellData) {
			n.Data = ux.UnmarshalRow[ManPieceworkEditData](rows, func(key, value string) (field any) {
				switch key { //todo
				//case "Size":
				//	return int64(datasize.Parse(value))
				//case "ModTime":
				//	return mylog.Check2(time.Parse(time.RFC3339, value))
				default:
					return nil //todo change reflect.Zero(field.Type()) != "" and panic it
				}
			})
		},
		RowSelectedCallback: func() {

		},
		RowDoubleClickCallback: func() {

		},
		SetRootRowsCallBack: func() {
			all := ux.NewContainerNode("2022-02-05", ManPieceworkEditData{
				Time: mylog.Check2(time.Parse("2006-01-02", "2022-02-05")),
			})
			t.Root.AddChild(all)
			fnWageDay := func(date string, wuhan, chengdu, erMei ManPieceworkEditData) {
				containerWageDay := ux.NewContainerNode(date, ManPieceworkEditData{
					Time: mylog.Check2(time.Parse("2006-01-02", date)),
				})
				all.AddChild(containerWageDay)
				wuhan.Destination = "武汉"
				chengdu.Destination = "成都"
				erMei.Destination = "峨眉"
				wuhan.Time = containerWageDay.Data.Time
				chengdu.Time = containerWageDay.Data.Time
				erMei.Time = containerWageDay.Data.Time

				if wuhan.NumberOfCompleteVehiclesSum() > 0 {
					containerWageDay.AddChildByData(wuhan)
				}
				if chengdu.NumberOfCompleteVehiclesSum() > 0 {
					containerWageDay.AddChildByData(chengdu)
				}
				if erMei.NumberOfCompleteVehiclesSum() > 0 {
					containerWageDay.AddChildByData(erMei)
				}
				if containerWageDay.LenChildren() == 0 {
					containerWageDay.Data.Note = "没发车"
				}
				if containerWageDay.LenChildren() == 1 {
					data := containerWageDay.Children[0].Data
					containerWageDay.Remove()
					all.AddChildByData(data)
				}
			}
			fnWageDay("2022-02-05", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    0,
				Level3:                   0,
				Supermarket3Labels:       497,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-06", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-07", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    503,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-08", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    0,
				Level3:                   503,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    0,
				Level3:                   0,
				Supermarket3Labels:       289,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  80,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-09", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    503,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-10", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    503,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-11", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    503,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    0,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          53,
				Mark1_J:                  80,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-12", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    0,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  41,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         22,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-13", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-14", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    0,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  78,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         22,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-15", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-16", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    45,
				Level3:                   0,
				Supermarket3Labels:       170,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  108,
				Mark2_J:                  20,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-17", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-18", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    296,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  84,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-19", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-20", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    247,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          26,
				Mark1_J:                  114,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-21", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    154,
				Level3:                   0,
				Supermarket3Labels:       67,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            49,
				SmallFruitMark1:          15,
				Mark1_J:                  102,
				Mark2_J:                  0,
				Mark3_J:                  2,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-22", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    0,
				Level3:                   507,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-23", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    40,
				Level3:                   0,
				Supermarket3Labels:       207,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            30,
				SmallFruitMark1:          35,
				Mark1_J:                  71,
				Mark2_J:                  0,
				Mark3_J:                  4,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-24", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    45,
				Level3:                   0,
				Supermarket3Labels:       226,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            20,
				SmallFruitMark1:          12,
				Mark1_J:                  80,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-25", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    56,
				Level3:                   0,
				Supermarket3Labels:       145,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            37,
				SmallFruitMark1:          59,
				Mark1_J:                  96,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-26", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-02-27", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    22,
				Level3:                   0,
				Supermarket3Labels:       125,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            16,
				SmallFruitMark1:          65,
				Mark1_J:                  87,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         42,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-02-28", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				SmallFruitMark1:          0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-01", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    5,
				Level3:                   0,
				Supermarket3Labels:       194,
				SmallFruit:               20,
				BigFruit:                 0,
				FlowersFruits:            55,
				SmallFruitMark1:          4,
				Mark1_J:                  48,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallBunchFruits:         6,
				BunchFruits:              65,
				FineBunchFruits:          26,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-03-02", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    1,
				Level3:                   0,
				Supermarket3Labels:       244,
				SmallFruit:               38,
				BigFruit:                 0,
				FlowersFruits:            22,
				Mark1_J:                  73,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallFruitMark1:          16,
				SmallBunchFruits:         0,
				BunchFruits:              2,
				FineBunchFruits:          17,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-03-03", ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    507,
				Level3:                   0,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            0,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallFruitMark1:          0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    18,
				Level3:                   0,
				Supermarket3Labels:       178,
				SmallFruit:               49,
				BigFruit:                 0,
				FlowersFruits:            62,
				Mark1_J:                  71,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallFruitMark1:          15,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{})
			fnWageDay("2022-03-04", ManPieceworkEditData{}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    0,
				Mark1Bar:                 0,
				Mark2:                    0,
				Mark3:                    11,
				Level3:                   0,
				Supermarket3Labels:       262,
				SmallFruit:               47,
				BigFruit:                 0,
				FlowersFruits:            18,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallFruitMark1:          35,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			}, ManPieceworkEditData{
				Time:                     time.Time{},
				Destination:              "",
				NumberOfCompleteVehicles: 0,
				FoamBox:                  0,
				PlasticFrame:             0,
				Mark1:                    13,
				Mark1Bar:                 0,
				Mark2:                    52,
				Mark3:                    162,
				Level3:                   18,
				Supermarket3Labels:       0,
				SmallFruit:               0,
				BigFruit:                 0,
				FlowersFruits:            25,
				Mark1_J:                  0,
				Mark2_J:                  0,
				Mark3_J:                  0,
				SmallFruitMark1:          0,
				SmallBunchFruits:         0,
				BunchFruits:              0,
				FineBunchFruits:          0,
				SweetBunchFruits:         0,
				BunchFruitsAA:            0,
				AAA:                      0,
				TurnoverBasket:           0,
				ElectronicCommerce:       0,
				Carpooling:               0,
				Note:                     "",
			})
			fnWageDay("2022-03-05", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-06", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-07", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-08", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-09", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-10", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-11", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-12", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-13", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-14", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-15", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-16", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-17", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-18", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-19", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-20", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-21", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-22", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-23", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-24", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-25", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-26", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-27", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-28", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-29", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-30", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-03-31", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-04-01", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-04-02", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-04-03", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-04-04", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})
			fnWageDay("2022-04-05", ManPieceworkEditData{}, ManPieceworkEditData{}, ManPieceworkEditData{})

			休息天数 := 0
			总车数 := 0
			总天数 := all.LenChildren()
			for _, n := range all.Walk() {
				if n.Data.Note == "没发车" {
					休息天数++
				}
				if n.Data.NumberOfCompleteVehiclesSum() > 0 {
					总车数++
				}
			}
			info := fmt.Sprintf("\t休息天数：%d\t总车数：%d\t总天数：%d", 休息天数, 总车数, 总天数)
			println(info)

		},
		JsonName:   "Piecework tomato loading",
		IsDocument: false,
	}
	return t

	//ux.RunTest("计件表", t)

}

type (
	BossInfo struct { // 收货地点的提示内容
		Name  string
		Phone string
		Place string // 地点
	}
	DriverInfo struct { // 驾驶员信息,日期的提示内容
		Name        string
		Phone       string
		destination string // 目的地
	}
	// 日结：root容器节点=FoamBox容器节点+PlasticFrame容器节点
	ManPieceworkEditData struct { // 男工计件,用于统计装车
		Time        time.Time `table:"日期"` // 装车日期
		Destination string    `table:"目的地"`

		NumberOfCompleteVehicles int `table:"总件数"`
		FoamBox                  int `table:"泡沫箱"`
		PlasticFrame             int `table:"胶框"`

		Mark1              int `table:"1标"`
		Mark1Bar           int `table:"1标带靶"`
		Mark2              int `table:"2标"`
		Mark3              int `table:"3标"`
		Level3             int `table:"三层"`
		Supermarket3Labels int `table:"超市3标"`
		SmallFruit         int `table:"小果"`
		BigFruit           int `table:"大果"`
		FlowersFruits      int `table:"花果"`

		Mark1_J            int `table:"1标(胶)"`
		Mark2_J            int `table:"2标(胶)"`
		Mark3_J            int `table:"3标(胶)"`
		SmallFruitMark1    int `table:"小果一标(胶)"`
		SmallBunchFruits   int `table:"小串"`
		BunchFruits        int `table:"串果"`
		FineBunchFruits    int `table:"精串"`
		SweetBunchFruits   int `table:"甜串"`
		BunchFruitsAA      int `table:"串AA"`
		AAA                int `table:"AAA"`
		TurnoverBasket     int `table:"周转筐"`
		ElectronicCommerce int `table:"电商"`

		Carpooling int    `table:"拼车"`
		Note       string `table:"备注"`

		// 公用日结结构体
		// 男工报表显示装车费用日结
		// 女工报表显示打包费用日结
		// todo 打包费用没想好怎么显示好一点
		// 打包和装车日结
		// DailySettlement       int       // 日结
		// SettlementDate        time.Date // 结算日期
		// NumberOfCheckoutItems int       // 结账件数,加减人会改变
		// Price                 float32   // 单价
		// TotalPeople           int       // 人数,加减人会改变,根据离职日期计算天数？
		// PerCapitaWage         float64   // 人均工资
		// Settled               float64   // 已结算
		// Unsettled             float64   // 未结算
		// Tags                  []string  // 过滤条件列表
		// Enable                bool      // 是否离职
		// Note                  string    // 备注
	}
	// Salary 底薪
	Cycle struct {
		StartDate   time.Time // 开始日期
		EndDate     time.Time // 结束日期
		Days        int       // 天数
		TotalPeople int       // 人数
	}
)

func (m ManPieceworkEditData) FoamBoxSum() int {
	return m.Mark1 + m.Mark1Bar + m.Mark2 + m.Mark3 + m.Level3 + m.Supermarket3Labels + m.SmallFruit + m.BigFruit + m.FlowersFruits
}

func (m ManPieceworkEditData) PlasticFrameSum() int {
	return m.Mark1_J + m.Mark2_J + m.Mark3_J + m.SmallFruitMark1 +
		m.SmallBunchFruits + m.BunchFruits + m.FineBunchFruits + m.SweetBunchFruits + m.BunchFruitsAA + m.AAA + m.TurnoverBasket + m.ElectronicCommerce
}

func (m ManPieceworkEditData) NumberOfCompleteVehiclesSum() int {
	return m.FoamBoxSum() + m.PlasticFrameSum()
}
