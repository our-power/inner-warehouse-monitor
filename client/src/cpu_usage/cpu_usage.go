package cpu_usage

import (
	"github.com/bitly/go-nsq"
	"github.com/influxdb/influxdb-go"
	"strconv"
	"strings"
	"util"
	"regexp"
	"log"
	"os"
)

type CPUUsageHandler struct {
	db_client  *influxdb.Client
	table_name string
	logger *log.Logger
}

var column_names = []string{"time", "date", "time_index", "ip", "host_name", "hardware_addr", "usage"}

func (h *CPUUsageHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
		实现队列消息处理功能
	*/
	regexPattern := `\A\d{8}\r\n\d{1,4}\r\n(?:(?:25[0-5]|2[0-4]\d|[01]?\d?\d)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d?\d)\r\n.+\r\n[:alnum:]{2}(:[:alnum:]{2})5\r\n\d{2}/\d{2}/\d{4}\s\d{2}:\d{2}:\d{2}\.\d{1,},.+\z`
	regex := regexp.MustCompile(regexPattern)
	matched := regex.MatchString(string(m.Body))
	if matched == false {
		h.logger.Println("***************************************************")
		h.logger.Println(m.Body)
		h.logger.Println("###################################################")
		return
	}

	bodyParts := strings.Split(string(m.Body), "\r\n")

	time_index, err := strconv.Atoi(bodyParts[1])
	time_int := util.FormatTime(bodyParts[0], time_index)
	ps := make([][]interface{}, 0, 1)
	ps = append(ps, []interface{}{time_int, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], strings.Split(bodyParts[5], ",")[1]})
	cpu_msg := influxdb.Series{
		Name:    h.table_name,
		Columns: column_names,
		Points:  ps,
	}

	err = h.db_client.WriteSeries([]*influxdb.Series{&cpu_msg})
	return err
}

func NewCPUUsageHandler(client *influxdb.Client) (cpuUsageHandler *CPUUsageHandler, err error) {
	fh, err := os.OpenFile("cpu_usage.log", os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	logger:= log.New(fh, "cpu_usage_logger", log.LstdFlags)

	cpuUsageHandler = &CPUUsageHandler{
		db_client:  client,
		table_name: "cpu_usage",
		logger: logger,
	}
	return cpuUsageHandler, err
}
