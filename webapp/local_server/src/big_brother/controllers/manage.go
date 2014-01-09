package controllers

import (
	"big_brother/models"
	"github.com/astaxie/beego"
)

type ManageController struct {
	beego.Controller
}

func (this *ManageController) GetManagePage() {
	var machineList []*models.Register
	o.Using("default")
	_, err := o.QueryTable("register").Limit(-1).All(&machineList, "ip", "host_name", "hardware_addr", "agent_version", "machine_role", "status")
	if err != nil {
		this.Data["machine_list"] = nil
	} else {
		this.Data["machine_list"] = machineList
	}
	this.Data["role_mapper"] = map[string]string{
		"ercifenjian": "二次分拣",
		"kaipiao":     "财务开票",
		"dabao":       "打包",
		"fenbo":       "分拨",
	}
	this.Data["status_mapper"] = map[int]string{
		0:  "已正常关机",
		1:  "正常运行中",
		-1: "运行异常",
		-2: "不再使用",
	}
	this.TplNames = "manage_machine.html"
}
