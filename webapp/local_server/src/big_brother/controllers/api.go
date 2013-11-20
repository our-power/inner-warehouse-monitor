package controllers

import (
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"big_brother/models"
)

type ApiController struct {
	beego.Controller
}

/*
GET /api/serverindicator?hardwareaddr=xxxx&indicator=xxx&date=xxxx
*/
func (this *ApiController) GetServerIndicator() {
	hardware_addr := this.GetString("hardwareaddr")
	indicator := this.GetString("indicator")
	date, _ := time.Parse("2006-01-02", this.GetString("date"))
	date_str := date.Format("20060102")

	var maps []orm.Params
	if indicator == "cpu_usage" {
		o.Using("cpu_usage")
		_, err := o.QueryTable(indicator).Filter("hardware_addr", hardware_addr).Filter("date", date_str).OrderBy("time_index").Values(&maps, "id", "time_index", "ip", "host_name", "usage")
		if err == nil {
			this.Data["json"] = &maps
		} else {
			this.Data["json"] = nil
		}
	}else if indicator == "mem_usage" {
		o.Using("mem_usage")
		_, err := o.QueryTable(indicator).Filter("hardware_addr", hardware_addr).Filter("date", date_str).OrderBy("time_index").Values(&maps, "id", "time_index", "ip", "host_name", "usage")
		if err == nil {
			this.Data["json"] = &maps
		} else {
			this.Data["json"] = nil
		}
	}else if indicator == "net_flow" {
		o.Using("net_flow")
		_, err := o.QueryTable(indicator).Filter("hardware_addr", hardware_addr).Filter("date", date_str).OrderBy("time_index").Values(&maps, "id", "time_index", "ip", "host_name", "out_bytes", "in_bytes", "out_packets", "in_packets")
		if err == nil {
			this.Data["json"] = &maps
		}else {
			this.Data["json"] = nil
		}
	}else {
		this.Data["json"] = nil
	}

	this.ServeJson()
}

/*
GET /api/serverinuse
*/
func (this *ApiController) GetServerList() {
	var serverList []*models.Register
	// 在ORM的数据库default中
	//o.Using("register")
	_, err := o.QueryTable("register").All(&serverList)
	if err == nil {
		this.Data["json"] = &serverList
	}else {
		this.Data["json"] = nil
	}
	this.ServeJson()
}
