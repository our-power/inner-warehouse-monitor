package controllers

import (
	"github.com/astaxie/beego"
)

type AdminController struct {
	beego.Controller
}

func validateUser(userName, passwd string) (bool) {
	return true
}

func (this *AdminController) Login() {
	if this.Ctr.Request.Method == "GET" {
		this.TplNames = "login.html"
	}else {
		userName := this.GetString("user_name")
		passwd := this.GetString("password")
		if validateUser(userName, passwd) {
			this.SetSession("login_name", userName)
			this.Redirect("/", 302)
		}else {
			this.Data["err_tips"] = "帐号错误！"
			this.TplNames = "login.html"
		}
	}
}

func (this *AdminController) Logout() {
	this.DelSession("user_name")
	this.Redirect("/login", 302)
}
