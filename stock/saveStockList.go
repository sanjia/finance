package stock
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sanjia/goquery"
	"github.com/sanjia/mahonia"
	"database/sql"
	"fmt"
	"strings"
)

//保存深市沪市股票列表到数据库中
func SaveStockList(db *sql.DB)  {
	url:="http://quote.eastmoney.com/stocklist.html"
	doc, err := goquery.NewDocument(url)
	if( err!=nil ) {
		fmt.Printf("find exception:%s\r\n",err.Error())
		return ;
	}

	stmt, err := db.Prepare(`INSERT stock_code_list (code,name) values (?,?)`)
	if( err!=nil ) {
		fmt.Printf("find exception:%s\r\n",err.Error())
		return ;
	}

	doc.Find(".quotebody").Find("ul").Find("a").Each(func(i int, s *goquery.Selection) {
		decode := mahonia.NewDecoder("gbk")
		str := decode.ConvertString(s.Text())
		name:= str
		data := strings.Split(name,"(")

		stockName := data[0]
		if( len(data)>1){
			stockCode := strings.Split(data[1],")")[0]
			//保留A股、B股、创业版股票
			if( strings.HasPrefix(stockCode,"6")|| strings.HasPrefix(stockCode,"0") || strings.HasPrefix(stockCode,"3") ){
				_,err= stmt.Exec(stockCode,stockName)
				if( err!=nil ) {
					fmt.Printf("find exception:%s\r\n",err.Error())
					return ;
				}
			}
		}
	})
}

