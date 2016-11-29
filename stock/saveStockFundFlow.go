package stock
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/goquery"
	"github.com/mahonia"
	"database/sql"
	"fmt"
	"strings"
	"strconv"
)

/*
	保存股票资金流向数据
*/
func SaveStockFundFlow(db *sql.DB,stockCode string)  {
	url := "http://data.eastmoney.com/zjlx/" + stockCode + ".html";
	doc, _ := goquery.NewDocument(url)
	doc.Find("table.tab1").Find("tbody>tr").Each(func(i int, s1 *goquery.Selection) {
		var data []string
		data = append(data, "601688")
		s1.Find("td").Each(func(i int, s *goquery.Selection) {
			text := s.Text()
			text = strings.Replace(text, "\n", "", -1)
			text = strings.Replace(text, "\t", "", -1)
			enc := mahonia.NewEncoder("gbk")
			str := enc.ConvertString("万")
			if ( strings.Contains(text, str) ) {
				text = strings.Replace(text, str, "", -1)
				value, _ := strconv.Atoi(text)
				//value = value * 10000
				//fmt.Printf("%d",value)
				text = strconv.Itoa(value)
			}
			str = enc.ConvertString("亿")
			if ( strings.Contains(text, str) ) {
				text = strings.Replace(text, str, "", -1)
				value, _ := strconv.ParseFloat(text, 64)
				value = value * 10000
				text = strconv.FormatFloat(value, 'f', 0, 64)
			}
			if ( strings.Contains(text, "%") ) {
				text = strings.Replace(text, "%", "", -1)
				value, _ := strconv.ParseFloat(text, 64)
				value = value / 100
				text = strconv.FormatFloat(value, 'f', 4, 64)
			}
			data = append(data, text)
		});
		stmt, err := db.Prepare(`INSERT fund_flow (stockCode,riqi,shoupan,zhangdiefu,
				zhuli_jinge,zhuli_jingzhanbi,
				chaodadan_jinge,chaodadan_jingzhanbi,
				dadan_jinge,dadan_jingzhanbi,
				zhongdan_jinge,zhongdan_jingzhanbi,
				xiaodan_jinge,xiaodan_jingzhanbi) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
		if( err!=nil ){
			fmt.Printf("find exception:%s\r\n",err.Error())
			return ;
		}

		_,err= stmt.Exec("601688",data[1],
			F(data[2]),
			F(data[3]),
			F(data[4]),
			F(data[5]),
			F(data[6]),
			F(data[7]),
			F(data[8]),
			F(data[9]),
			F(data[10]),
			F(data[11]),
			F(data[12]),
			F(data[13]))
		if( err!=nil ) {
			fmt.Printf("find exception:%s\r\n",err.Error())
			return ;
		}
	});
}