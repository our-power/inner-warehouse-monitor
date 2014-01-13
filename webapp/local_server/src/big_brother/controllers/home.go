package controllers

type HomeController struct {
	BaseController
}

func (this *HomeController) Get() {
	this.TplNames = "index.html"
}
