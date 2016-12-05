package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sanjia/finance/util"
	"net/smtp"
	"sanjia/finance/stock"
	"strconv"
	"strings"
)

func SendStockMonitorEmail(body string) {
	user := "wl97@yeah.net"
	password := "hosthost1234"
	host := "smtp.yeah.net:25"
	to := "6844357@qq.com"
	subject := "股价告警邮件"
	err := SendMail(user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("send mail success!")
	}
}
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func main() {

	cfg := new(util.Config)
	cfg.InitConfig("app.conf")
	connectURI := cfg.Read("mysql", "connectURI")
	var db *sql.DB
	db, _ = sql.Open("mysql", connectURI)
	rows, err := db.Query("select stockCode,topWarningPrice,lowWarningPrice from warning_stock;")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	body := `<html>
			<body>
				<h3>"告警邮件列表"</h3>
				<table>%s</table>
			</body>
		</html>`

	content := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", "股票名称", "股票代码", "最低价", "最高价", "当前价格")
	defer rows.Close()
	for rows.Next() {
		var stockCode string
		var topWarningPrice, lowWarningPrice float64
		err := rows.Scan(&stockCode, &topWarningPrice, &lowWarningPrice)
		if err != nil {
			fmt.Println(err)
		}

		var price float64
		data := stock.GetStockQuote(stockCode)
		if len(data) <= 3 {
			continue
		}

		price, _ = strconv.ParseFloat(data[3], 32)
		if price > topWarningPrice || price < lowWarningPrice {
			content = content + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%6.2f</td><td>%6.2f</td><td>%6.2f</td></tr>", data[1], stockCode, lowWarningPrice, topWarningPrice, price)
		}
		fmt.Print(content)
	}

	html := fmt.Sprintf(body, content)
	fmt.Print(html)
	SendStockMonitorEmail(html)
}
