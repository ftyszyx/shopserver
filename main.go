package main

import (
	"flag"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/baiduai"
	"github.com/zyx/shop_server/models"
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

//shopserver -backupsql true
//shopserver -initpaycode true
//shopserver -conf home.conf
func main() {

	defer func() {
		if err := recover(); err != nil {
			logs.Error("err:%+v\n statck:\n %s", err, string(debug.Stack()))
		}
	}()

	initParam()

	err := beego.LoadAppConfig("ini", "conf/"+configfile)
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
	}

	//当时加了支付码后，为了给老订单生成支付码用的
	if initpaycode == true {
		models.InitDatabase()
		models.InitModel()
		initpaycodeData()
		return
	}

	if updateship == true {
		models.InitDatabase()
		models.InitModel()
		updateAllShip()
		return
	}

	if cleanSystem == true {
		models.InitDatabase()
		models.InitModel()
		CleanSystem()
		return
	}

	//跨域访问
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{
			"http://localhost:8100",
			"http://localhost:8200",
			"http://localhost:8300",
			"http://localhost:8400",
			"http://localhost:8500",
			"https://open.weixin.qq.com",
			"http://adminshop.bqmarket.com",
			"http://adminhome.bqmarket.com",
			"http://tt9pbr.natappfree.cc",
			"http://shop.bqmarket.com",
			"http://ship.bqmarket.com",
			"http://shoptest.bqmarket.com",
			"http://adminship.bqmarket.com",
			"http://testapi.bqmarket.com"},
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
	models.InitModel()
	routers.InitAllRoute()
	wechat.InitWechat()
	baiduai.InitBaiduAiIDcard() //身份证识别
	//模板函数
	beego.AddFuncMap("formatTimestamp", formatTimestamp)
	beego.AddFuncMap("minus", TempleMinus)
	beego.AddFuncMap("add", TempleAdd)
	beego.Run()

}
