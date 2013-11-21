package controllers

import (
	"github.com/astaxie/beego/orm"
)

var o orm.Ormer

func InitControllers () {
	o = orm.NewOrm()
}
