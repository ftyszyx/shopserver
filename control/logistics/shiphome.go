package logistics

import (
	"strconv"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type ShipHomeControl struct {
	beego.Controller
}

func (self *ShipHomeControl) Prepare() {

}

func (self *ShipHomeControl) setdata() {
	senddata := make(map[string]interface{})
	self.Data["data"] = senddata

}

//首页
func (self *ShipHomeControl) Home() {
	self.setdata()
	adsmodel := models.GetModel(names.ADS)
	confcache := adsmodel.Cache()
	self.Data["news"] = confcache.Get("shipnews")
	self.Layout = "shiphome/layout.html"
	self.TplName = "shiphome/index.html"
}

//产品
func (self *ShipHomeControl) Service() {
	self.setdata()

	self.Layout = "shiphome/layout.html"
	self.TplName = "shiphome/service.html"
}

func (self *ShipHomeControl) Track() {
	self.setdata()
	id := self.Ctx.Input.Param(":id")
	self.Data["id"] = id
	logs.Info("trackid:%s", id)

	self.Layout = "shiphome/layout.html"
	self.TplName = "shiphome/track.html"
}

func (self *ShipHomeControl) About() {
	self.setdata()
	self.Layout = "shiphome/layout.html"
	self.TplName = "shiphome/about.html"
}

func (self *ShipHomeControl) Upload() {
	self.setdata()
	self.Layout = "shiphome/layout.html"
	self.TplName = "shiphome/upload.html"
}

func (self *ShipHomeControl) NewsList() {
	self.setdata()
	// logs.Info("get news list")
	dboper := db.NewOper()
	page := self.Ctx.Input.Param(":page")
	postmodel := models.GetModel(names.POST)
	pagenum, err := strconv.Atoi(page)
	if err != nil {
		libs.AjaxReturn(&self.Controller, libs.ErrorCode, err.Error(), nil)
	}
	var numperpage int = 10
	newstype := beego.AppConfig.String("newsposttype")
	search := map[string]interface{}{"post.type": newstype, "post.is_del": 0}
	var reqdata = models.AllReqData{Search: search, Rownum: numperpage, And: true, Page: pagenum}
	var totalnum int
	err, totalnum, _ = postmodel.AllExcCommon(dboper, postmodel, reqdata, libs.GetAll_type_num)
	if err != nil {
		libs.AjaxReturn(&self.Controller, libs.ErrorCode, err.Error(), nil)
	}
	var dataList []db.Params
	err, _, dataList = postmodel.AllExcCommon(dboper, postmodel, reqdata, libs.GetAll_type)
	if err != nil {
		libs.AjaxReturn(&self.Controller, libs.ErrorCode, err.Error(), nil)
	}
	self.Data["newlist"] = dataList
	self.Data["totalpage"] = (totalnum / numperpage) + 1
	self.Data["page"] = pagenum
	self.Data["total"] = totalnum
	self.Layout = "shiphome/layout.html"
	self.TplName = "shiphome/newslist.html"
}

//文章
func (self *ShipHomeControl) New() {
	self.setdata()
	dboper := db.NewOper()
	postid := self.Ctx.Input.Param(":id")
	postmodel := models.GetModel(names.POST)
	postinfo := postmodel.GetInfoById(dboper, postid)
	if postinfo != nil {
		self.Data["postinfo"] = postinfo
	}
	self.Layout = "shiphome/layout.html"
	self.TplName = "shiphome/new.html"
}
