package main

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"

	"big_brother/controllers"
	"big_brother/models"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	dbSourceList := strings.Split(beego.AppConfig.String("dbsourcename"), ";")
	for _, dbSource := range dbSourceList {
		dbName_DbSource := strings.Split(dbSource, ",")
		if dbName_DbSource[0] == "register" {
			// beego的ORM要求必须要有个default的数据库
			orm.RegisterDataBase("default", beego.AppConfig.String("dbdriver"), dbName_DbSource[1])
		}else{
			orm.RegisterDataBase(dbName_DbSource[0], beego.AppConfig.String("dbdriver"), dbName_DbSource[1])
		}
	}
	models.InitModels()
	controllers.InitControllers()
}

func main() {
	run_mode := beego.AppConfig.String("runmode")
	addr := beego.AppConfig.String("httpaddr")
	port, _ := beego.AppConfig.Int("httpport")

	fmt.Println("启动应用....")
	fmt.Printf("运行模式：%s\n", run_mode)
	fmt.Printf("请访问：%s:%d\n", addr, port)

	beego.Router("/", &controllers.HomeController{})
	beego.Router("/api/serverindicator", &controllers.ApiController{}, "GET:GetServerIndicator")
	beego.Router("/api/serverlist", &controllers.ApiController{}, "GET:GetServerList")
	beego.Router("/servergroupbystep", &controllers.NavItemsController{}, "GET:GetServerDataGroupByStep")
	beego.Run()
}
