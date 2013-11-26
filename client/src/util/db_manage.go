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

	dbLink, ok := link.Links[date]
	// 如果已经存在，则认为没有日期变更，且数据库连接已经打开
	if ok {
		return dbLink, nil
	} else {
		// 否则为新的日期打开新的数据库连接，并延时关闭原有日期对应的数据库连接，且删除其在本结构体中的注册条目
		link.Changing = true
		link.Today = date
		newLink, err := sql.Open("sqlite3", dbSourceName)
		fmt.Println(err)
		link.Links[date] = newLink
		if err != nil {
			return nil, err
		}
		go func() {
			c := time.Tick(5 * time.Minute)
			for _ = range c {
				for k, v := range link.Links {
					if k != date {
						v.Close()
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

