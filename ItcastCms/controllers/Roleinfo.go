package controllers

import (
	"github.com/astaxie/beego"
	"time"
	"github.com/astaxie/beego/orm"
	"ItcastCms/models"
)

type RoleInfoController struct {
	beego.Controller
}

func (this *RoleInfoController)Index()  {
	this.TplName="RoleInfo/Index.html"
}
func (this *RoleInfoController)ShowAddRole()  {
	this.TplName="RoleInfo/ShowAddRole.html"
}
func (this *RoleInfoController)AddRole()  {
	var roleInfo=models.RoleInfo{}
	roleInfo.RoleName=this.GetString("roleName")
	roleInfo.Remark=this.GetString("roleRemark")
	roleInfo.DelFlag=0
	roleInfo.AddDate=time.Now()
	roleInfo.ModifDate=time.Now()
	o:=orm.NewOrm()
	_,err:=o.Insert(&roleInfo)
	if err==nil{
		this.Data["json"]=map[string]interface{}{"flag":"ok"}
	}else{
		this.Data["json"]=map[string]interface{}{"flag":"no"}
	}
	this.ServeJSON()

}
func (this *RoleInfoController)GetRoleInfo()  {
	pageIndex,_:=this.GetInt("page")
	pageSize,_:=this.GetInt("rows")
	start:=(pageIndex-1)*pageSize
	o:=orm.NewOrm()
	var roles []models.RoleInfo
	o.QueryTable("role_info").Filter("del_flag",0).OrderBy("Id").Limit(pageSize,start).All(&roles)
	count,_:=o.QueryTable("role_info").Filter("del_flag",0).Count()
	this.Data["json"]=map[string]interface{}{"rows":roles,"total":count}
	this.ServeJSON()
}