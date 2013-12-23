package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var o orm.Ormer
var multidbRoot string

func InitControllers() {
	o = orm.NewOrm()
	multidbRoot = beego.AppConfig.String("multidb")
}
