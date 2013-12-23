package mem_usage

import (
	//"fmt"
	"strings"
	"strconv"
	"github.com/bitly/go-nsq"
    "github.com/influxdb/influxdb-go"
    "util"
)

type MemUsageHandler struct {
	db_client *influxdb.Client
    table_name string
}

var column_names = [5]string{"time", "date", "time_index", "ip", "host_name", "hardware_addr", "usage"}

func (h *MemUsageHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/
	bodyParts := strings.Split(string(m.Body), "\r\n")
	if len(bodyParts) == 6 {
		time_index, err := strconv.Atoi(bodyParts[1])
        time_int := util.FormatTime(bodyParts[0], time_index)
        point := Point_Type {
            Time: time_int,
            Ip: bodyParts[2],
            Host_name: bodyParts[3],
            Hardware_addr: bodyParts[4],
            Usage: 
        }
        mem_msg := influxdb.Series{
            Name: h.table_name,
            Columns: column_names,
            Points: [][]interface{}{time_int, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], strings.Split(bodyParts[5], ",")[1]},
        }

        err = h.db_client.WriteSeries([]influxdb.Series{mem_msg})
        return err
	}
	return nil
}

func NewMemUsageHandler(client *influxdb.Client) (memUsageHandler *MemUsageHandler, err error) {
	memUsageHandler = &MemUsageHandler {
		db_client: client,
        table_name: "mem_usage"
	}
	return memUsageHandler, err
}
