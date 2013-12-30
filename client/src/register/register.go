package register

import (
	"database/sql"
	"github.com/bitly/go-nsq"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
	"util"
)

type RegisterToDBHandler struct {
	db *sql.DB
	exception_handler *util.ExceptionHandler
}

func (h *RegisterToDBHandler) tryHandleIt(m *nsq.Message)(err error){
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
		return err
	} else {
		status = 1
		version_role := strings.Split(bodyParts[5], ",")
		sql := `
		REPLACE INTO register (date, time_index, ip, host_name, hardware_addr, agent_version, machine_role, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
		`
		_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], version_role[0], version_role[1], status)
		return err
	}
	return
}

func (h *RegisterToDBHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
		实现队列消息处理功能
	*/

	defer h.exception_handler.HandleException(string(m.Body))

	err = h.tryHandleIt(m)

	return err
}

func NewRegisterToDBHandler(dbLink *sql.DB) (registerToDBHandler *RegisterToDBHandler, err error) {
	registerToDBHandler = &RegisterToDBHandler{
		db: dbLink,
		exception_handler: util.InitHandler("/var/log/register.log", "register_logger"),
	}
	return registerToDBHandler, err
}
