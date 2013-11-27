package cpu_usage

import (
	"database/sql"
	"strings"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
)

type CPUUsageHandler struct {
	db *sql.DB
}

func (h *CPUUsageHandler) HandleMessage(m *nsq.Message)(err error){
	/*
	实现队列消息处理功能
	*/
	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])
	sql := `
	INSERT INTO cpu_usage (date, time_index, ip, host_name, hardware_addr, usage) VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], strings.Split(bodyParts[5],",")[1])
	return err
}

func NewCPUUsageHandler(dbLink *sql.DB)(cpuUsageHandler *CPUUsageHandler, err error){
	cpuUsageHandler = &CPUUsageHandler {
		db: dbLink,
	}
	return cpuUsageHandler, err
}



