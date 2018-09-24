package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
	"ItcastCms/models"
	"strconv"
	"strings"
)

type UserInfoController struct {
	beego.Controller
}

//搜索的数据
type UserSearchData struct {
	UserName   string
	Remark     string
	PageIndex  int
	PageSize   int
	TotalCount int64
}

func (this *UserInfoController) Index() {
	this.TplName = "UserInfo/Index.html"
}
func (this *UserInfoController) AddUser() {
	//1:接收表单传递过来的数据
	userName := this.GetString("userName")
	userPwd := this.GetString("userPwd")
	userRemark := this.GetString("userRemark")
	//2:进行表单的校验
	//3:完成数据的保存
	o := orm.NewOrm()
	var userInfo = models.UserInfo{}
	userInfo.AddDate = time.Now()
	userInfo.DelFlag = 0
	userInfo.ModifDate = time.Now()
	userInfo.Remark = userRemark
	userInfo.UserName = userName
	userInfo.UserPwd = userPwd
	_, err := o.Insert(&userInfo)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"flag": "yes"}
	} else {
		this.Data["json"] = map[string]interface{}{"flag": "no"}
	}
	this.ServeJSON()
}
func (this *UserInfoController) GetUserInfo() {
	pageIndex, _ := strconv.Atoi(this.GetString("page"))
	pageSize, _ := strconv.Atoi(this.GetString("rows"))
	username := this.GetString("username")
	remark := this.GetString("remark")

	var userSearch = UserSearchData{} //通过该对象构建相应的搜索数据。
	userSearch.UserName = username
	userSearch.Remark = remark
	userSearch.PageIndex = pageIndex
	userSearch.PageSize = pageSize
	response := userSearch.GetSearchData()
	//start:=(pageIndex-1)*pageSize
	//var users []models.UserInfo
	//o:=orm.NewOrm()
	//o.QueryTable("user_info").Filter("del_flag",0).Limit(pageSize,start).All(&users)
	//count,_:=o.QueryTable("user_info").Filter("del_flag",0).Count()
	//beego.Info(users)
	this.Data["json"] = map[string]interface{}{"rows": response, "total": userSearch.TotalCount}
	this.ServeJSON()

}
func (this *UserSearchData) GetSearchData() ([]models.UserInfo) {
	o := orm.NewOrm()
	//构建搜索的条件
	var temp = o.QueryTable("user_info")
	if this.UserName != "" {
		temp = temp.Filter("user_name__icontains", this.UserName)
	}
	if this.Remark != "" {
		temp = temp.Filter("remark__icontains", this.Remark)
	}
	start := (this.PageIndex - 1) * this.PageSize
	count, _ := temp.Filter("del_flag", 0).Count()
	this.TotalCount = count     //注意，总的记录数的计算，这里是this,表示的是*UserSearchData,也就是引用传递。
	var users []models.UserInfo //构建UserInfo对象
	temp.OrderBy("Id").Filter("del_flag", 0).Limit(this.PageSize, start).All(&users)
	return users //返回搜索查询的数据
}

//完成用户信息的删除操作
func (this *UserInfoController) DeleteUser() {
	var ids = this.GetString("id")    //接收要删除的记录的id
	strIds := strings.Split(ids, ",") //按照逗号进行分隔，返回的是字符串切片
	o := orm.NewOrm()
	//var userInfo =models.UserInfo{}
	for i := 0; i < len(strIds); i++ {
		id, _ := strconv.Atoi(strIds[i])
		var userInfo = models.UserInfo{Id: id}
		//o.QueryTable("user_info").Filter("id",id).All(&userInfo)
		userInfo.DelFlag = 1
		o.Update(&userInfo, "del_flag")
	}
	this.Data["json"] = map[string]interface{}{"flag": "ok"}
	this.ServeJSON()
}

//展示要编辑的用户信息
func (this *UserInfoController) ShowEdit() {
	id, err := strconv.Atoi(this.GetString("id"))
	beego.Info("sgsagggsags", id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"flag": "no", "msg": "类型转换错误"}
	} else {
		//查询出要修改的数据
		o := orm.NewOrm()
		var userInfo = models.UserInfo{Id: id}
		err := o.Read(&userInfo)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"flag": "no", "msg": "查询失败"}
		} else {
			beego.Info("sgjksg", userInfo)
			this.Data["json"] = map[string]interface{}{"flag": "yes", "serverData": userInfo}
		}

	}
	this.ServeJSON()
}
func (this *UserInfoController) EditUser() {
	var userInfo models.UserInfo
	o := orm.NewOrm()
	userInfo.Id, _ = strconv.Atoi(this.GetString("Eidtid"))
	userInfo.UserName = this.GetString("Editname")
	//beego.Info(this.GetString("Editname"))
	userInfo.UserPwd = this.GetString("Editpwd")
	userInfo.Remark = this.GetString("Editremark")
	userInfo.ModifDate = time.Now()
	adddate := this.GetString("EditaddDate")
	//beego.Info(adddate)
	t, _ := time.Parse("2006-01-02T15:04:05+08:00", adddate)
	userInfo.AddDate = t

	//完成数据的更新
	_, err := o.Update(&userInfo)
	if err != nil {
		beego.Info(err)
		this.Data["json"] = map[string]interface{}{"flag": "no", "msg": err}
	} else {
		this.Data["json"] = map[string]interface{}{"flag": "yes", "msg": "数据更新成功!"}
	}
	this.ServeJSON() //生成JSON返回

}

//展示要给用户分配的角色信息
func (this *UserInfoController) ShowSetUserRole() {
	//1：接收传递过来的用户编号
	userId, _ := this.GetInt("userId")
	//2:查询出用户已经有的角色。
	var userInfo models.UserInfo
	o := orm.NewOrm()
	o.QueryTable("user_info").Filter("id", userId).One(&userInfo)
	var userExtRoles []*models.RoleInfo //表示用户已有的角色。
	o.LoadRelated(&userInfo, "Roles")
	for _, role := range userInfo.Roles {
		userExtRoles = append(userExtRoles, role)
	}
	//3:查询出所有角色。
	var allRoles []models.RoleInfo
	o.QueryTable("role_info").Filter("del_flag", 0).All(&allRoles)
	this.Data["allRoles"] = allRoles
	this.Data["userExtRoles"] = userExtRoles
	this.Data["userInfo"] = userInfo

	this.TplName = "UserInfo/ShowSetUserRole.html"
}
//完成用户角色的分配
func (this *UserInfoController)SetUserRole()  {
	userId,_:=this.GetInt("userId")//接收用户的编号
	beego.Info("用户",userId)
	var roleIdList[]int


	allKeys:=this.Ctx.Request.PostForm//获取所有post请求发送过来的表单数据
	//但是，这里我们只需要以cba_开头的数据
	for key,_:=range allKeys {
		if strings.Contains(key,"cba_"){
			id:=strings.Replace(key,"cba_","",-1)
			roleId,_:=strconv.Atoi(id)
			roleIdList=append(roleIdList,roleId)
		}
	}
	//查询用户信息
	o:=orm.NewOrm()
	var userInfo models.UserInfo
	o.QueryTable("user_info").Filter("Id",userId).One(&userInfo)
	//查询用户具有的角色信息
	o.LoadRelated(&userInfo,"Roles")
	m2m:=o.QueryM2M(&userInfo,"Roles")
	//删除用户已经有的角色信息
	o.Begin()//开启事务
	var err1 error
	var err2 error
	for _,role:=range userInfo.Roles{
		_,err1=m2m.Remove(role)
	}
	//重新给用户分配角色信息
	var roleInfo models.RoleInfo
	for i:=0;i<len(roleIdList) ;i++  {
		o.QueryTable("role_info").Filter("Id",roleIdList[i]).One(&roleInfo)
		_,err2=m2m.Add(roleInfo)
	}
	if err2!=nil||err1!=nil{
		o.Rollback()
		this.Data["json"]=map[string]interface{}{"flag":"no"}
	}else{

		o.Commit()
		this.Data["json"]=map[string]interface{}{"flag":"yes"}
	}
	this.ServeJSON()

}