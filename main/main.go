package main

import (
	"fmt"
	"sanjia/finance/qq"
)

func main() {
	fmt.Printf("stat")
	//stock.StatStockFundFlowByTick()
	qq.SaveTickAsExcel("002124")
}
