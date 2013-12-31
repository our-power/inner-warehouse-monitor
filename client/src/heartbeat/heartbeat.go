package heartbeat

import (
	"strings"
	"strconv"
	"util"
	"time"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
	"path"
)

type HeartBeatHandler struct {
	db *util.DbLink
}

func (h *HeartBeatHandler) tryHandleIt(m *nsq.Message)(err error) {
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

func (h *HeartBeatHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/
	defer util.HandleException(path.Join(util.LogRoot, "heartbeat.log"), string(m.Body))
	err = h.tryHandleIt(m)
	return err
}

func NewHeartBeatHandler(dbLink *util.DbLink) (heartBeatHandler *HeartBeatHandler, err error) {
	heartBeatHandler = &HeartBeatHandler {
		db: dbLink,
	}
	return heartBeatHandler, err
}

func updateMachineStatus(h *HeartBeatHandler, registerDB *sql.DB) {
	for {
		c := time.Tick(1*time.Minute)
		for _ = range c {
			sql := "SELECT hardware_addr, status FROM register"
			rows, err := registerDB.Query(sql)
			if err != nil {
				fmt.Println(err)
			}
			type machineStatus struct {
				Hardware_addr string
				Status        int
			}
			// 这里的参数100应根据不同仓库机器数目进行调整
			machineListWithStatus := make([]machineStatus,0, 100)
			var hardwareAddr string
			var status int
			for rows.Next() {
				rows.Scan(&hardwareAddr, &status)
				machineListWithStatus = append(machineListWithStatus, machineStatus{Hardware_addr: hardwareAddr, Status: status})
			}
			rows.Close()
			now := time.Now()
			dateStr := now.Format("20060102")
			nowTimeIndex := now.Hour()*60*2 + now.Minute()*2 + now.Second()/30
			criticalTimeIndex := nowTimeIndex - 6
			for _, item := range machineListWithStatus {
				if item.Status == 1 || item.Status == -1 {
					db, err := h.db.GetLink(dateStr, item.Hardware_addr, "heartbeat")
					if err != nil {
						fmt.Println(err)
					}
					var count int
					sql = "SELECT count(*) FROM heartbeat WHERE time_index > ?"
					err = db.QueryRow(sql, criticalTimeIndex).Scan(&count)
					if err != nil {
						continue
					}
					sql = "UPDATE register SET date=?, time_index=?, status=? WHERE hardware_addr=?"
					newStatus := 0
					if count == 0 && item.Status == 1 {
						newStatus = -1
					}
					if count > 0 && item.Status == -1 {
						newStatus = 1
					}
					if newStatus == 1 || newStatus == -1 {
						_, err := registerDB.Exec(sql, dateStr, nowTimeIndex, newStatus, item.Hardware_addr)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}
}

func (h *HeartBeatHandler) CheckPeriodically(registerDB *sql.DB) {
	go updateMachineStatus(h, registerDB)
}
