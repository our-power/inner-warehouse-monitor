package controllers

import (
	"github.com/astaxie/beego"
	"big_brother/models"
)

type AdminController struct {
	beego.Controller
}

func validateUser(userName, passwd string) (user *models.User, role *models.Role, exist bool) {
	o.Using("admin")
	err := o.QueryTable("user").Filter("name", userName).Filter("passwd", passwd).One(user)
	if err != nil {
		return nil, nil, false
	}
	err = o.QueryTable("role").Filter("id", user.Role_id).One(role)
	if err != nil {
		return user, nil, true
	}
	return user, role, true
}

func (this *AdminController) Login() {
	if this.Ctx.Request.Method == "GET" {
		this.TplNames = "login.html"
	}else {
		userName := this.GetString("user_name")
		passwd := this.GetString("password")
		user, role, exist := validateUser(userName, passwd)
		if exist {
			this.SetSession("login_name", userName)
			this.SetSession("email", user.Email)
			this.SetSession("role_type", role.Role_type)
			this.SetSession("pemission", role.Permission)
			this.Redirect("/", 302)
		}else {
			this.Data["err_tips"] = "帐号错误！"
			this.TplNames = "login.html"
		}
	}
}

func (this *AdminController) Logout() {
	this.DestroySession()
	//this.DelSession("login_name")
	this.Redirect("/login", 302)
}
