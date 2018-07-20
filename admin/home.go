package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/models"
)

type HomeControl struct {
	beego.Controller
}

func (self *HomeControl) Prepare() {

}

func (self *HomeControl) setdata() {
	senddata := make(map[string]interface{})
	senddata["site_name"] = models.ConfigCache.Get("site_name")
	senddata["logo"] = models.ConfigCache.Get("logo")
	senddata["site_phone"] = models.ConfigCache.Get("site_phone")
	senddata["site_email"] = models.ConfigCache.Get("site_email")
	senddata["site_address"] = models.ConfigCache.Get("site_address")
	senddata["site_icp"] = models.ConfigCache.Get("site_icp")
	senddata["site_desc"] = models.ConfigCache.Get("site_desc")
	senddata["site_ower"] = models.ConfigCache.Get("site_ower")
	senddata["site_codepic"] = models.ConfigCache.Get("site_codepic")

	senddata["swipe"] = models.AdsHomeCache.Get("swipecopany")
	senddata["cases"] = models.AdsHomeCache.Get("cases")
	senddata["news"] = models.AdsHomeCache.Get("news")
	senddata["product1"] = models.AdsHomeCache.Get("product1")
	if senddata["product1"] != nil {
		arr := senddata["product1"].([]orm.Params)
		senddata["productname1"] = arr[0]["ads_pos_title"]
	}

	senddata["product2"] = models.AdsHomeCache.Get("product2")
	if senddata["product2"] != nil {
		arr := senddata["product2"].([]orm.Params)
		senddata["productname2"] = arr[0]["ads_pos_title"]
	}
	senddata["about"] = models.AdsHomeCache.Get("about")
	senddata["contact"] = models.AdsHomeCache.Get("contact")
	senddata["joinus"] = models.AdsHomeCache.Get("joinus")

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
	postid := self.Ctx.Input.Param(":id")
	postmodel := models.GetModel(models.POST)
	postinfo := postmodel.GetInfoById(postid)
	if postinfo != nil {
		self.Data["postinfo"] = postinfo
	}
	self.Layout = "home/layout.html"
	self.TplName = "home/post.html"
}
