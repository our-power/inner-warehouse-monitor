package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type NavItemsController struct {
	beego.Controller
}

func (this *NavItemsController) GetServerDataGroupByStep() {
	step := this.GetString("step")
	var role string
	switch step {
	case "kaipiao":
		role = "test"
	case "ercifenjian":
		role = "test"
	case "dabao":
		role = "test"
	case "fenbo":
		role = "test"
	default:
		role = ""
	}

	if role == "" {
		this.Abort("404")
	}else {
		var maps []orm.Params
		o.Using("default")
		_, err := o.QueryTable("register").Filter("machine_role", role).Values(&maps)
		if err == nil {
			this.Data["nav_now"] = step
			this.Data["server_info_list"] = &maps
		}else{
			this.Data["server_info_list"] = nil
		}
		this.TplNames = "serverDataGroupByStep.html"
	}
}
