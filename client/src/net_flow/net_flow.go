package net_flow

import (
	"database/sql"
	//"fmt"
	"strings"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
)

type NetFlowHandler struct {
	db *sql.DB
}

func (h *NetFlowHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/

	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])
	netFlowDataParts := strings.Split(bodyParts[5], ",")[1:]
	networkCardNum := len(netFlowDataParts)/4
	beginIndex := 0
	endIndex := networkCardNum
	outBytes := strings.Join(netFlowDataParts[beginIndex:endIndex], ",")
	beginIndex = endIndex
	endIndex += networkCardNum
	inBytes := strings.Join(netFlowDataParts[beginIndex:endIndex], ",")
	beginIndex = endIndex
	endIndex += networkCardNum
	outPackets := strings.Join(netFlowDataParts[beginIndex:endIndex], ",")
	beginIndex = endIndex
	endIndex += networkCardNum
	inPackets := strings.Join(netFlowDataParts[beginIndex:endIndex], ",")
	sql := `
	INSERT INTO net_flow (date, time_index, ip, host_name, hardware_addr, out_bytes, in_bytes, out_packets, in_packets) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], outBytes, inBytes, outPackets, inPackets)

	return err
}

func NewNetFlowHandler(dbLink *sql.DB) (netFlowHandler *NetFlowHandler, err error) {
	netFlowHandler = &NetFlowHandler {
		db: dbLink,
	}
	return netFlowHandler, err
}




