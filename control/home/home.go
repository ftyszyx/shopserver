package home

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type HomeControl struct {
	beego.Controller
}

func (self *HomeControl) Prepare() {

}

func (self *HomeControl) setdata() {
	senddata := make(map[string]interface{})
	configmodel := models.GetModel(names.CONFIG)
	adsmodel := models.GetModel(names.ADS)
	confcache := configmodel.Cache()
	senddata["site_name"] = confcache.Get("site_name")
	senddata["logo"] = confcache.Get("logo")
	senddata["site_phone"] = confcache.Get("site_phone")
	senddata["site_email"] = confcache.Get("site_email")
	senddata["site_address"] = confcache.Get("site_address")
	senddata["site_icp"] = confcache.Get("site_icp")
	senddata["site_desc"] = confcache.Get("site_desc")
	senddata["site_ower"] = confcache.Get("site_ower")
	senddata["site_codepic"] = confcache.Get("site_codepic")
	adscache := adsmodel.Cache()
	senddata["swipe"] = adscache.Get("swipecopany")
	senddata["cases"] = adscache.Get("cases")
	senddata["news"] = adscache.Get("news")
	senddata["product1"] = adscache.Get("product1")
	if senddata["product1"] != nil {
		arr := senddata["product1"].([]db.Params)
		senddata["productname1"] = arr[0]["ads_pos_title"]
	}

	senddata["product2"] = adscache.Get("product2")
	if senddata["product2"] != nil {
		arr := senddata["product2"].([]db.Params)
		senddata["productname2"] = arr[0]["ads_pos_title"]
	}
	senddata["about"] = adscache.Get("about")
	senddata["contact"] = adscache.Get("contact")
	senddata["joinus"] = adscache.Get("joinus")

	// logs.Info("sendata:%#v", senddata)
	self.Data["data"] = senddata

}

//首页
func (self *HomeControl) Home() {
	self.setdata()
	self.Layout = "home/layout.html"
	self.TplName = "home/index.html"
}

//新闻
func (self *HomeControl) News() {
	self.setdata()
	self.Layout = "home/layout.html"
	self.TplName = "home/news.html"
}

//产品
func (self *HomeControl) Products() {
	self.setdata()
	self.Layout = "home/layout.html"
	self.TplName = "home/products.html"
}

//文章
func (self *HomeControl) Post() {
	self.setdata()
	dboper := db.NewOper()
	postid := self.Ctx.Input.Param(":id")
	postmodel := models.GetModel(names.POST)
	postinfo := postmodel.GetInfoById(dboper, postid)
	if postinfo != nil {
		self.Data["postinfo"] = postinfo
	}
	self.Layout = "home/layout.html"
	self.TplName = "home/post.html"
}
