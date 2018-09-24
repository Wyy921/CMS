package models

import (
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

//用户信息
type UserInfo struct {
	Id        int       //id
	UserName  string    //用户名
	UserPwd   string    //密码
	Remark    string    //备注
	AddDate   time.Time //添加日期
	ModifDate time.Time //修改日期
	DelFlag   int       //删除标记
	Roles  []*RoleInfo `orm:"rel(m2m)"`
}

//角色信息
type RoleInfo struct {
	Id        int
	RoleName  string      `orm:"size(32)"`
	Remark    string
	DelFlag   int
	AddDate   time.Time
	ModifDate time.Time
	Users     []*UserInfo `orm:"reverse(many)"`
}

func init() {
	var dbhost string
	var dbport string
	var dbuser string
	var dbpassword string
	var db string
	//获取配置文件中对应的配置信息
	dbhost = beego.AppConfig.String("dbhost")
	dbport = beego.AppConfig.String("dbport")
	dbuser = beego.AppConfig.String("dbuser")
	dbpassword = beego.AppConfig.String("dbpassword")
	db = beego.AppConfig.String("db")
	orm.RegisterDriver("mysql", orm.DRMySQL) //注册mysql Driver
	//构造conn连接
	conn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + db + "?charset=utf8&loc=Asia%2FShanghai"
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)

	orm.RegisterModel(new(UserInfo),new(RoleInfo)) //注册模型
	orm.RunSyncdb("default", false, true)

}
