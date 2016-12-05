package stock

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sanjia/goquery"
	"github.com/sanjia/mahonia"
	"sanjia/finance/util"
	"strconv"
	"strings"
)

/*
	从东方财富网上获取资金流向数据，并保存到数据库中
*/
func SaveStockFundFlow(db *sql.DB, stockCode string) {
	url := "http://data.eastmoney.com/zjlx/" + stockCode + ".html"
	doc, _ := goquery.NewDocument(url)
	doc.Find("table.tab1").Find("tbody>tr").Each(func(i int, s1 *goquery.Selection) {
		var data []string
		data = append(data, stockCode)
		s1.Find("td").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			text = strings.Replace(text, "\n", "", -1)
			text = strings.Replace(text, "\t", "", -1)
			enc := mahonia.NewEncoder("gbk")
			str := enc.ConvertString("万")
			if strings.Contains(text, str) {
				text = strings.Replace(text, str, "", -1)
				value, _ := strconv.Atoi(text)
				//value = value * 10000
				//fmt.Printf("%d",value)
				text = strconv.Itoa(value)
			}
			str = enc.ConvertString("亿")
			if strings.Contains(text, str) {
				text = strings.Replace(text, str, "", -1)
				value, _ := strconv.ParseFloat(text, 64)
				value = value * 10000
				text = strconv.FormatFloat(value, 'f', 0, 64)
			}
			if strings.Contains(text, "%") {
				text = strings.Replace(text, "%", "", -1)
				value, _ := strconv.ParseFloat(text, 64)
				value = value / 100
				text = strconv.FormatFloat(value, 'f', 4, 64)
			}
			data = append(data, text)
		})
		stmt, err := db.Prepare(`INSERT fund_flow (stockCode,riqi,shoupan,zhangdiefu,
				zhuli_jinge,zhuli_jingzhanbi,
				chaodadan_jinge,chaodadan_jingzhanbi,
				dadan_jinge,dadan_jingzhanbi,
				zhongdan_jinge,zhongdan_jingzhanbi,
				xiaodan_jinge,xiaodan_jingzhanbi) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
		if err != nil {
			//fmt.Printf("find exception:%s\r\n", err.Error())
			return
		}
		if len(data) <= 13 {
			//fmt.Printf("data invalid\r\n")
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(stockCode, data[1],
			util.ToFloat(data[2]),
			util.ToFloat(data[3]),
			util.ToFloat(data[4]),
			util.ToFloat(data[5]),
			util.ToFloat(data[6]),
			util.ToFloat(data[7]),
			util.ToFloat(data[8]),
			util.ToFloat(data[9]),
			util.ToFloat(data[10]),
			util.ToFloat(data[11]),
			util.ToFloat(data[12]),
			util.ToFloat(data[13]))
		if err != nil {
			//fmt.Printf("find exception:%s\r\n", err.Error())
			return
		}
	})
}
