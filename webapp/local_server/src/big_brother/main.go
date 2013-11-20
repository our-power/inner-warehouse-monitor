package main

import (
	"fmt"

	"github.com/astaxie/beego"

	"big_brother/controllers"
	"big_brother/models"

	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)


func init() {
	orm.RegisterDataBase("default", beego.AppConfig.String("dbdriver"), beego.AppConfig.String("dbsourcename"))
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
	beego.Router("/api/serverinuse", &controllers.ApiController{}, "GET:GetServerListInUse")
	beego.Run()
}
