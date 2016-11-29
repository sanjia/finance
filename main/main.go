package main

import _ "github.com/go-sql-driver/mysql"
import (
	"database/sql"
	"sanjia/stock"
)

func main() {
	var db *sql.DB
	db, _ = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	defer db.Close()
	stockCode := "6000000"
	////保存深市沪市股票列表
	stock.SaveStockList(db)
	//保存股票摘要数据
	stock.SaveStockAbstractData(db,stockCode)
	//保存股票资金流向数据
	stock.SaveStockFundFlow(db,stockCode)
}
