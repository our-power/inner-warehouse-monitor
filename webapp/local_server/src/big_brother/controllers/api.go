package controllers

import (
	"time"
	"strings"
	"strconv"

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
					Data	  []NcDataType
				}
				results := make([]ResultType,0, 50)
				for _, machine := range maps {
					var rows []orm.Params;
					num, err := o.QueryTable(dataTable).Filter("hardware_addr", machine["Hardware_addr"]).Filter("date", dateStr).OrderBy("time_index").Limit(-1).Values(&rows, "time_index", "out_bytes", "in_bytes", "out_packets", "in_packets")
					if err == nil && num > 0 {
						dataContainerLength := int(rows[num - 1]["Time_index"].(int64)) + 1
						ncNum := len(strings.Split(rows[num-1]["Out_bytes"].(string), ","))
						ncData := make([]NcDataType, ncNum)
						for index := 0; index < ncNum; index++ {
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
							ncData[index] = NcDataType{
								Out_bytes: outBytes,
								In_bytes: inBytes,
								Out_packets: outPackets,
								In_packets: inPackets,
							}
						}
						for _, row := range rows {
							ncsOutByte := strings.Split(row["Out_bytes"].(string), ",")
							ncsInByte := strings.Split(row["In_bytes"].(string), ",")
							ncsOutPacket := strings.Split(row["Out_packets"].(string), ",")
							ncsInPacket := strings.Split(row["In_packets"].(string), ",")
							for i := 0; i < ncNum; i++ {
								ob, _ := strconv.ParseFloat(ncsOutByte[i], 32)
								ncData[i].Out_bytes[int(row["Time_index"].(int64))] = int(ob)
								ib, _ := strconv.ParseFloat(ncsInByte[i], 32)
								ncData[i].In_bytes[int(row["Time_index"].(int64))] = int(ib)
								op, _ := strconv.ParseFloat(ncsOutPacket[i], 32)
								ncData[i].Out_packets[int(row["Time_index"].(int64))] = int(op)
								ip, _ := strconv.ParseFloat(ncsInPacket[i], 32)
								ncData[i].In_packets[int(row["Time_index"].(int64))] = int(ip)
							}
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
		num, err := o.QueryTable("cpu_usage").Filter("hardware_addr", hardwareAddr).Filter("date", dateStr).OrderBy("time_index").All(&cpuUsageData, "time_index", "usage")
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
		num, err := o.QueryTable("mem_usage").Filter("hardware_addr", hardwareAddr).Filter("date", dateStr).OrderBy("time_index").All(&memUsageData, "time_index", "usage")
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
		num, err := o.QueryTable("net_flow").Filter("hardware_addr", hardwareAddr).Filter("date", dateStr).OrderBy("time_index").All(&netFlowData, "time_index", "out_bytes", "in_bytes", "out_packets", "in_packets")
		if err == nil && num > 0 {
			dataContainerLength := netFlowData[num - 1].Time_index + 1
			networkCardNum := len(strings.Split(netFlowData[num-1].Out_bytes, ","))
			type ResultType struct {
				Out_bytes   []int
				In_bytes    []int
				Out_packets []int
				In_packets  []int
			}
			results := make([]ResultType, networkCardNum)
			for index := 0; index < networkCardNum; index++ {
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
				results[index].Out_bytes = outBytes
				results[index].In_bytes = inBytes
				results[index].Out_packets = outPackets
				results[index].In_packets = inPackets
			}
			for _, row := range netFlowData {
				ncsOutBytes := strings.Split(row.Out_bytes, ",")
				ncsInBytes := strings.Split(row.In_bytes, ",")
				ncsOutPackets := strings.Split(row.Out_packets, ",")
				ncsInPackets := strings.Split(row.In_packets, ",")
				for index := 0; index < networkCardNum; index++ {
					outBytesFloat, _ := strconv.ParseFloat(ncsOutBytes[index], 32)
					results[index].Out_bytes[row.Time_index] = int(outBytesFloat)
					inBytesFloat, _ := strconv.ParseFloat(ncsInBytes[index], 32)
					results[index].In_bytes[row.Time_index] = int(inBytesFloat)
					outPacketsFloat, _ := strconv.ParseFloat(ncsOutPackets[index], 32)
					results[index].Out_packets[row.Time_index] = int(outPacketsFloat)
					inPacketsFloat, _ := strconv.ParseFloat(ncsInPackets[index], 32)
					results[index].In_packets[row.Time_index] = int(inPacketsFloat)
				}
			}
			this.Data["json"] = results
		}else {
			this.Data["json"] = nil
		}
	}
	this.ServeJson()
}
