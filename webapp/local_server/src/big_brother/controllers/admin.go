package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"big_brother/models"
	"crypto/md5"
	"io"
	"fmt"
	"time"
	"strconv"
)

type AdminController struct {
	beego.Controller
}

func validateUser(userName, passwd string) (user *models.User, role *models.Role, exist bool) {
	o.Using("admin")
	h := md5.New()
	io.WriteString(h, passwd)
	user = new(models.User)
	err := o.QueryTable("user").Filter("name", userName).Filter("passwd", fmt.Sprintf("%x", h.Sum(nil))).One(user)
	if err != nil {
		return nil, nil, false
	}
	role = new(models.Role)
	err = o.QueryTable("role").Filter("id", user.Role_id).One(role)
	if err != nil {
		return nil, nil, false
	}
	_, _ = o.QueryTable("user").Filter("id", user.Id).Update(orm.Params{"last_login": time.Now().Format("2006-01-02 15:04:05")})
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
			this.SetSession("permission", role.Permission)
			this.Redirect("/", 302)
		}else {
			this.Data["err_tips"] = "帐号登录错误！"
			this.TplNames = "login.html"
		}
	}
}

func (this *AdminController) Logout() {
	this.DestroySession()
	//this.DelSession("login_name")
	this.Redirect("/login", 302)
}

func (this *AdminController) GetAdminPage() {
	if this.GetSession("role_type") != "admin_user" {
		this.Redirect("/", 302)
		return
	}
	o.Using("admin")
	var userList []*models.User
	_, err := o.QueryTable("user").Limit(-1).All(&userList, "id", "name", "email", "register_time", "last_login", "role_id")
	if err == nil {
		this.Data["user_list"] = userList
	}
	var roleList []*models.Role
	_, err = o.QueryTable("role").Limit(-1).All(&roleList, "id", "role_type", "permission")
	if err == nil {
		this.Data["role_list"] = roleList
	}
	var traceList []*models.Trace
	_, err = o.QueryTable("trace").Limit(-1).All(&traceList, "id", "user", "do_what", "that_time")
	if err == nil {
		this.Data["trace_list"] = traceList
	}

	if this.GetSession("role_type") == "admin_user" {
		this.Data["admin"] = true
	}
	this.Data["login_name"] = this.GetSession("login_name")
	this.TplNames = "admin.html"
}
