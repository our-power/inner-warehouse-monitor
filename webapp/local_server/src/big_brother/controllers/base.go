package controllers

import (
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare(){
	loginName := this.GetSession("login_name")
	if loginName == nil {
		this.Redirect("/login", 302)
		return
	}
}
