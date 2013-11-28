package controllers

import (
	"time"
	"runtime"
	"strings"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type ApiController struct {
	beego.Controller
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

func getContainerLength(indicator string, db *sql.DB) int {
	rows, err := db.Query("select time_index from " + indicator + " order by time_index desc limit 1")
	if err != nil {
		return 0
	}
	var dataContainerLength int = 0
	for rows.Next() {
		var time_index int
		err = rows.Scan(&time_index)
		if err == nil {
			dataContainerLength = time_index + 1
		}
	}
	return dataContainerLength
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

	var dataContainerLength int
	path := beego.AppConfig.String("multidb") + dateStr + "/" + strings.Replace(hardwareAddr, ":", "_", -1) + "/"
	dbName := path
	dbName = dbName + indicator + ".db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		this.Data["json"] = nil
		goto RETURN
	}
	defer db.Close()
	dataContainerLength = getContainerLength(indicator, db)

	if indicator == "cpu_usage" {
		var results []float32
		if dataContainerLength > 0 {
			results = make([]float32, dataContainerLength)
			for index := 0; index < dataContainerLength; index++ {
				results[index] = -1
			}
		} else {
			this.Data["json"] = nil
			goto RETURN
		}
		rows, err := db.Query("select time_index,usage from cpu_usage order by time_index")
		if err != nil {
			this.Data["json"] = nil
			goto RETURN
		}
		for rows.Next() {
			var time_index int
			var usage float32
			err = rows.Scan(&time_index, &usage)
			if err == nil {
				results[time_index] = usage
			}
		}
		this.Data["json"] = results

	} else if indicator == "mem_usage" {
		var results []float32
		if dataContainerLength > 0 {
			results = make([]float32, dataContainerLength)
			for index := 0; index < dataContainerLength; index++ {
				results[index] = -1
			}
		} else {
			this.Data["json"] = nil
			goto RETURN
		}
		rows, err := db.Query("select time_index,usage from mem_usage order by time_index")
		if err != nil {
			this.Data["json"] = nil
			goto RETURN
		}
		for rows.Next() {
			var time_index int
			var usage float32
			err = rows.Scan(&time_index, &usage)
			if err == nil {
				results[time_index] = usage
			}
		}
		this.Data["json"] = results
	} else if indicator == "net_flow" {
		if dataContainerLength > 0 {
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
			rows, err := db.Query("select time_index,out_bytes,in_bytes,out_packets,in_packets from net_flow order by time_index")
			if err != nil {
				this.Data["json"] = nil
				goto RETURN
			}
			for rows.Next() {
				var time_index int
				var out_bytes, in_bytes, out_packets, in_packets int
				err = rows.Scan(&time_index, &out_bytes, &in_bytes, &out_packets, &in_packets)
				if err == nil {
					results.Out_bytes[time_index] = out_bytes
					results.In_bytes[time_index] = in_bytes
					results.Out_packets[time_index] = out_packets
					results.In_packets[time_index] = in_packets
				}
			}
			this.Data["json"] = results
		} else {
			this.Data["json"] = nil
		}
	}
RETURN:
	this.ServeJson()
}
