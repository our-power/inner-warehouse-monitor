package controllers

type HomeController struct {
	BaseController
}

func (this *HomeController) Get() {
	this.Data["user_name"] = this.GetSession("user_name")
	this.TplNames = "index.html"
}
