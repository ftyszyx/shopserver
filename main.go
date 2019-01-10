package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/baiduai"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/initdata"
	"github.com/zyx/shop_server/routers"
	"github.com/zyx/shop_server/wechat"
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

func TempleMinus(a, b int) string {
	return strconv.Itoa(a - b)
}

func TempleAdd(a, b int) string {
	return strconv.Itoa(a + b)
}

var backupsql bool
var initpaycode bool
var updateship bool
var configfile string
var cleanSystem bool
var SERVER_CONF_DATA_PATH string

func initParam() {
	flag.BoolVar(&backupsql, "backupsql", false, "specify backupsql defaults to false.")
	flag.BoolVar(&initpaycode, "initpaycode", false, "specify initpaycode defaults to false.")
	flag.BoolVar(&updateship, "updateship", false, "specify updateship defaults to false.")
	flag.BoolVar(&cleanSystem, "cleansystem", false, "specify cleansystem defaults to false.")
	flag.StringVar(&configfile, "conf", "shop.conf", "specify port defaults to shop.conf")
	flag.Parse()
}

func initlog() {
	//日志
	logpath := beego.AppConfig.String("server.logpath")
	fullpath := fmt.Sprintf(SERVER_CONF_DATA_PATH+"logs/%s.log", logpath)

	logs.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"%s","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`, fullpath))

	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
}

func initPath() {
	libs.MakePath("logs/")
	//系统临时文件夹
	libs.MakePath(beego.AppConfig.String("site.tempfolder"))
}

//shopserver -backupsql true
//shopserver -initpaycode true
//shopserver -conf home.conf
func main() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error("err:%+v\n statck:\n %s", err, string(debug.Stack()))
		}
	}()

	SERVER_CONF_DATA_PATH = os.Getenv("SERVER_CONF_DATA_PATH")

	initParam()

	err := beego.LoadAppConfig("ini", SERVER_CONF_DATA_PATH+"conf/"+configfile)
	if err != nil {
		panic(err.Error())
	}
	initPath()
	initlog()
	//备份数据库
	if backupsql == true {
		models.InitDatabase()
		saveDatabaseTask()
		return
	} else if initpaycode == true {
		//当时加了支付码后，为了给老订单生成支付码用的
		models.InitDatabase()
		initModel()
		initpaycodeData()
		return
	} else if updateship == true {
		models.InitDatabase()
		initModel()
		updateAllShip()
		return
	} else if cleanSystem == true {
		models.InitDatabase()
		initModel()
		CleanSystem()
		return
	}

	allowliststr := beego.AppConfig.String("http.allow_orgin")
	allowlist := strings.Split(allowliststr, ",")
	logs.Info("allowlist:%v", allowlist)
	//跨域访问
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     allowlist,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "token", "X_Requested_With", "uid", "x-requested-with", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
	public := beego.AppConfig.String("site.publicname")
	//定static目录
	beego.SetStaticPath("/static", "static/"+public)
	beego.SetStaticPath("/MP_verify_MYO5aAi6qGBYezdL.txt", "static/MP_verify_MYO5aAi6qGBYezdL.txt")
	//模块初始化
	models.InitDatabase()
	models.InitCaptchaCode()
	initModel()
	routers.InitAllRoute()
	wechat.InitWechat()
	baiduai.InitBaiduAiIDcard() //身份证识别
	//模板函数
	beego.AddFuncMap("formatTimestamp", formatTimestamp)
	beego.AddFuncMap("minus", TempleMinus)
	beego.AddFuncMap("add", TempleAdd)
	beego.Run()
}

func initModel() {
	appname := beego.AppConfig.String("appname")
	if appname == "shop" {
		initdata.InitShopModel()
	} else if appname == "ship" {
		initdata.InitLogisticsModel()
	} else if appname == "home" {
		initdata.InitHomeModel()
	}
}
