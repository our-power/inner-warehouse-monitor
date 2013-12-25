package controllers

import (
	"big_brother/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"runtime"
	"strconv"
	"time"
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
	} else {
		this.Data["json"] = nil
	}
	this.ServeJson()
}

/*
	GET /api/status_overview?role=xxx
*/
func (this *ApiController) GetStatusOverview() {
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
		statistics := []interface{}{available, shutdown, exception, &serverList}
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
	switch indicator {
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
	} else {
		var maps []orm.Params
		o.Using("default")
		_, err := o.QueryTable("register").Filter("machine_role", role).Limit(-1).Values(&maps, "host_name", "hardware_addr")
		if err == nil {
			if dataTable == "cpu_usage" || dataTable == "mem_usage" {
				type ResultType struct {
					Host_name string
					Data      []float32
				}
				results := make([]ResultType, 0, 50)
				for _, machine := range maps {
					query := fmt.Sprintf("SELECT time_index, usage FROM %s WHERE hardware_addr='%s' AND date='%s'", dataTable, machine["Hardware_addr"], dateStr)
					series, _ := client.Query(query)
					if len(series) > 0 && len(series[0].Points) > 0 {

						column_index_mapper := make(map[string]int)
						for index, value := range series[0].Columns {
							column_index_mapper[value] = index
						}
						//By default, InfluxDB returns data in time descending order.
						dataContainerLength := int(series[0].Points[0][column_index_mapper["time_index"]].(float64)) + 1

						usageData := make([]float32, dataContainerLength)
						for index := 0; index < dataContainerLength; index++ {
							usageData[index] = -1
						}
						for _, point := range series[0].Points {
							usage, _ := strconv.ParseFloat(point[column_index_mapper["usage"]].(string), 32)
							time_index := int(point[column_index_mapper["time_index"]].(float64))
							usageData[time_index] = float32(usage)
						}
						host_name, _ := machine["Host_name"].(string)
						results = append(results, ResultType{Host_name: host_name, Data: usageData})
					}
				}
				this.Data["json"] = results
			} else if dataTable == "net_flow" {
				type NcDataType struct {
					Out_bytes   []int
					In_bytes    []int
					Out_packets []int
					In_packets  []int
				}
				type ResultType struct {
					Host_name string
					Data      NcDataType
				}
				results := make([]ResultType, 0, 50)
				for _, machine := range maps {
					query := fmt.Sprintf("SELECT time_index, out_bytes, in_bytes, out_packets, in_packets FROM net_flow WHERE hardware_addr='%s' AND date='%s'", machine["Hardware_addr"], dateStr)
					series, _ := client.Query(query)
					if len(series) > 0 && len(series[0].Points) > 0 {
						column_index_mapper := make(map[string]int)
						for index, value := range series[0].Columns {
							column_index_mapper[value] = index
						}
						dataContainerLength := int(series[0].Points[0][column_index_mapper["time_index"]].(float64)) + 1
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
							Out_bytes:   outBytes,
							In_bytes:    inBytes,
							Out_packets: outPackets,
							In_packets:  inPackets,
						}
						for _, point := range series[0].Points {
							time_index := int(point[column_index_mapper["time_index"]].(float64))
							ncData.Out_bytes[time_index] = int(point[column_index_mapper["out_bytes"]].(int64))
							ncData.In_bytes[time_index] = int(point[column_index_mapper["in_bytes"]].(int64))
							ncData.Out_packets[time_index] = int(point[column_index_mapper["out_packets"]].(int64))
							ncData.In_packets[time_index] = int(point[column_index_mapper["in_packets"]].(int64))
						}
						results = append(results, ResultType{Host_name: machine["Host_name"].(string), Data: ncData})
					}
				}
				this.Data["json"] = results
			}
		} else {
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
		query := fmt.Sprintf("SELECT time_index, usage FROM cpu_usage WHERE hardware_addr='%s' AND date='%s'", hardwareAddr, dateStr)
		series, _ := client.Query(query)
		if len(series) > 0 && len(series[0].Points) > 0 {
			column_index_mapper := make(map[string]int)
			for index, value := range series[0].Columns {
				column_index_mapper[value] = index
			}
			dataContainerLength := int(series[0].Points[0][column_index_mapper["time_index"]].(float64)) + 1
			results := make([]float32, dataContainerLength)
			for index := 0; index < dataContainerLength; index++ {
				results[index] = -1
			}
			for _, point := range series[0].Points {
				usage, _ := strconv.ParseFloat(point[column_index_mapper["usage"]].(string), 32)
				results[int(point[column_index_mapper["time_index"]].(float64))] = float32(usage)
			}
			this.Data["json"] = results
		} else {
			this.Data["json"] = nil
		}
	} else if indicator == "mem_usage" {
		query := fmt.Sprintf("SELECT time_index, usage FROM mem_usage WHERE hardware_addr='%s' AND date='%s'", hardwareAddr, dateStr)
		series, _ := client.Query(query)
		if len(series) > 0 && len(series[0].Points) > 0 {
			column_index_mapper := make(map[string]int)
			for index, value := range series[0].Columns {
				column_index_mapper[value] = index
			}
			dataContainerLength := int(series[0].Points[0][column_index_mapper["time_index"]].(float64)) + 1
			results := make([]float32, dataContainerLength)
			for index := 0; index < dataContainerLength; index++ {
				results[index] = -1
			}
			for _, point := range series[0].Points {
				usage, _ := strconv.ParseFloat(point[column_index_mapper["usage"]].(string), 32)
				results[int(point[column_index_mapper["time_index"]].(float64))] = float32(usage)
			}
			this.Data["json"] = results
		} else {
			this.Data["json"] = nil
		}
	} else {
		query := fmt.Sprintf("SELECT time_index, out_bytes, in_bytes, out_packets, in_packets FROM net_flow WHERE hardware_addr='%s' AND date='%s'", hardwareAddr, dateStr)
		series, _ := client.Query(query)
		if len(series) > 0 && len(series[0].Points) > 0 {
			column_index_mapper := make(map[string]int)
			for index, value := range series[0].Columns {
				column_index_mapper[value] = index
			}
			dataContainerLength := int(series[0].Points[0][column_index_mapper["time_index"]].(float64)) + 1

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
			for _, point := range series[0].Points {
				time_index := int(point[column_index_mapper["time_index"]].(float64))
				results.Out_bytes[time_index] = int(point[column_index_mapper["out_bytes"]].(int64))
				results.In_bytes[time_index] = int(point[column_index_mapper["in_bytes"]].(int64))
				results.Out_packets[time_index] = int(point[column_index_mapper["out_packets"]].(int64))
				results.In_packets[time_index] = int(point[column_index_mapper["in_packets"]].(int64))
			}
			this.Data["json"] = results
		} else {
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

	type ResultType struct {
		Hardware_addr     string
		Date              string
		Ping_time_index   int
		Ping_results      []PingResultType
		Telnet_time_index int
		Telnet_results    []TelnetResultType
	}

	pingResults := make([]PingResultType, 0, 100)
	telnetResults := make([]TelnetResultType, 0, 100)
	var pingTimeIndex int
	var telnetTimeIndex int
	query := fmt.Sprintf("SELECT time_index, target_ip, response_time FROM ping_accessibility WHERE hardware_addr='%s' AND date='%s'", hardwareAddr, dateStr)
	series, _ := client.Query(query)
	if len(series) > 0 && len(series[0].Points) > 0 {
		column_index_mapper := make(map[string]int)
		for index, value := range series[0].Columns {
			column_index_mapper[value] = index
		}
		newestTimeIndex := int(series[0].Points[0][column_index_mapper["time_index"]].(float64)) + 1
		pingTimeIndex = newestTimeIndex
		for _, point := range series[0].Points {
			if int(point[column_index_mapper["time_index"]].(float64)) < newestTimeIndex {
				break
			} else {
				pingResults = append(pingResults, PingResultType{Target_ip: point[column_index_mapper["target_ip"]].(string), Response_time: int(point[column_index_mapper["response_time"]].(int64))})
			}
		}
	}

	query = fmt.Sprintf("SELECT time_index, target_url, status FROM telnet_accessibility WHERE hardware_addr='%s' AND date='%s'", hardwareAddr, dateStr)
	series, _ = client.Query(query)
	if len(series) > 0 && len(series[0].Points) > 0 {
		column_index_mapper := make(map[string]int)
		for index, value := range series[0].Columns {
			column_index_mapper[value] = index
		}
		newestTimeIndex := int(series[0].Points[0][column_index_mapper["time_index"]].(float64)) + 1
		telnetTimeIndex = newestTimeIndex
		for _, point := range series[0].Points {
			if int(point[column_index_mapper["time_index"]].(float64)) < newestTimeIndex {
				break
			} else {
				telnetResults = append(telnetResults, TelnetResultType{Target_url: point[column_index_mapper["target_url"]].(string), Status: point[column_index_mapper["status"]].(string)})
			}
		}
	}
	this.Data["json"] = ResultType{Hardware_addr: hardwareAddr, Date: now.Format("2006-01-02"), Ping_time_index: pingTimeIndex, Ping_results: pingResults, Telnet_time_index: telnetTimeIndex, Telnet_results: telnetResults}
	this.ServeJson()
}
