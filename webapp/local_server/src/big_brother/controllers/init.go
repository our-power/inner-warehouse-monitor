package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/influxdb/influxdb-go"
)

var o orm.Ormer
var client *influxdb.Client

func InitControllers() {
	o = orm.NewOrm()
	config := influxdb.ClientConfig{
		Host:     beego.AppConfig.String("influxdb_host"),
		Username: beego.AppConfig.String("influxdb_user"),
		Password: beego.AppConfig.String("influxdb_passwd"),
		Database: beego.AppConfig.String("influxdb_database"),
	}
	client, err := influxdb.NewClient(&config)
}
