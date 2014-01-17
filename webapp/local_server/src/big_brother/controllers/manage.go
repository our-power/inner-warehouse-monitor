package controllers

import (
	"big_brother/models"
	"strconv"
	"strings"
)

type ManageController struct {
	BaseController
}

func (this *ManageController) GetManagePage() {
	var machineList []*models.Register
	o.Using("default")
	_, err := o.QueryTable("register").Limit(-1).All(&machineList, "id", "ip", "host_name", "hardware_addr", "agent_version", "machine_role", "status")
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

	statusLabelMapper := map[int]string{
		0:  "label",
		1:  "label label-success",
		-1: "label label-important",
		-2: "label label-inverse",
	}
	machineStatusLabels := make(map[string]string)
	for _, machine := range machineList {
		machineStatusLabels[machine.Hardware_addr] = statusLabelMapper[machine.Status]
	}
	this.Data["status_labels"] = machineStatusLabels
	this.Data["login_name"] = this.GetSession("login_name")
	if this.GetSession("role_type") == "admin_user" {
		this.Data["admin"] = true
	}
	del_permission := false
	permissions := strings.Split(this.GetSession("permission").(string), "|")
	for _, value := range permissions {
		if value == "del" {
			del_permission = true
			break
		}
	}
	this.Data["del_permission"] = del_permission
	this.TplNames = "manage_machine.html"
}

func (this *ManageController) DelMachine() {
	del_permission := false
	permissions := strings.Split(this.GetSession("permission").(string), "|")
	for _, value := range permissions {
		if value == "del" {
			del_permission = true
			break
		}
	}
	if del_permission {
		id := this.GetString("id")
		if id == "" {
			this.Data["json"] = map[string]string{
				"Status": "failure",
				"Msg": "未提供作业机器的数据库ID，未能删除机器！",
			}
		}else {
			id_int, _ := strconv.Atoi(id)
			o.Using("default")
			num, err := o.Delete(&models.Register{Id: id_int})
			if err != nil {
				this.Data["json"] = map[string]string{
					"Status": "failure",
					"Msg": "数据库操作出错！",
				}
			}else if num == 0 {
				this.Data["json"] = map[string]string{
					"Status": "failure",
					"Msg": "未能删除该机器，可能数据库中不存在该机器。",
				}
			}else {
				this.Data["json"] = map[string]string{
					"Status": " success",
					"Msg": "成功删除该机器！",
				}
			}
		}
	} else {
		this.Data["json"] = map[string]string{
			"Status": " failure",
			"Msg": "没有删除操作的权限！",
		}
	}
	this.ServeJson()
}
