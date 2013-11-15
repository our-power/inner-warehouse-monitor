package accessibility

import (
	"database/sql"
	"fmt"
	//"strings"
	//"strconv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/bitly/go-nsq"
)

/*
 *	将可达性数据(ping)存入数据库的类
 */

type AccessibilityToDBHandler struct {
	db *sql.DB
}

func (h *AccessibilityToDBHandler) HandleMessage(m *nsq.Message)(err error){
	/*
	实现队列消息处理功能
	*/
	fmt.Printf("%s\n", m.Body)
	/*
	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])
	sql := `
	INSERT INTO cpu_usage (date, time_index, ip, host_name, hardware_addr, usage) VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], strings.Split(bodyParts[5],",")[1])
	*/
	return err
}

func NewAccessibilityToDBHandler(dbLink *sql.DB)(accessibilityToDBHandler *AccessibilityToDBHandler, err error){
	accessibilityToDBHandler = &AccessibilityToDBHandler {
		db: dbLink,
	}
	return accessibilityToDBHandler, err
}


/*
 *	检测可达性是否异常的类
 */

type AccessibilityCheckHandler struct {
	db *sql.DB
}

func (h *AccessibilityCheckHandler) HandleMessage(m *nsq.Message)(err error){
	/*
	实现队列消息处理功能
	*/
	fmt.Printf("%s\n", m.Body)
	/*
	bodyParts := strings.Split(string(m.Body), "\r\n")
	time_index, err := strconv.Atoi(bodyParts[1])
	sql := `
	INSERT INTO cpu_usage (date, time_index, ip, host_name, hardware_addr, usage) VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err = h.db.Exec(sql, bodyParts[0], time_index, bodyParts[2], bodyParts[3], bodyParts[4], strings.Split(bodyParts[5],",")[1])
	*/
	return err
}

func NewAccessibilityCheckHandler(dbLink *sql.DB)(accessibilityCheckHandler *AccessibilityCheckHandler, err error){
	accessibilityCheckHandler = &AccessibilityCheckHandler {
		db: dbLink,
	}
	return accessibilityCheckHandler, err
}
