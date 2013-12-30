package register

import (
	"database/sql"
	"strings"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
	"util"
	"path"
)

type RegisterToDBHandler struct {
	db *sql.DB
}

func (h *RegisterToDBHandler) tryHandleIt(m *nsq.Message)(err error) {
	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])

	/*
	0:机器正常关闭
	1:机器正在运行
	-1:机器没有正常发送心跳数据
	-2:删除机器（该机器不再投入使用）
	*/
	var status int
	if bodyParts[5] == "shutdown" {
		status = 0
		sql := `
		UPDATE register SET date=?, time_index=?, ip=?, host_name=?, status=? WHERE hardware_addr=?;
		`
		_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], status, bodyParts[4])
	} else {
		status = 1
		version_role := strings.Split(bodyParts[5], ",")
		sql := `
		REPLACE INTO register (date, time_index, ip, host_name, hardware_addr, agent_version, machine_role, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
		`
		_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], version_role[0], version_role[1], status)
	}

	return err
}

func (h *RegisterToDBHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/
	defer util.HandleException(path.Join(util.LogRoot, "register.log"), string(m.Body))
	err = h.tryHandleIt(m)
	return err
}

func NewRegisterToDBHandler(dbLink *sql.DB) (registerToDBHandler *RegisterToDBHandler, err error) {
	registerToDBHandler = &RegisterToDBHandler {
		db: dbLink,
	}
	return registerToDBHandler, err
}
