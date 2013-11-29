package util

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DbLink struct {
	Today string
	Changing bool
	Links map[string]*sql.DB
}

func NewDbLink(date string) (link *DbLink) {
	links := make(map[string]*sql.DB)
	link = &DbLink{
		date,
		false,
		links,
	}
	return
}

func (link *DbLink) GetLink(date string, hardware_addr string, indicator string) (dbLink *sql.DB, err error) {
	key := date + "_" + hardware_addr
	dbLink, ok := link.Links[key]
	// 如果已经存在，则认为没有日期变更，且数据库连接已经打开
	if ok {
		// fmt.Println("bingo!") //命中已经打开的数据库连接
		return dbLink, nil
	} else {
		// 否则为新的日期打开新的数据库连接，并延时关闭原有日期对应的数据库连接，且删除其在本结构体中的注册条目
		var dbPath, dbSourceName string
		dbPath = "../db/" + date + "/" + strings.Replace(hardware_addr, ":", "_", -1) + "/"
		dbSourceName = dbPath + indicator + ".db"
		os.MkdirAll(dbPath, 0666)
		link.Changing = true
		inComingDate, _ := strconv.Atoi(date)
		currentDate, _ := strconv.Atoi(link.Today)
		if inComingDate > currentDate {
			// 仅当后来的日期比保存的日期更晚时，更新结构体中的Today值
			link.Today = date
		}
		newLink, err := sql.Open("sqlite3", dbSourceName)
		link.Links[key] = newLink
		if err != nil {
			return nil, err
		}
		createTable(indicator, newLink)
		go func() {
			c := time.Tick(5 * time.Minute)
			for _ = range c {
				for k, v := range link.Links {
					// 如果缓存中有非Today的日期，表示已经过期，可以执行延时关闭
					if !strings.HasPrefix(k, link.Today) {
						v.Close()
						fmt.Println(k, "to be deleted")
						delete(link.Links, k)
						link.Changing = false
					}
				}
			}
		}()
		return newLink, nil
	}
	return nil, nil
}

func createTable(indicator string, link *sql.DB) {
	var sql string
	switch indicator {
	case "cpu_usage":
		sql = `
		CREATE TABLE IF NOT EXISTS cpu_usage (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT, usage REAL);
		DELETE FROM cpu_usage;
		`

	case "mem_usage":
		sql = `
		CREATE TABLE IF NOT EXISTS mem_usage (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT, usage REAL);
		DELETE FROM mem_usage;
		`

	case "net_flow":
		sql = `
		CREATE TABLE IF NOT EXISTS net_flow (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT, out_bytes INTEGER, in_bytes INTEGER, out_packets INTEGER, in_packets  INTEGER);
		DELETE FROM net_flow;
		`

	case "heartbeat":
		sql = `
		CREATE TABLE IF NOT EXISTS heartbeat (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT, alive INTEGER NOT NULL);
		DELETE FROM heartbeat;
		`

	case "register":
		sql = `
		CREATE TABLE IF NOT EXISTS register (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT UNIQUE, agent_version TEXT, machine_role TEXT, status INTEGER);
		DELETE FROM register;
		`

	case "accessibility":
		sql = `
		CREATE TABLE IF NOT EXISTS ping_accessibility (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT, target_ip TEXT, response_time TEXT);
		DELETE FROM ping_accessibility;
		`
		sql0 := `
		CREATE TABLE IF NOT EXISTS telnet_accessibility (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT, target_url TEXT, status TEXT);
		DELETE FROM telnet_accessibility;
		`
		_, err := link.Exec(sql0)
		if err != nil {
			fmt.Println(err)
		}
	default:
		return
	}
	_, err := link.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
}
