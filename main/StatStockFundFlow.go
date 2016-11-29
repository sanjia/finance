package main

import _ "github.com/go-sql-driver/mysql"
import (
	"database/sql"
	"fmt"
	"os"
	"sanjia/finance/stock"
	"sanjia/finance/util"
	"strings"
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
		fmt.Errorf("demo:./StatStockFundFlow 600000|ALL|SH|SZ")
		return
	}
	var stockList []string
	if os.Args[1] == "SH" || os.Args[1] == "SZ" || os.Args[1] == "ALL" {
		rows, err := db.Query("select code from stock_code_list")
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			var stockCode string
			err := rows.Scan(&stockCode)
			if err != nil {
				fmt.Println(err)
			}
			if os.Args[1] == "ALL" {
				stockList = append(stockList, stockCode)
			} else if os.Args[1] == "SH" {
				if strings.HasPrefix(stockCode, "6") == true {
					stockList = append(stockList, stockCode)
				}
			} else if os.Args[1] == "SZ" {
				if strings.HasPrefix(stockCode, "6") == false {
					stockList = append(stockList, stockCode)
				}
			}
		}
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		stockList = append(stockList, os.Args[1])
	}

	//保存股票资金流向数据
	fmt.Printf("start to stat fund flow,stock sum:%d \r\n", len(stockList))
	for i := 0; i < len(stockList); i++ {
		fmt.Printf("process:%s (%d/%d)", stockList[i], i, len(stockList))
		stock.SaveStockFundFlow(db, stockList[i])
	}
}
