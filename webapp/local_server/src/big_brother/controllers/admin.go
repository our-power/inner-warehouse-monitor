package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"big_brother/models"
	"time"
	"strconv"
)

type AdminController struct {
	beego.Controller
}

var hasNoAdminPermissionMsg = map[string]string{
	"Status": "failure",
	"Msg": "你没有该操作权限!",
}

func validateUser(userName, passwd string) (user *models.User, role *models.Role, exist bool) {
	o.Using("admin")
	user = new(models.User)
	err := o.QueryTable("user").Filter("name", userName).Filter("passwd", GenMd5Passwd(passwd)).One(user)
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


func (this *AdminController) ChangePasswd() {
	hasAdminPermission := HasTheRight("admin", this.GetSession("permission"))
	if (hasAdminPermission) {
		this.Data["json"] = hasNoAdminPermissionMsg
	}else {
		userId, err := this.GetInt("user_id")
		if err != nil {
			this.Data["json"] = map[string]string{
				"Status": "failure",
				"Msg": "提供的用户ID不正确！",
			}
		}else {
			newPasswd := this.GetString("new_passwd")
			_, err := o.QueryTable("user").Filter("id", userId).Update(orm.Params{"passwd": GenMd5Passwd(newPasswd)})
			if err != nil {
				this.Data["json"] = map[string]string{
					"Status": "failure",
					"Msg": "数据库更新操作出错！",
				}
			}else {
				this.Data["json"] = map[string]string{
					"Status": "success",
				}
			}
		}
	}
	this.ServeJson()
}

func (this *AdminController) DelUser() {
	hasAdminPermission := HasTheRight("admin", this.GetSession("permission"))
	if (hasAdminPermission) {
		this.Data["json"] = hasNoAdminPermissionMsg
	}else {
		user_id, err := strconv.Atoi(this.GetString("user_id"))
		if err != nil {
			this.Data["json"] = map[string]string{
				"Status": "failure",
				"Msg": "用户ID不对！",
			}
		}else {
			o.Using("admin")
			_, err = o.Delete(&models.User{Id: user_id})
			if err != nil {
				this.Data["json"] = map[string]string{
					"Status": "failure",
					"Msg": "未能从数据库中删除该用户",
				}
			}else {
				this.Data["json"] = map[string]string{
					"Status": "success",
				}
			}
		}
	}
	this.ServeJson()
}

func (this *AdminController) DelRole() {
	hasAdminPermission := HasTheRight("admin", this.GetSession("permission"))
	if (hasAdminPermission) {
		this.Data["json"] = hasNoAdminPermissionMsg
	}else {
		role_id, err := strconv.Atoi(this.GetString("role_id"))
		if err != nil {
			this.Data["json"] = map[string]string{
				"Status": "failure",
				"Msg": "角色ID不对！",
			}
		}else {
			o.Using("admin")
			_, err = o.Delete(&models.Role{Id: role_id})
			if err != nil {
				this.Data["json"] = map[string]string{
					"Status": "failure",
					"Msg": "未能从数据库中删除该角色",
				}
			}else {
				this.Data["json"] = map[string]string{
					"Status": "success",
				}
			}
		}
	}
	this.ServeJson()
}

func (this *AdminController) AddUser() {
	/*
	"data": {
		"user_name": userName,
        "passwd": newPasswd,
        "email": email,
        "role_type": roleType
    }
	*/
	hasAdminPermission := HasTheRight("admin", this.GetSession("permission"))
	if (hasAdminPermission) {
		this.Data["json"] = hasNoAdminPermissionMsg
	}else {
		userName := this.GetString("user_name")
		md5Passwd := GenMd5Passwd(this.GetString("passwd"))
		email := this.GetString("email")
		roleType, _ := strconv.Atoi(this.GetString("role_type"))
		now := time.Now().Format("2006-01-02 15:04:05")
		o.Using("admin")
		_, err := o.Insert(&models.User{Name: userName, Passwd: md5Passwd, Email: email, Role_id: roleType, Register_time: now})
		if err != nil {
			this.Data["json"] = map[string]string{
				"Status": "failure",
				"Msg": "数据库插入该用户失败！",
			}
		}else {
			this.Data["json"] = map[string]string{
				"Status": "success",
			}
		}
	}
	this.ServeJson()
}

func (this *AdminController) AddRole() {
	/*
	"data": {
		"role_name": roleName,
		"permissions": permissions
	}
	*/
	hasAdminPermission := HasTheRight("admin", this.GetSession("permission"))
	if (hasAdminPermission) {
		this.Data["json"] = hasNoAdminPermissionMsg
	}else {
		roleName := this.GetString("role_name")
		permissions := this.GetString("permissions")
		o.Using("admin")
		_, err := o.Insert(&models.Role{Role_type: roleName, Permission: permissions})
		if err != nil {
			this.Data["json"] = map[string]string{
				"Status": "failure",
				"Msg": "数据库插入该角色失败!",
			}
		}else {
			this.Data["json"] = map[string]string{
				"Status": "success",
			}
		}
	}
	this.ServeJson()
}

func (this *AdminController) ModifyUser() {
	/*
	"data": {
		"user_id": userId,
		"user_name": userName,
		"email": email,
		"role_type": roleType
	}
	*/
	hasAdminPermission := HasTheRight("admin", this.GetSession("permission"))
	if (hasAdminPermission) {
		this.Data["json"] = hasNoAdminPermissionMsg
	}else {
		userId, _ := strconv.Atoi(this.GetString("user_id"))
		userName := this.GetString("user_name")
		email := this.GetString("email")
		roleId, _ := strconv.Atoi(this.GetString("role_type"))
		o.Using("admin")
		_, err := o.QueryTable("user").Filter("id", userId).Update(orm.Params{"name": userName, "email": email, "role_id": roleId})
		if err != nil {
			this.Data["json"] = map[string]string{
				"Status": "failure",
				"Msg": "数据库更新操作失败",
			}
		}else {
			this.Data["json"] = map[string]string{
				"Status": "success",
			}
		}
	}
	this.ServeJson()
}

func (this *AdminController) ModifyRole() {
	/*
	"data": {
		"role_id": roleId,
		"role_name": roleName,
		"permissions": permissions
	}
	*/
	hasAdminPermission := HasTheRight("admin", this.GetSession("permission"))
	if (hasAdminPermission) {
		this.Data["json"] = hasNoAdminPermissionMsg
	}else {
		roleId, _ := strconv.Atoi(this.GetString("role_id"))
		roleName := this.GetString("role_name")
		permissions := this.GetString("permissions")
		_, err := o.QueryTable("role").Filter("id", roleId).Update(orm.Params{"role_type": roleName, "permissions": permissions})
		if err != nil {
			this.Data["json"] = map[string]string{
				"Status": "failure",
				"Msg": "数据库更新操作失败！",
			}
		}else {
			this.Data["json"] = map[string]string{
				"Status": "success",
			}
		}
	}
	this.ServeJson()
}
