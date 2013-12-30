package cpu_usage

import (
	"github.com/bitly/go-nsq"
	"github.com/influxdb/influxdb-go"
	"strconv"
	"strings"
	"util"
	"os"
	"log"
)

type CPUUsageHandler struct {
	db_client  *influxdb.Client
	exception_handler *util.ExceptionHandler
	table_name string
}

var column_names = []string{"time", "date", "time_index", "ip", "host_name", "hardware_addr", "usage"}

func (h *CPUUsageHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
		实现队列消息处理功能
	*/

	defer h.exception_handler.HandleException(string(m.Body))

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
	fh, _ := os.OpenFile("/var/log/cpu_usage.log", os.O_RDWR | os.O_APPEND | os.O_CREATE, 0777)
	defer file.Close()

	l := log.New(fh, "cpu_usage_logger", LstdFlags)

	eh := &util.ExceptionHandler {
		Logger: l,
	}

	cpuUsageHandler = &CPUUsageHandler{
		db_client:  client,
		exception_handler: eh,
		table_name: "cpu_usage",
	}
	return cpuUsageHandler, err
}
