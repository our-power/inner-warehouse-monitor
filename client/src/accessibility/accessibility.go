package accessibility

import (
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/influxdb/influxdb-go"
	"strconv"
	"strings"
	"time"
	"util"
)

/*
 *	将可达性数据(ping)存入数据库的类
 */

type AccessibilityToDBHandler struct {
	db_client         *influxdb.Client
	ping_table_name   string
	telnet_table_name string
	exception_handler *util.ExceptionHandler
}

var ping_column_names = []string{"time", "date", "time_index", "ip", "host_name", "hardware_addr", "target_ip", "response_time"}
var telnet_column_names = []string{"time", "date", "time_index", "ip", "host_name", "hardware_addr", "target_url", "status"}

func (h *AccessibilityToDBHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
		实现队列消息处理功能
	*/
	//fmt.Printf("%s\n", m.Body)

	defer h.exception_handler.HandleException(string(m.Body))

	bodyParts := strings.Split(string(m.Body), "\r\n")

	time_index, err := strconv.Atoi(bodyParts[1])
	time_int := util.FormatTime(bodyParts[0], time_index)

	if bodyParts[5] == "1" {
		validData := bodyParts[6 : len(bodyParts)-1]
		series := make([][]interface{}, 0, 10)
		for _, item := range validData {
			targetAndPingResult := strings.Split(item, ",")
			pingResult := strings.Split(targetAndPingResult[1], "=")
			responseTime := -1
			if pingResult[0] == "ResponseTime" {
				responseTime, err = strconv.Atoi(pingResult[1])
			}
			series = append(series, []interface{}{time_int, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], targetAndPingResult[0], responseTime})
		}
		ping_msg := influxdb.Series{
			Name:    h.ping_table_name,
			Columns: ping_column_names,
			Points:  series,
		}
		err = h.db_client.WriteSeries([]*influxdb.Series{&ping_msg})
	}
	if bodyParts[5] == "2" {
		validData := bodyParts[6 : len(bodyParts)-1]
		series := make([][]interface{}, 0, 10)
		for _, item := range validData {
			targetAndTelnetResult := strings.Split(item, ",")
			status := targetAndTelnetResult[1]
			if status == "" {
				status = "NotOK"
			}
			series = append(series, []interface{}{time_int, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], targetAndTelnetResult[0], status})
		}
		telnet_msg := influxdb.Series{
			Name:    h.telnet_table_name,
			Columns: telnet_column_names,
			Points:  series,
		}
		err = h.db_client.WriteSeries([]*influxdb.Series{&telnet_msg})
	}
	return err
}

func NewAccessibilityToDBHandler(client *influxdb.Client) (accessibilityToDBHandler *AccessibilityToDBHandler, err error) {
	accessibilityToDBHandler = &AccessibilityToDBHandler{
		db_client:         client,
		ping_table_name:   "ping_accessibility",
		telnet_table_name: "telnet_accessibility",
		exception_handler: util.InitHandler("/var/log/accessibility.log", "accessibility_logger"),
	}
	return accessibilityToDBHandler, err
}

/*
 *	检测可达性是否异常的类
 */

type AccessibilityCheckHandler struct{
	exception_handler *util.ExceptionHandler
}

func (h *AccessibilityCheckHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
		实现队列消息处理功能
	*/
	//fmt.Printf("%s\n", m.Body)

	defer h.exception_handler.HandleException(string(m.Body))

	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])

	secondsToNow := time_index * 30
	hour := secondsToNow / 3600
	minutes := secondsToNow % 3600 / 60
	seconds := secondsToNow % 60

	yearTime, _ := time.Parse("20060102", bodyParts[0])
	date := time.Date(yearTime.Year(), yearTime.Month(), yearTime.Day(), hour, minutes, seconds, 0, time.Local).Format("2006-01-02 15:04:05")

	if bodyParts[5] == "1" {
		validData := bodyParts[6 : len(bodyParts)-1]
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
		validData := bodyParts[6 : len(bodyParts)-1]
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

func NewAccessibilityCheckHandler() (accessibilityCheckHandler *AccessibilityCheckHandler, err error) {
	accessibilityCheckHandler = &AccessibilityCheckHandler{
		exception_handler: util.InitHandler("/var/log/check_accessibility.log", "check_accessibility_logger"),
	}
	return accessibilityCheckHandler, err
}
