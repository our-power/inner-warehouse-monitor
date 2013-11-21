package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type IndicatorDataController struct {
	beego.Controller
}

func (this *IndicatorDataController) GetStepIndicatorData() {
	step := this.GetString("step")
	queryDate := this.GetString("date")
	indicator := this.GetString("indicator")

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
		this.Data["json"] = nil
	}else {
		var maps []orm.Params
		o.Using("default")
		_, err := o.QueryTable("register").Filter("machine_role", role).Values(&maps, "host_name", "hardware_addr")
		if err == nil {
			for _, machine := range maps {

			}
		}else{
			this.Data["json"] = nil
		}
	}
	this.ServeJson()
}


func (this *IndicatorDataController) GetMachineIndicatorsData() {

}
