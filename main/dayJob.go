package main

import _ "github.com/go-sql-driver/mysql"
import (
	"database/sql"
	"sanjia/stock"
	"fmt"
)

func main() {
	cfg := new(stock.Config)
	cfg.InitConfig("main/app.conf")
	connectURI := cfg.Read("mysql", "connectURI")
	var db *sql.DB
	db, _ = sql.Open("mysql", connectURI)
	defer db.Close()
	////保存深市沪市股票列表
	stock.SaveStockList(db)
	fmt.Printf("SaveStockList ok")
}
