package cpu_usage

import (
	"strings"
	"strconv"
	"util"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
)

type CPUUsageHandler struct {
	db *util.DbLink
}

func (h *CPUUsageHandler) HandleMessage(m *nsq.Message)(err error){
	/*
	实现队列消息处理功能
	*/

	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])
	db, err := h.db.GetLink(bodyParts[0], bodyParts[4], "cpu_usage")
	if err != nil {
		return err
	}
	sql := `
	INSERT INTO cpu_usage (date, time_index, ip, host_name, hardware_addr, usage) VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err = db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], strings.Split(bodyParts[5],",")[1])

	return err
}

func NewCPUUsageHandler(dbLink *util.DbLink)(cpuUsageHandler *CPUUsageHandler, err error){
	cpuUsageHandler = &CPUUsageHandler {
		db: dbLink,
	}
	return cpuUsageHandler, err
}



