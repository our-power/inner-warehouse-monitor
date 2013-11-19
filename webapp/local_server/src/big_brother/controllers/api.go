package controllers

import (
	"time"
	"github.com/astaxie/beego"

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

	if indicator == "cpu_usage" {
		var cpuUsageData []*models.Cpu_usage
		o.QueryTable(indicator).Filter("hardware_addr", hardware_addr).Filter("date", date_str).All(&cpuUsageData)
		this.Data["json"] = &cpuUsageData
	}else if indicator == "mem_usage" {
		var memUsageData []*models.Mem_usage
		o.QueryTable(indicator).Filter("hardware_addr", hardware_addr).Filter("date", date_str).All(&memUsageData)
		this.Data["json"] = &memUsageData
	}else if indicator == "net_flow" {
		var netflowUsageData []*models.Net_flow
		o.QueryTable(indicator).Filter("hardware_addr", hardware_addr).Filter("date", date_str).All(&netflowUsageData)
		this.Data["json"] = &netflowUsageData
	}else{
		this.Data["json"] = nil
	}

	this.ServeJson()
}

/*
GET /api/serverinuse
*/
func (this *ApiController) GetServerListInUse() {

}
