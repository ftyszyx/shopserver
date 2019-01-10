package models

import (
	"fmt"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils/captcha"
	"github.com/zyx/shop_server/libs/db"
)

var allModels map[string]ModelInterface //存储所有的数据
var CaptchaCode *captcha.Captcha
var PhoneCodeCache cache.Cache

func InitDatabase() {
	dbhost := beego.AppConfig.String("db.host")
	dbport := beego.AppConfig.String("db.port")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8mb4,utf8"
	fmt.Println(dsn)

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	db.RegisterDataBase("default", "mysql", dsn)
	dbinfo, err := db.GetDB("default")
	if err != nil {
		panic(err)
	}
	dbinfo.DB.SetConnMaxLifetime(time.Minute * 100)
	dbinfo.SetMaxIdleConns(10)
	dbinfo.SetMaxOpenConns(30)
}

func InitCaptchaCode() {
	//验证码
	store := cache.NewMemoryCache()
	CaptchaCode = captcha.NewCaptcha("/captcha/", store)
	beego.InsertFilter("/captcha/*", beego.BeforeRouter, CaptchaCode.Handler)
	//手机验证码
	PhoneCodeCache, _ = cache.NewCache("memory", `{"interval":360}`)
}

func GetModel(modelname string) ModelInterface {
	return allModels[modelname]
}

func GetAllModel() map[string]ModelInterface {
	return allModels
}

func InitAllModel() {
	allModels = make(map[string]ModelInterface)
}

//刷新
func RefrshCache(modelname string) {
	model := GetModel(modelname)
	if model != nil {
		logs.Info("clear  cache:%s", modelname)
		model.ClearCache()
	}
}

func RefrshAllCache() {
	logs.Info("clear all cache")
	for _, value := range allModels {
		value.ClearCache()
	}
}
