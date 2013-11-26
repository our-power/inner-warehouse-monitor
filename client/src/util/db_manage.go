package util

import (
	"os"
	"fmt"
	"time"
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
	var dbPath, dbSourceName string
	dbPath = "D:/" + date + "/" + strings.Replace(hardware_addr, ":", "_", -1) + "/"
	dbSourceName = dbPath + indicator + ".db"
	os.MkdirAll(dbPath, 0666)
	key := date + "_" + hardware_addr
	dbLink, ok := link.Links[key]
	// 如果已经存在，则认为没有日期变更，且数据库连接已经打开
	if ok {
		fmt.Println("bingo!") //命中已经打开的数据库连接
		return dbLink, nil
	} else {
		// 否则为新的日期打开新的数据库连接，并延时关闭原有日期对应的数据库连接，且删除其在本结构体中的注册条目
		link.Changing = true
		link.Today = date
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
					if k != key {
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

	case "ping_accessibility":
		sql = `
		CREATE TABLE IF NOT EXISTS ping_accessibility (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT, target_ip TEXT, response_time TEXT);
		DELETE FROM ping_accessibility;
		`

	case "telnet_accessibility":
		sql = `
		CREATE TABLE IF NOT EXISTS telnet_accessibility (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, date TEXT, time_index INTEGER, ip TEXT, host_name TEXT, hardware_addr TEXT, target_url TEXT, status TEXT);
		DELETE FROM telnet_accessibility;
		`
	default:
		return
	}
	_, err := link.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
}
