package heartbeat

import (
	"strings"
	"strconv"
	"util"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
)

type HeartBeatHandler struct {
	db *util.DbLink
}

func (h *HeartBeatHandler) HandleMessage(m *nsq.Message)(err error){
	/*
	实现队列消息处理功能
	*/

	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])
	db, err := h.db.GetLink(bodyParts[0], bodyParts[4], "heartbeat")
	if err != nil {
		return err
	}
	sql := `
	INSERT INTO heartbeat (date, time_index, ip, host_name, hardware_addr, alive) VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err = db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], 1)

	return err
}

func NewHeartBeatHandler(dbLink *util.DbLink)(heartBeatHandler *HeartBeatHandler, err error){
	heartBeatHandler = &HeartBeatHandler {
		db: dbLink,
	}
	return heartBeatHandler, err
}
