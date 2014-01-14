package models

import (
	"github.com/astaxie/beego/orm"
)

/*
type Common_fields struct {
	Id int
	Date string
	Time_index int
	Ip string
	Host_name string
	Hardware_addr string
}
*/

/*
type Cpu_usage struct {
	Id int
	Date string
	Time_index int
	Ip string
	Host_name string
	Hardware_addr string
	Usage float32
}

type Mem_usage struct {
	Id int
	Date string
	Time_index int
	Ip string
	Host_name string
	Hardware_addr string
	Usage float32
}

type Net_flow struct {
	Id int
	Date string
	Time_index int
	Ip string
	Host_name string
	Hardware_addr string
	Out_bytes int
	In_bytes int
	Out_packets int
	In_packets int
}

type Heart_beat struct {
	Id int
	Date string
	Time_index int
	Ip string
	Host_name string
	Hardware_addr string
	Alive int
}
*/
type Register struct {
	Id            int
	Date          string
	Time_index    int
	Ip            string
	Host_name     string
	Hardware_addr string
	Agent_version string
	Machine_role  string
	Status        int
}

type User struct {
	Id int
	Name string
	Passwd string
	Email string
	RegisterTime string
	LastLogin string
	Role string
}

type Role struct {
	Id int
	RoleType string
	Permission string
}

type Trace struct {
	Id int
	User string
	DoWhat string
	ThatTime string
}

/*
type Ping_accessibility struct {
	Id            int
	Date          string
	Time_index    int
	Ip            string
	Host_name     string
	Hardware_addr string
	Target_ip     string
	Response_time int
}

type Telnet_accessibility struct {
	Id            int
	Date          string
	Time_index    int
	Ip            string
	Host_name     string
	Hardware_addr string
	Target_url    string
	Status        string
}
*/
func InitModels() {
	//orm.RegisterModel(new(Cpu_usage), new(Mem_usage), new(Net_flow), new(Heart_beat), new(Register), new(Ping_accessibility), new(Telnet_accessibility))
	orm.RegisterModel(new(Register), new(User), new(Role), new(Trace))
}
