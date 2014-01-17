package controllers

type HomeController struct {
	BaseController
}

func (this *HomeController) Get() {
	this.Data["login_name"] = this.GetSession("login_name")
	if this.GetSession("role_type") == "admin_user" {
		this.Data["admin"] = true
	}
	this.TplNames = "index.html"
}
