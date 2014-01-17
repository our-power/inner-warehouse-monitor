package controllers

import (
	//"fmt"
	"big_brother/models"
	"strings"
)

type SearchController struct {
	BaseController
}

func (this *SearchController) GetSearchPage() {
	this.Data["login_name"] = this.GetSession("login_name")
	if this.GetSession("role_type") == "admin_user" {
		this.Data["admin"] = true
	}
	this.TplNames = "search_machine.html"
}

type ResultType struct {
	IsExisted     bool
	SearchItem    string
	Ip            string
	Host_name     string
	Hardware_addr string
	Machine_role  string
	Status        string
}

func queryMachine(col string, items string) (partialResults []ResultType) {
	o.Using("default")
	if items != "" {
		item_s := strings.Split(items, ",")
		itemNum := len(item_s)
		for index := 0; index < itemNum; index++ {
			var rows []*models.Register
			num, err := o.QueryTable("register").Filter(col, item_s[index]).Limit(-1).All(&rows, "ip", "host_name", "hardware_addr", "machine_role", "status")
			if err == nil {
				if num == 0 {
					partialResults = append(partialResults, ResultType{IsExisted: false, SearchItem: item_s[index], Ip: "", Host_name: "", Hardware_addr: "", Machine_role: "", Status: ""})
				}
				if num > 0 {
					for _, row := range rows {
						var machineRole string
						switch row.Machine_role {
						case "kaipiao":
							machineRole = "财务开票"
						case "ercifenjian":
							machineRole = "二次分拣"
						case "dabao":
							machineRole = "打包"
						case "fenbo":
							machineRole = "分拨"
						default:
							machineRole = "其它"
						}
						var status string
						switch row.Status {
						case 0:
							status = "已正常关机"
						case 1:
							status = "正常运行中"
						case -1:
							status = "运行异常"
						case -2:
							status = "不再使用"
						default:
							status = "未知状态"
						}
						partialResults = append(partialResults, ResultType{IsExisted: true, SearchItem: item_s[index], Ip: row.Ip, Host_name: row.Host_name, Hardware_addr: row.Hardware_addr, Machine_role: machineRole, Status: status})
					}
				}
			}
		}
	}
	return
}

func (this *SearchController) FilterMachineList() {
	ipList := this.GetString("iplist")
	hostNameList := this.GetString("hostnamelist")
	hardwareAddrList := this.GetString("hardwareaddrlist")
	results := make([]ResultType, 0, 10)
	results = append(results, queryMachine("ip", ipList)...)
	results = append(results, queryMachine("host_name", hostNameList)...)
	results = append(results, queryMachine("hardware_addr", hardwareAddrList)...)

	this.Data["json"] = results
	this.ServeJson()
}

func (this *SearchController)GetIndicatorsByMac() {
	mac := this.GetString("mac")
	if mac == "" {
		this.Abort("404")
	}
	this.Data["targets"] = mac
	this.Data["login_name"] = this.GetSession("login_name")
	if this.GetSession("role_type") == "admin_user" {
		this.Data["admin"] = true
	}
	this.TplNames = "search_machine.html"
}
