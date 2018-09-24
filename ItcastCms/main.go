package main

import (
	_ "ItcastCms/routers"
	_ "ItcastCms/models"
	"github.com/astaxie/beego"
	"ItcastCms/models"
)

func main() {
	beego.AddFuncMap("checkId",CheckId)
	beego.AddFuncMap("checkRoleAction",CheckRoleAction)
	beego.Run()
}
func CheckId(userExitRoleList []*models.RoleInfo,roleId int)(b bool)  {
	b=false
	for  i:=0;i<len(userExitRoleList);i++  {
		if userExitRoleList[i].Id==roleId {
			b=true
			break
		}
	}
	return
}
//判断角色对应的权限
func CheckRoleAction(roleExtActions []*models.ActionInfo,actionId int)(b bool)  {
	b=false
	for i:=0; i<len(roleExtActions);i++  {
		if roleExtActions[i].Id==actionId{
			b=true
			break
		}
	}
	return
}