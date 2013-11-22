package controllers

import (
	"fmt"
	"time"
	"strings"
	"strconv"

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
				fmt.Println("In net_flow handle process...")
				type NcDataType struct {
					Out_bytes   []float64
					In_bytes    []float64
					Out_packets []float64
					In_packets  []float64
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
						ncNum := len(strings.Split(rows[0]["Out_bytes"].(string), ","))
						ncData := make([]NcDataType, ncNum)
						for index := 0; index < ncNum; index++ {
							outBytes := make([]float64, dataContainerLength)
							inBytes := make([]float64, dataContainerLength)
							outPackets := make([]float64, dataContainerLength)
							inPackets := make([]float64, dataContainerLength)
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
								ncData[i].Out_bytes[int(row["Time_index"].(int64))] = ob
								ib, _ := strconv.ParseFloat(ncsInByte[i], 32)
								ncData[i].In_bytes[int(row["Time_index"].(int64))] = ib
								op, _ := strconv.ParseFloat(ncsOutPacket[i], 32)
								ncData[i].Out_packets[int(row["Time_index"].(int64))] = op
								ip, _ := strconv.ParseFloat(ncsInPacket[i], 32)
								ncData[i].In_packets[int(row["Time_index"].(int64))] = ip
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


func (this *IndicatorDataController) GetMachineIndicatorsData() {

}
