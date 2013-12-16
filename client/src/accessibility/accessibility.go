package accessibility

import (
	"database/sql"
	"fmt"
	"strings"
	"strconv"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
)

/*
 *	将可达性数据(ping)存入数据库的类
 */

type AccessibilityToDBHandler struct {
	db *sql.DB
}

func (h *AccessibilityToDBHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/
	//fmt.Printf("%s\n", m.Body)

	bodyParts := strings.Split(string(m.Body), "\r\n")
	if len(bodyParts) >= 6 {
		time_index, err := strconv.Atoi(bodyParts[1])
		if bodyParts[5] == "1" {
			tx, err := h.db.Begin()
			if err != nil {
				fmt.Println(err)
			}
			stmt, err := tx.Prepare("INSERT INTO ping_accessibility (date, time_index, ip, host_name, hardware_addr, target_ip, response_time) VALUES (?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				fmt.Println(err)
			}
			defer stmt.Close()
			validData := bodyParts[6:len(bodyParts) - 1]
			for _, item := range validData {
				targetAndPingResult := strings.Split(item, ",")
				pingResult := strings.Split(targetAndPingResult[1], "=")
				responseTime := -1
				if pingResult[0] == "ResponseTime" {
					responseTime, err = strconv.Atoi(pingResult[1])
				}
				_, err = stmt.Exec(bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], targetAndPingResult[0], responseTime)
			}

			tx.Commit()
		}
		if bodyParts[5] == "2" {

			tx, err := h.db.Begin()
			if err != nil {
				fmt.Println(err)
			}
			stmt, err := tx.Prepare("INSERT INTO telnet_accessibility (date, time_index, ip, host_name, hardware_addr, target_url, status) VALUES (?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				fmt.Println(err)
			}
			defer stmt.Close()
			validData := bodyParts[6:len(bodyParts) - 1]
			for _, item := range validData {
				targetAndTelnetResult := strings.Split(item, ",")
				status := targetAndTelnetResult[1]
				if status == "" {
					status = "NotOK"
				}
				_, err = stmt.Exec(bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], targetAndTelnetResult[0], status)
			}
			tx.Commit()
		}
		return err
	}
}

func NewAccessibilityToDBHandler(dbLink *sql.DB) (accessibilityToDBHandler *AccessibilityToDBHandler, err error) {
	accessibilityToDBHandler = &AccessibilityToDBHandler {
		db: dbLink,
	}
	return accessibilityToDBHandler, err
}

/*
 *	检测可达性是否异常的类
 */

type AccessibilityCheckHandler struct {}

func (h *AccessibilityCheckHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/
	//fmt.Printf("%s\n", m.Body)

	bodyParts := strings.Split(string(m.Body), "\r\n")
	if len(bodyParts) >= 6 {
		time_index, err := strconv.Atoi(bodyParts[1])

		secondsToNow := time_index*30
		hour := secondsToNow/3600
		minutes := secondsToNow%3600/60
		seconds := secondsToNow%60

		yearTime, _ := time.Parse("20060102", bodyParts[0])
		date := time.Date(yearTime.Year(), yearTime.Month(), yearTime.Day(), hour, minutes, seconds, 0, time.Local).Format("2006-01-02 15:04:05")

		if bodyParts[5] == "1" {
			validData := bodyParts[6:len(bodyParts) - 1]
			for _, item := range validData {
				targetAndPingResult := strings.Split(item, ",")
				pingResult := strings.Split(targetAndPingResult[1], "=")
				responseTime := -1
				if pingResult[0] == "ResponseTime" {
					responseTime, err = strconv.Atoi(pingResult[1])
				}
				if responseTime == -1 {
					fmt.Printf("%s  %s 无法ping通 %s\n", date, bodyParts[2], targetAndPingResult[0])
				}
			}

		}
		if bodyParts[5] == "2" {
			validData := bodyParts[6:len(bodyParts) - 1]
			for _, item := range validData {
				targetAndTelnetResult := strings.Split(item, ",")
				status := targetAndTelnetResult[1]
				if status == "" {
					fmt.Printf("%s  %s 无法成功连接到 %s\n", date, bodyParts[2], targetAndTelnetResult[0])
				}
			}
		}
		return err
	}
}

func NewAccessibilityCheckHandler() (accessibilityCheckHandler *AccessibilityCheckHandler, err error) {
	accessibilityCheckHandler = &AccessibilityCheckHandler {}
	return accessibilityCheckHandler, err
}
