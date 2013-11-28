package controllers

import (
	//"fmt"
	"time"

	"runtime"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"big_brother/models"
)

type ApiController struct {
	beego.Controller
}

/*
GET /api/get_machine_list
*/
func (this *ApiController) GetMachineList() {
	var serverList []*models.Register
	// 在ORM的数据库default中
	o.Using("default")
	_, err := o.QueryTable("register").All(&serverList)
	if err == nil {
		this.Data["json"] = &serverList
	}else {
		this.Data["json"] = nil
	}
	this.ServeJson()
}

/*
	GET /api/status_overview?role=xxx
*/
func (this * ApiController) GetStatusOverview() {
	machineRole := this.GetString("role")
	var serverList []orm.Params
	var available, shutdown, exception int
	var one, zero int64
	one = 1
	zero = 0
	o.Using("default")
	_, err := o.QueryTable("register").Filter("machine_role", machineRole).Limit(-1).Values(&serverList, "ip", "host_name", "hardware_addr", "status")
	if err != nil {
		this.Data["json"] = nil
	} else {
		for _, item := range serverList {
			switch item["Status"] {
			case one:
				available++
			case zero:
				shutdown++
			default:
				exception++
			}
		}
		statistics := []interface {}{available, shutdown, exception, &serverList}
		this.Data["json"] = statistics
	}
	this.ServeJson()
}

// GET  /api/get_step_indicator_data?step=xxx&date=xxx&indicator=xxx
func (this *ApiController) GetStepIndicatorData() {

	// 主动触发垃圾回收，但还没测过效果
	runtime.GC()

	step := this.GetString("step")
	queryDate := this.GetString("date")
	indicator := this.GetString("indicator")

	var role string
	switch step {
	case "kaipiao":
		role = "kaipiao"
	case "ercifenjian":
		role = "ercifenjian"
	case "dabao":
		role = "dabao"
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

	if role == "" || dataTable == "" {
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
					Data	  []float64
				}
				results := make([]ResultType,0, 50)
				for _, machine := range maps {
					var rows []orm.Params
					num, err := o.QueryTable(dataTable).Filter("hardware_addr", machine["Hardware_addr"]).Filter("date", dateStr).OrderBy("time_index").Limit(-1).Values(&rows, "time_index", "usage")
					if err == nil && num > 0 {
						dataContainerLength := int(rows[num - 1]["Time_index"].(int64)) + 1
						usageData := make([]float64, dataContainerLength)
						for index := 0; index < dataContainerLength; index++ {
							usageData[index] = -1;
						}
						for _, row := range rows {
							time_index, _ := row["Time_index"].(int64)
							usageData[time_index] = row["Usage"].(float64)
						}

						host_name, _ := machine["Host_name"].(string)
						results = append(results, ResultType{Host_name: host_name, Data: usageData})
					}
				}
				this.Data["json"] = results
			}else if dataTable == "net_flow" {
				type NcDataType struct {
					Out_bytes   []int
					In_bytes    []int
					Out_packets []int
					In_packets  []int
				}
				type ResultType struct {
					Host_name string
					Data	  NcDataType
				}
				results := make([]ResultType,0, 50)
				for _, machine := range maps {
					var rows []orm.Params;
					num, err := o.QueryTable(dataTable).Filter("hardware_addr", machine["Hardware_addr"]).Filter("date", dateStr).OrderBy("time_index").Limit(-1).Values(&rows, "time_index", "out_bytes", "in_bytes", "out_packets", "in_packets")
					if err == nil && num > 0 {
						dataContainerLength := int(rows[num - 1]["Time_index"].(int64)) + 1

						outBytes := make([]int, dataContainerLength)
						inBytes := make([]int, dataContainerLength)
						outPackets := make([]int, dataContainerLength)
						inPackets := make([]int, dataContainerLength)
						for index := 0; index < dataContainerLength; index++ {
							outBytes[index] = -1
							inBytes[index] = -1
							outPackets[index] = -1
							inPackets[index] = -1
						}
						ncData := NcDataType{
							Out_bytes: outBytes,
							In_bytes: inBytes,
							Out_packets: outPackets,
							In_packets: inPackets,
						}
						for _, row := range rows {
							ncData.Out_bytes[int(row["Time_index"].(int64))] = int(row["Out_bytes"].(int64))
							ncData.In_bytes[int(row["Time_index"].(int64))] = int(row["In_bytes"].(int64))
							ncData.Out_packets[int(row["Time_index"].(int64))] = int(row["Out_packets"].(int64))
							ncData.In_packets[int(row["Time_index"].(int64))] = int(row["In_packets"].(int64))
						}
						results = append(results, ResultType{Host_name: machine["Host_name"].(string), Data: ncData})
					}
				}
				this.Data["json"] = results
			}
		}else {
			this.Data["json"] = nil
		}
	}
	this.ServeJson()
}

// 获取某一天某机器机器 CPU使用率 或 内存使用量 或 网卡数据
// GET /api/get_machine_indicator_data?hardware_addr=xxx&date=xxx&indicator=xxx
func (this *ApiController) GetMachineIndicatorData() {

	runtime.GC()

	hardwareAddr := this.GetString("hardware_addr")
	indicator := this.GetString("indicator")
	queryDate := this.GetString("date")
	date, _ := time.Parse("2006-01-02", queryDate)
	dateStr := date.Format("20060102")

	if indicator == "cpu_usage" {
		o.Using("cpu_usage")
		var cpuUsageData []*models.Cpu_usage
		num, err := o.QueryTable("cpu_usage").Filter("hardware_addr", hardwareAddr).Filter("date", dateStr).OrderBy("time_index").Limit(-1).All(&cpuUsageData, "time_index", "usage")
		if err == nil && num > 0 {
			dataContainerLength := cpuUsageData[num - 1].Time_index + 1
			results := make([]float32, dataContainerLength)
			for index := 0; index < dataContainerLength; index++ {
				results[index] = -1
			}
			for _, row := range cpuUsageData {
				results[row.Time_index] = row.Usage
			}
			this.Data["json"] = results
		}else {
			this.Data["json"] = nil
		}
	}else if indicator == "mem_usage" {
		o.Using("mem_usage")
		var memUsageData []*models.Mem_usage
		num, err := o.QueryTable("mem_usage").Filter("hardware_addr", hardwareAddr).Filter("date", dateStr).OrderBy("time_index").Limit(-1).All(&memUsageData, "time_index", "usage")
		if err == nil && num > 0 {
			dataContainerLength := memUsageData[num - 1].Time_index + 1
			results := make([]float32, dataContainerLength)
			for index := 0; index < dataContainerLength; index++ {
				results[index] = -1
			}
			for _, row := range memUsageData {
				results[row.Time_index] = row.Usage
			}
			this.Data["json"] = results
		}else {
			this.Data["json"] = nil
		}
	}else {
		o.Using("net_flow")
		var netFlowData []*models.Net_flow
		num, err := o.QueryTable("net_flow").Filter("hardware_addr", hardwareAddr).Filter("date", dateStr).OrderBy("time_index").Limit(-1).All(&netFlowData, "time_index", "out_bytes", "in_bytes", "out_packets", "in_packets")
		if err == nil && num > 0 {
			dataContainerLength := netFlowData[num - 1].Time_index + 1
			type ResultType struct {
				Out_bytes   []int
				In_bytes    []int
				Out_packets []int
				In_packets  []int
			}
			outBytes := make([]int, dataContainerLength)
			inBytes := make([]int, dataContainerLength)
			outPackets := make([]int, dataContainerLength)
			inPackets := make([]int, dataContainerLength)
			for index := 0; index < dataContainerLength; index++ {
				outBytes[index] = -1
				inBytes[index] = -1
				outPackets[index] = -1
				inPackets[index] = -1
			}
			results := ResultType{Out_bytes: outBytes, In_bytes: inBytes, Out_packets: outPackets, In_packets: inPackets}
			for _, row := range netFlowData {
				results.Out_bytes[row.Time_index] = row.Out_bytes
				results.In_bytes[row.Time_index] = row.In_bytes
				results.Out_packets[row.Time_index] = row.Out_packets
				results.In_packets[row.Time_index] = row.In_packets
			}
			this.Data["json"] = results
		}else {
			this.Data["json"] = nil
		}
	}
	this.ServeJson()
}

// 获取机器最新的可达性数据
func (this *ApiController) GetMachineAccessibilityData() {
	hardwareAddr := this.GetString("hardware_addr")
	now := time.Now()
	dateStr := now.Format("20060102")

	type PingResultType struct {
		Target_ip     string
		Response_time int
	}

	type TelnetResultType struct {
		Target_url string
		Status     string
	}

	type ResultType struct{
		Hardware_addr string
		Date              string
		Ping_time_index   int
		Ping_results      []PingResultType
		Telnet_time_index int
		Telnet_results    []TelnetResultType
	}

	pingResults := make([]PingResultType,0, 100)
	telnetResults := make([]TelnetResultType,0, 100)
	var pingTimeIndex int
	var telnetTimeIndex int
	o.Using("accessibility")
	var pingItems []*models.Ping_accessibility
	// 这里使用的Limit(20)是假设最多有20个服务
	num, err := o.QueryTable("ping_accessibility").Filter("hardware_addr", hardwareAddr).Filter("date", dateStr).OrderBy("-time_index").Limit(100).All(&pingItems)
	if err == nil && num > 0 {
		newestTimeIndex := pingItems[0].Time_index
		pingTimeIndex = newestTimeIndex
		for _, item := range pingItems {
			if item.Time_index < newestTimeIndex {
				break
			}else {
				pingResults = append(pingResults, PingResultType{Target_ip: item.Target_ip, Response_time: item.Response_time})
			}
		}
	}

	var telnetItems []*models.Telnet_accessibility
	num, err = o.QueryTable("telnet_accessibility").Filter("hardware_addr", hardwareAddr).Filter("date", dateStr).OrderBy("-time_index").Limit(100).All(&telnetItems)
	if err == nil && num > 0 {
		newestTimeIndex := telnetItems[0].Time_index
		telnetTimeIndex = newestTimeIndex
		for _, item := range telnetItems {
			if item.Time_index < newestTimeIndex {
				break
			}else {
				telnetResults = append(telnetResults, TelnetResultType{Target_url: item.Target_url, Status: item.Status})
			}
		}
	}
	this.Data["json"] = ResultType{Hardware_addr: hardwareAddr, Date: now.Format("2006-01-02"), Ping_time_index: pingTimeIndex, Ping_results: pingResults, Telnet_time_index: telnetTimeIndex, Telnet_results: telnetResults}
	this.ServeJson()
}
