package main

import (
	_ "ItcastCms/routers"
	_ "ItcastCms/models"
	"github.com/astaxie/beego"
	"ItcastCms/models"
)

func main() {
	beego.AddFuncMap("checkId",CheckId)
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