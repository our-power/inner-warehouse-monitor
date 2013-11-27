package heartbeat

import (
	"database/sql"
	"fmt"
	"strings"
	"strconv"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
)

type HeartBeatHandler struct {
	db *sql.DB
}

func (h *HeartBeatHandler) HandleMessage(m *nsq.Message) (err error) {
	/*
	实现队列消息处理功能
	*/
	//fmt.Printf("%s\n", m.Body)

	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])
	sql := `
	INSERT INTO heartbeat (date, time_index, ip, host_name, hardware_addr, alive) VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], 1)

	return err
}

func NewHeartBeatHandler(dbLink *sql.DB) (heartBeatHandler *HeartBeatHandler, err error) {
	heartBeatHandler = &HeartBeatHandler {
		db: dbLink,
	}
	return heartBeatHandler, err
}

func updateMachineStatus(h *HeartBeatHandler, registerDB *sql.DB) {
	for {
		c := time.Tick(3*time.Minute)
		for _ = range c {
			sql := "SELECT hardware_addr, status FROM register"
			rows, _ := registerDB.Query(sql)
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
					sql = "SELECT count(*) FROM heartbeat WHERE hardware_addr = ? AND date = ? AND time_index > ?"
					rows, _ := h.db.Query(sql, item.Hardware_addr, dateStr, criticalTimeIndex)
					var count int
					for rows.Next() {
						rows.Scan(&count)
						break
					}
					rows.Close()
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
