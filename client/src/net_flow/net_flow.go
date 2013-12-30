package net_flow

import (
	//"fmt"
	"github.com/bitly/go-nsq"
	"github.com/influxdb/influxdb-go"
	"strconv"
	"strings"
	"util"
)

type NetFlowHandler struct {
	db_client  *influxdb.Client
	exception_handler *util.ExceptionHandler
	table_name string
}

var column_names = []string{"time", "date", "time_index", "ip", "host_name", "hardware_addr", "out_bytes", "in_bytes", "out_packets", "in_packets"}

func (h *NetFlowHandler) tryHandleIt(m *nsq.Message) ([][]interface{}) {
	bodyParts := strings.Split(string(m.Body), "\r\n")

	time_index, err := strconv.Atoi(bodyParts[1])
	time_int := util.FormatTime(bodyParts[0], time_index)

	netFlowDataParts := strings.Split(bodyParts[5], ",")[1:]
	networkCardNum := len(netFlowDataParts) / 4

	beginIndex := 0
	endIndex := networkCardNum
	outBytes := 0
	for index := beginIndex; index < endIndex; index++ {
		oneOutByteFloat, _ := strconv.ParseFloat(netFlowDataParts[index], 32)
		outBytes += int(oneOutByteFloat)
	}
	beginIndex = endIndex
	endIndex += networkCardNum
	inBytes := 0
	for index := beginIndex; index < endIndex; index++ {
		oneInByteFloat, _ := strconv.ParseFloat(netFlowDataParts[index], 32)
		inBytes += int(oneInByteFloat)
	}

	beginIndex = endIndex
	endIndex += networkCardNum
	outPackets := 0
	for index := beginIndex; index < endIndex; index++ {
		oneOutPacketFloat, _ := strconv.ParseFloat(netFlowDataParts[index], 32)
		outPackets += int(oneOutPacketFloat)
	}

	beginIndex = endIndex
	endIndex += networkCardNum
	inPackets := 0
	for index := beginIndex; index < endIndex; index++ {
		oneInPacketsFloat, _ := strconv.ParseFloat(netFlowDataParts[index], 32)
		inPackets += int(oneInPacketsFloat)
	}

	ps := make([][]interface{}, 0, 1)
	ps = append(ps, []interface{}{time_int, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], outBytes, inBytes, outPackets, inPackets})
	return ps
}

func (h *NetFlowHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
		实现队列消息处理功能

		按指标叠加所有网卡的数据
	*/

	defer h.exception_handler.HandleException(string(m.Body))

	ps := h.tryHandleIt(m)

	netflow_msg := influxdb.Series{
		Name:    h.table_name,
		Columns: column_names,
		Points:  ps,
	}
	err = h.db_client.WriteSeries([]*influxdb.Series{&netflow_msg})
	return err
}

func NewNetFlowHandler(client *influxdb.Client) (netFlowHandler *NetFlowHandler, err error) {
	netFlowHandler = &NetFlowHandler{
		db_client:  client,
		exception_handler: util.InitHandler("/var/log/net_flow.log", "net_flow_logger"),
		table_name: "net_flow",
	}
	return netFlowHandler, err
}
