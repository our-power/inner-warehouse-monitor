package controllers

import (
	"github.com/astaxie/beego"
)

type SearchController struct {
	beego.Controller
}

func (this *SearchController) GetSearchPage() {
	this.Data["nav_now"] = "search_machine"
	this.TplNames = "search_machine.html"
}
