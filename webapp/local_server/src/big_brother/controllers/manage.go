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
	this.TplNames = "manage_machine.html"
}
