package main

import (
	"big_brother/controllers"
	"big_brother/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"path/filepath"
)

func init() {
	appPath, _ := filepath.Abs(os.Args[0])
	beego.AppConfigPath = path.Join(path.Dir(appPath), "conf/app.conf")

	registerDB := beego.AppConfig.String("multidb") + "register.db"
	orm.RegisterDataBase("default", beego.AppConfig.String("dbdriver"), registerDB)
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
	beego.Router("/search_machine", &controllers.SearchController{}, "GET:GetSearchPage")
	beego.Router("/filter_machine_list", &controllers.SearchController{}, "GET:FilterMachineList")

	beego.Router("/manage_machine", &controllers.ManageController{}, "GET:GetManagePage")

	beego.Router("/api/get_machine_indicator_data", &controllers.ApiController{}, "GET:GetMachineIndicatorData")
	beego.Router("/api/get_machine_accessibility_data", &controllers.ApiController{}, "GET:GetMachineAccessibilityData")
	beego.Router("/api/status_overview", &controllers.ApiController{}, "GET:GetStatusOverview")

	beego.Run()
}
