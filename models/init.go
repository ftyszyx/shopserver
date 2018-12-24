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

//模块名
const USER = "user"
const USERGROUP = "usergroup"
const MODULE = "module"
const LOG = "log"
const POSTTYPE = "posttype"
const POST = "post"
const PHOTO = "photo"
const ALBUM = "album"
const MEMBER = "member"
const CONFIG = "config"
const ADS = "ads"
const ADSPOS = "adspos"
const SHOP_BRAND = "shopbrand"
const SHOP_ITEM = "shopitem"
const SHOP_ITEMTYPE = "shopitemtype"
const SHOP_NOTICE = "shopnotice"
const SHOP_TAG = "shoptag"
const SHOP_ORDER = "shoporder"
const EXPORT_TASK = "exporttask"
const EXPORT = "export"
const DATABASE = "database"
const LOGISTICS = "logistics"
const PAYCODE = "paycode"
const LOGISTICS_TASK = "logisticstask"

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

//初始化数据库
func Init() {
	logs.Info("init models")

	//initModel()
}

func GetModel(modelname string) ModelInterface {
	return allModels[modelname]
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

func InitModel() {
	allModels = make(map[string]ModelInterface)
	allModels[USER] = &User{Model{"aq_user", make(map[string]db.Params)}}
	allModels[USERGROUP] = &UserGroup{Model{"aq_usergroup", make(map[string]db.Params)}}
	allModels[MODULE] = &Module{Model{"aq_module", make(map[string]db.Params)}}
	allModels[LOG] = &Log{Model{"aq_log", make(map[string]db.Params)}}
	allModels[POST] = &Post{Model{"aq_post", make(map[string]db.Params)}}
	allModels[POSTTYPE] = &PostType{Model{"aq_post_type", make(map[string]db.Params)}}
	allModels[ALBUM] = &Album{Model{"aq_album", make(map[string]db.Params)}}
	allModels[PHOTO] = &Photo{Model{"aq_photo", make(map[string]db.Params)}}
	allModels[CONFIG] = &Config{Model{"aq_config", make(map[string]db.Params)}}
	allModels[ADS] = &Ads{Model{"aq_ads", make(map[string]db.Params)}}
	allModels[ADSPOS] = &AdsPos{Model{"aq_ads_pos", make(map[string]db.Params)}}
	allModels[SHOP_BRAND] = &ShopBrand{Model{"aq_brand", make(map[string]db.Params)}}
	allModels[SHOP_ITEM] = &ShopItem{Model{"aq_item", make(map[string]db.Params)}}
	allModels[SHOP_ITEMTYPE] = &ShopItemType{Model{"aq_item_type", make(map[string]db.Params)}}
	allModels[SHOP_NOTICE] = &ShopNotice{Model{"aq_notice", make(map[string]db.Params)}}
	allModels[SHOP_ORDER] = &ShopOrder{Model{"aq_order", make(map[string]db.Params)}}
	allModels[SHOP_TAG] = &ShopTag{Model{"aq_tag", make(map[string]db.Params)}}
	allModels[EXPORT_TASK] = &ExportTask{Model{"aq_export_task", make(map[string]db.Params)}}
	allModels[EXPORT] = &Export{Model{"aq_export", make(map[string]db.Params)}}
	allModels[DATABASE] = &DataBase{Model{"aq_database", make(map[string]db.Params)}}
	allModels[LOGISTICS] = &Logistics{Model{"aq_logistics", make(map[string]db.Params)}}
	allModels[PAYCODE] = &PayCode{Model{"aq_paycode", make(map[string]db.Params)}}
	allModels[LOGISTICS_TASK] = &LogisticsTask{Model{"aq_logistics_task", make(map[string]db.Params)}}

	for _, value := range allModels {
		value.Init()
	}
}
