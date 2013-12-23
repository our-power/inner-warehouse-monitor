package mem_usage

import (
	"strings"
	"strconv"
	"util"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
)

type MemUsageHandler struct {
	db *util.DbLink
}

func (h *MemUsageHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/
	bodyParts := strings.Split(string(m.Body), "\r\n")
	if len(bodyParts) == 6 {
		time_index, err := strconv.Atoi(bodyParts[1])
        db, err := h.db.GetLink(bodyParts[0], bodyParts[4], "mem_usage")
        if err != nil {
            return err
        }
        sql := `
        INSERT INTO mem_usage (date, time_index, ip, host_name, hardware_addr, usage) VALUES (?, ?, ?, ?, ?, ?);
        `
		_, err = db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], strings.Split(bodyParts[5], ",")[1])
		return err
	}
	return nil
}

func NewMemUsageHandler(dbLink *util.DbLink)(memUsageHandler *MemUsageHandler, err error) {
	memUsageHandler = &MemUsageHandler {
		db: dbLink,
	}
	return memUsageHandler, err
}
