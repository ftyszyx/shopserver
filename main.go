package main

import (
	"flag"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zyx/shop_server/admin"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
	_ "github.com/zyx/shop_server/routers"
	_ "github.com/zyx/shop_server/wechat"
)

func formatTimestamp(stamp string, format string) (datestring string) {
	timeint, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		logs.Error("err:%s", err.Error())
		return
	}
	timeinfo := time.Unix(timeint, 0)
	datestring = beego.Date(timeinfo, format)
	return
}

func saveDatabaseTask() {
	path, err := admin.SaveDatabase()
	if err != nil {
		logs.Error("back up system err:%", err.Error())
	}
	adddata := make(map[string]interface{})
	adddata["name"] = fmt.Sprintf("系统备份 %s", time.Now().Format("2006-01-02"))
	adddata["user_id"] = "1"
	adddata["build_time"] = time.Now().Unix()
	adddata["path"] = path
	o := orm.NewOrm()
	keys, values := libs.SqlGetInsertInfo(adddata)
	_, err = o.Raw(fmt.Sprintf("insert into aq_database (%s) values (%s)", keys, values)).Exec()
	if err != nil {
		logs.Error("back up system err:%", err.Error())
	}
}

var logpath string
var backupsql bool
var public string
var port string

func initParam() {
	flag.StringVar(&logpath, "log", "project", "specify logpath defaults to project.log.")
	flag.BoolVar(&backupsql, "backupsql", false, "specify backupsql defaults to false.")
	flag.StringVar(&public, "public", "", "specify static path defaults to /.")
	flag.StringVar(&port, "port", "9000", "specify port defaults to 9000.")
	flag.Parse()
}

func initlog() {
	//日志
	fullpath := fmt.Sprintf("logs/%s.log", logpath)
	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`, fullpath))

	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
}

func initPath() {
	libs.MakePath("logs/")
	//系统临时文件夹
	libs.MakePath(beego.AppConfig.String("site.tempfolder"))
}

// public 默认"  static目录
// port 默认  9000  端口号
//log   默认 project
//shopserver -port 8000 -public home -log project
//shopserver -backupsql true
func main() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error("err:%s\n statck:\n %s", err, string(debug.Stack()))
		}
	}()
	initParam()
	initPath()
	initlog()

	if backupsql == true {
		models.InitDatabase()
		saveDatabaseTask()
		return
	}

	//跨域访问
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:8100", "http://localhost:8200", "https://open.weixin.qq.com", "http://adminshop.bqmarket.com", "http://shop.bqmarket.com", "http://testapi.bqmarket.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "token", "X_Requested_With", "uid", "x-requested-with", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	//定static目录
	beego.SetStaticPath("/static", "static/"+public)
	beego.SetStaticPath("/MP_verify_MYO5aAi6qGBYezdL.txt", "static/MP_verify_MYO5aAi6qGBYezdL.txt")
	//模块初始化
	models.InitDatabase()
	models.InitCaptchaCode()
	models.InitModel()
	//模板函数
	beego.AddFuncMap("formatTimestamp", formatTimestamp)
	beego.Run(":" + port)

}
