package main

import (
	"big_brother/controllers"
	"big_brother/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	registerDB := beego.AppConfig.String("multidb") + "register.db"
	adminDB := beego.AppConfig.String("admindb")
	orm.RegisterDataBase("default", beego.AppConfig.String("dbdriver"), registerDB)
	orm.RegisterDataBase("admin", beego.AppConfig.String("dbdriver"), adminDB)
	models.InitModels()
	controllers.InitControllers()
}

func addInTpl(x, y int) int {
	return x + y
}

/*
func inForTpl(item interface{}, item_array interface{}) bool {
	val := reflect.ValueOf(item_array)
	exist := false
	if val.Kind() == reflect.Slice {
		slice_len := val.Len()
		for index := 0; index < slice_len; index++ {
			if item.(string) == val.Index(index).Interface().(string) {
				exist = true
				break
			}
		}
	}
	return exist
}
*/

func main() {
	run_mode := beego.AppConfig.String("runmode")
	addr := beego.AppConfig.String("httpaddr")
	port, _ := beego.AppConfig.Int("httpport")

	fmt.Println("启动应用....")
	fmt.Printf("运行模式：%s\n", run_mode)
	fmt.Printf("请访问：%s:%d\n", addr, port)

	beego.SessionOn = true;
	beego.SessionProvider = "file"
	beego.SessionSavePath = "./sessions"

	beego.Router("/", &controllers.HomeController{})
	beego.Router("/search_machine", &controllers.SearchController{}, "GET:GetSearchPage")
	beego.Router("/filter_machine_list", &controllers.SearchController{}, "GET:FilterMachineList")
	beego.Router("/indicators_shortcut", &controllers.SearchController{}, "GET:GetIndicatorsByMac")

	beego.Router("/manage/list_machine", &controllers.ManageController{}, "GET:GetManagePage")
	beego.Router("/manage/del_machine", &controllers.ManageController{}, "POST:DelMachine")

	beego.Router("/api/get_machine_indicator_data", &controllers.ApiController{}, "GET:GetMachineIndicatorData")
	beego.Router("/api/get_machine_accessibility_data", &controllers.ApiController{}, "GET:GetMachineAccessibilityData")
	beego.Router("/api/status_overview", &controllers.ApiController{}, "GET:GetStatusOverview")

	beego.Router("/login", &controllers.AdminController{}, "GET,POST:Login")
	beego.Router("/logout", &controllers.AdminController{}, "GET:Logout")

	beego.Router("/admin", &controllers.AdminController{}, "GET:GetAdminPage")
	beego.Router("/admin/api/change_passwd", &controllers.AdminController{}, "POST:ChangePasswd")
	beego.Router("/admin/api/del_user", &controllers.AdminController{}, "POST:DelUser")
	beego.Router("/admin/api/del_role", &controllers.AdminController{}, "POST:DelRole")
	beego.Router("/admin/api/add_user", &controllers.AdminController{}, "POST:AddUser")
	beego.Router("/admin/api/modify_user", &controllers.AdminController{}, "POST:ModifyUser")
	beego.Router("/admin/api/add_role", &controllers.AdminController{}, "POST:AddRole")
	beego.Router("/admin/api/modify_role", &controllers.AdminController{}, "POST:ModifyRole")

	beego.AddFuncMap("add", addInTpl)
	beego.Run()
}
