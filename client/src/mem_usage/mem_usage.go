package mem_usage

import (
	"database/sql"
	//"fmt"
	"strings"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
	"util"
	"path"
)

type MemUsageHandler struct {
	db *sql.DB
}

func (h *MemUsageHandler) tryHandleIt(m *nsq.Message) (err error) {
	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])
	sql := `
	INSERT INTO mem_usage (date, time_index, ip, host_name, hardware_addr, usage) VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], strings.Split(bodyParts[5], ",")[1])
	return err
}

func (h *MemUsageHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/
	defer util.HandleException(path.Join(util.LogRoot, "mem_usage.log"), string(m.Body))
	err = h.tryHandleIt(m)
	return err
}

func NewMemUsageHandler(dbLink *sql.DB) (memUsageHandler *MemUsageHandler, err error) {
	memUsageHandler = &MemUsageHandler {
		db: dbLink,
	}
	return memUsageHandler, err
}
