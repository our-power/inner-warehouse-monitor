package controllers

import (
	"time"

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

	var role string
	switch step {
	case "kaipiao":
		role = "test"
	case "ercifenjian":
		role = "test"
	case "dabao":
		role = "test"
	case "fenbo":
		role = "fenbo"
	default:
		role = ""
	}

	var dataTable string
	switch indicator{
	case "cpu_view":
		dataTable = "cpu_usage"
	case "memory_view":
		dataTable = "mem_usage"
	case "netflow_view":
		dataTable = "net_flow"
	default:
		dataTable = ""
	}

	date, _ := time.Parse("2006-01-02", queryDate)
	dateStr := date.Format("20060102")

	if role == "" || dataTable == ""{
		this.Data["json"] = nil
	}else {
		database := dataTable
		var maps []orm.Params
		o.Using("default")
		_, err := o.QueryTable("register").Filter("machine_role", role).Limit(-1).Values(&maps, "host_name", "hardware_addr")
		if err == nil {
			o.Using(database)
			if dataTable == "cpu_usage" || dataTable == "mem_usage" {
				type ResultType struct {
					Host_name string
					Data	[]float64
				}

				results := make([]ResultType, 0, 50)
				for _, machine := range maps {
					var rows []orm.Params
					num, err := o.QueryTable(dataTable).Filter("hardware_addr", machine["Hardware_addr"]).Filter("date", dateStr).OrderBy("time_index").Limit(-1).Values(&rows, "time_index", "usage")
					usageData := make([]float64, rows[num-1]["Time_index"].(int64) + 1)
					if err == nil {
						for _, row := range rows {
							time_index, _ := row["Time_index"].(int64)
							usageData[time_index] = row["Usage"].(float64)
						}
					}
					host_name, _ := machine["Host_name"].(string)
					results = append(results, ResultType{Host_name: host_name, Data: usageData})
				}
				this.Data["json"] = results
			}else if dataTable == "net_flow" {
				/*
				type UsageDataType struct {
					Time_index int
					Out_bytes string
					In_bytes string
					Out_packets string
					In_packets string
				}
				type ResultType struct {
					Host_name string
					Data	[]UsageDataType
				}
				for _, machine := range maps {
					_, err = o.QueryTable(dataTable).Filter("hardware_addr", machine["hardware_addr"]).Filter("date", dateStr)
				}
				*/
			}

		}else{
			this.Data["json"] = nil
		}
	}
	this.ServeJson()
}


func (this *IndicatorDataController) GetMachineIndicatorsData() {

}
