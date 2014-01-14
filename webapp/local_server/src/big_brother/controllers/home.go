package controllers

type HomeController struct {
	BaseController
}

func (this *HomeController) Get() {
	this.Data["login_name"] = this.GetSession("login_name")
	this.TplNames = "index.html"
}
