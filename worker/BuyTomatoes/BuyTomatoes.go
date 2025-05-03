package BuyTomatoes

import (
	"time"
)

type BuyTomato struct {
	// 列标题加入cell的二级文本显示，内容为列合计，如果是cell类型是可计算类型
	// Place string //地点
	Date           time.Time // 进货日期，每天为一个容器节点，这样就不用像wps那样进货明细需要一堆sheet了
	Name           string    // 卖家姓名
	GrossWeight    int64     // 毛重
	Wagon          int64     // 车皮
	NetWeight      int64     // 净重 sum
	Price          float64   // 单价
	Payment        float64   // 货款 sum//(毛重-车皮)x2(市斤变公斤)x0.9（除皮10%）x单价=货款,因西红柿损耗大（坏果多），结账不再四舍五入
	Phone          string    // 电话
	SettlementDate time.Time // 结算日期
	Settler        string    // 结算人 列入tag过滤查询每个结算人的合计
	Settled        float64   // 已结算 sum
	Unsettled      float64   // 未结算 sum
	Tags           []string  // 过滤条件列表
	Score          float64   // 评分图标，西红柿硬度，色度，卖家性格（价格沟通，损耗除皮结账）
	Enable         bool      // 排除个别卖家 toogle ？？参考gcs表格右键设置状态变灰色，svg交易完成图标图标
	Note           string    // 备注
}

// BuyTomatoesEditData 西红柿进货erp表格设计
type (

	// Salary 底薪
	Cycle struct {
		StartDate   time.Time // 开始日期
		EndDate     time.Time // 结束日期
		Days        int       // 天数
		TotalPeople int       // 人数
	}
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
	ManPieceworkEditData struct { // 男工计件,用于统计装车
		Date                     time.Time // 装车日期
		FoamBox                  int       // 泡沫箱
		PlasticFrame             int       // 胶框
		Mark1                    int       // 1标
		Mark2                    int       // 2标
		Mark3                    int       // 3标
		Level3                   int       // 三层
		Supermarket3Labels       int       // 超市3标
		SmallFruit               int       // 小果
		BigFruit                 int       // 大果
		FlowersFruits            int       // 花果
		SmallFruitMark1          int       // 小果一标
		Mark1Bar                 int       // 1标带靶
		SmallBunchFruits         int       // 小串
		BunchFruits              int       // 串果
		FineBunchFruits          int       // 精串
		SweetBunchFruits         int       // 甜串
		BunchFruitsAA            int       // 串AA
		AAA                      int       // AAA
		TurnoverBasket           int       // 周转筐
		Sum                      int       // 合计
		ElectronicCommerce       int       // 电商
		Carpooling               int       // 拼车
		NumberOfCompleteVehicles int       // 整车件数

		// 公用日结结构体
		// 男工报表显示装车费用日结
		// 女工报表显示打包费用日结
		// todo 打包费用没想好怎么显示好一点
		// 打包和装车日结
		DailySettlement       int       // 日结
		SettlementDate        time.Time // 结算日期
		NumberOfCheckoutItems int       // 结账件数,加减人会改变
		Price                 float32   // 单价
		TotalPeople           int       // 人数,加减人会改变,根据离职日期计算天数？
		PerCapitaWage         float64   // 人均工资
		Settled               float64   // 已结算
		Unsettled             float64   // 未结算
		Tags                  []string  // 过滤条件列表
		Enable                bool      // 是否离职
		Note                  string    // 备注
	}
	ManWagesEditData struct { // 男工工资
	}
	WoManPieceworkEditData struct { // 女工计件
	}
	WoManWagesEditData struct { // 女工工资
	}
)

/*
卖家个数 Number of sellers ,卖家个数显示，deep+非容器节点
车数 Number of cars

转换为容器节点，当天只有一笔交易为非容器节点，多笔交易为容器节点，这样显示更明了和清晰
容器节点的车数统计显示，其实就是孩子的个数

如果是女工计件统计表的话加入行合计，该场景不需要

装车表格模型和这个差不多，改一下结构体完事，加入平均日结工资和周期加减人划分结算，按泡沫箱和胶框分类合计，小计
女工计件加入离职图标和状态，工资结算状态，列标题的二级文本显示工作天数，总件数，单价，工资

男女工的工资结算方式分为：大合唱，加减人，周期结算，多批男女工的话按照打包和装车结算，一批的话按最后一天的余货结算

米易报表的分类和标题合并

周期计件必须写明人数，起始日期
女工计件加入比较底薪

数据引用考虑dock查询或者一个窗口多张表显示
如果加入更多的场景的话，考虑将业务代码文件和数据存储文件放到文件夹，然后加入树形导航

财务类的固定资产和损益表设计

webview？安卓，云同步

因为泛型类型约束的原因，无法通过配置修改列标题，增加业务场景就要修改代码重新编译，是否合理

因为部分平台前置窗口失效的原因不用模态窗口而是非模态的dock，是否合理，滚动问题，编辑文本性能问题，着色和索引问题

进销存+财务系统还是自己设计的好用
*/
