package main

import _ "github.com/go-sql-driver/mysql"
import (
	"database/sql"
	"fmt"
	"os"
	"sanjia/finance/stock"
	"sanjia/finance/util"
)

func main() {
	cfg := new(util.Config)
	cfg.InitConfig("app.conf")
	connectURI := cfg.Read("mysql", "connectURI")
	var db *sql.DB
	db, _ = sql.Open("mysql", connectURI)
	defer db.Close()

	arg_num := len(os.Args)
	if arg_num != 2 {
		fmt.Errorf("demo:./StatStockFundFlow 600000")
		return
	}

	//保存股票资金流向数据
	stock.SaveStockFundFlow(db, os.Args[1])
}
