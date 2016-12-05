package stock

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sanjia/mahonia"
	"io/ioutil"
	"net/http"
	"sanjia/finance/util"
	"strings"
)

//获取股票摘要数据:市盈率、市净率、流通市值、总市值
func GetStockQuote(stockCode string) []string {
	prefix := "sh"
	if strings.HasPrefix(stockCode, "6") == false {
		prefix = "sz"
	}
	url := "http://qt.gtimg.cn/q=" + prefix + stockCode
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("find exception:%s\r\n", err.Error())
	}
	//v_sz000858="51~五 粮 液~000858~35.47~35.17~35.60~341218~159976~181242~35.46~1607~35.45~579~35.44~68~35.43~135~35.42~149~35.47~3~35.48~215~35.49~284~35.50~733~35.51~111~15:00:03/35.47/3772/S/13379284/15158|14:57:00/35.46/51/B/180846/15064|14:56:57/35.46/43/S/152478/15060|14:56:54/35.46/12/S/42552/15057|14:56:51/35.47/143/B/507111/15053|14:56:48/35.46/165/S/585111/15050~20161128150133~0.30~0.85~36.17~35.45~35.46/337446/1204718191~341218~121810~0.90~20.17~~36.17~35.45~2.05~1346.36~1346.43~2.97~38.69~31.65~1.48";
	fmt.Println(string(body))
	str := string(body)
	//将gbk转化为utf8
	decode := mahonia.NewDecoder("gbk")
	str = decode.ConvertString(str)
	data := strings.Split(str, "~")
	/*
		0: 未知
		 1: 名字 2: 代码 3: 当前价格 4: 昨收 5: 今开 6: 成交量（手） 7: 外盘 8: 内盘 9: 买一 10: 买一量（手） 11-18: 买二 买五
		19: 卖一 20: 卖一量 21-28: 卖二 卖五 29: 最近逐笔成交 30: 时间 31: 涨跌 32: 涨跌% 33: 最高 34: 最低 35: 价格/成交量（手）/成交额
		36: 成交量（手）37: 成交额（万）38: 换手率 39: 市盈率
		40:
		41: 最高 42: 最低 43: 振幅 44: 流通市值 45: 总市值 46: 市净率 47: 涨停价 48: 跌停价
	*/
	return data
}

func SaveStockAbstractData(db *sql.DB, stockCode string) {
	//沪市股票
	data := GetStockQuote(stockCode)
	stmt, err := db.Prepare("update stock_code_list set shiyinglv=?,liutongshizhi=?,zongshizhi=?,shijinglv=? where code=?")
	if err != nil {
		fmt.Printf("find exception:%s\r\n", err.Error())
		return
	}
	_, err = stmt.Exec(util.ToFloat(data[39]), util.ToFloat(data[44]), util.ToFloat(data[45]), util.ToFloat(data[46]), data[2])
	if err != nil {
		fmt.Printf("find exception:%s\r\n", err.Error())
		return
	}
}
