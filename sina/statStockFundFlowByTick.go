package qq

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"sanjia/finance/util"
	"strings"
)

//
func calcTickDetailUrl(stockCode string, pageIndex int) string {
	if strings.HasPrefix(stockCode, "6") == true {
		return fmt.Sprintf("http://stock.gtimg.cn/data/index.php?appn=detail&action=data&c=sh%s&p=%d", stockCode, pageIndex)
	}

	return fmt.Sprintf("http://stock.gtimg.cn/data/index.php?appn=detail&action=data&c=sz%s&p=%d", stockCode, pageIndex)
}

func SaveTickAsExcel(stockCode string) {
	//stockCode := "600666"//os.Args[1];
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("成交明细")
	row := sheet.AddRow()
	var names = [6]string{"时间", "价格", "涨跌", "成交量", "成交额", "成交性质"}
	for x := 0; x < 6; x++ {
		cell := row.AddCell()
		cell.Value = names[x]
	}
	var buy_cdd int64
	var sell_cdd int64
	var buy_dd int64
	var sell_dd int64
	var buy_zd int64
	var sell_zd int64
	var buy_xd int64
	var sell_xd int64

	for n := 1; ; n++ {
		url := calcTickDetailUrl(stockCode, n)
		body, err := util.HttpGet(url)
		if err != nil {
			break
		}

		data := strings.Split(string(body), "\"")
		if len(data) != 3 {
			break
		}
		fmt.Printf("response:%s\r\n", data[1])
		tick := strings.Split(string(data[1]), "|")

		for m := 0; m < len(tick); m++ {
			list := strings.Split(tick[m], "/")
			if len(list) < 7 {
				continue
			}
			row := sheet.AddRow()

			//v := util.ToInt(list[4])
			a := util.ToInt(list[5])

			if list[6] == "B" {
				if a >= 1000000 {
					buy_cdd = buy_cdd + a
				} else if a >= 200000 {
					buy_dd = buy_dd + a
				} else if a >= 50000 {
					buy_zd = buy_zd + a
				} else {
					buy_xd = buy_xd + a
				}
			} else if list[6] == "S" {
				if a >= 1000000 {
					sell_cdd = sell_cdd + a
				} else if a >= 200000 {
					sell_dd = sell_dd + a
				} else if a >= 50000 {
					sell_zd = sell_zd + a
				} else {
					sell_xd = sell_xd + a
				}
			}

			for k := 1; k < 7; k++ {
				cell := row.AddCell()
				cell.Value = list[k]
				if k >= 2 && k <= 5 {
					value := util.ToFloat(list[k])
					cell.SetFloat(value)
					if value < 0 {
						//cell.GetStyle().Font.Size=32
						cell.GetStyle().Font.Color = "FF0000"
					}
				}

			}
		}
	}

	sheet, _ = file.AddSheet("统计数据")
	row = sheet.AddRow()
	cell := row.AddCell()
	cell.Value = "超大单买入"
	cell = row.AddCell()
	cell.Value = "超大单卖出"
	cell = row.AddCell()
	cell.Value = "大单买入"
	cell = row.AddCell()
	cell.Value = "大单卖出"
	cell = row.AddCell()
	cell.Value = "中单买入"
	cell = row.AddCell()
	cell.Value = "中单卖出"
	cell = row.AddCell()
	cell.Value = "小单买入"
	cell = row.AddCell()
	cell.Value = "小单卖出"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetValue(buy_cdd / 10000.0)
	cell = row.AddCell()
	cell.SetValue(sell_cdd / 10000.0)
	cell = row.AddCell()
	cell.SetValue(buy_dd / 10000.0)
	cell = row.AddCell()
	cell.SetValue(sell_dd / 10000.0)
	cell = row.AddCell()
	cell.SetValue(buy_zd / 10000.0)
	cell = row.AddCell()
	cell.SetValue(sell_zd / 10000.0)
	cell = row.AddCell()
	cell.SetValue(buy_xd / 10000.0)
	cell = row.AddCell()
	cell.SetValue(sell_xd / 10000.0)

	err := file.Save(stockCode + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
