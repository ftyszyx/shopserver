package shop

import (
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type ShopController struct {
	base.BaseController
}

func (self *ShopController) GetInfo() {
	var data = make(map[string]interface{})

	var webconfig = make(map[string]interface{})
	confmodel := models.GetModel(names.CONFIG)
	confcache := confmodel.Cache()
	webconfig["site_name"] = confcache.Get("site_name")
	webconfig["logo"] = confcache.Get("logo")
	webconfig["site_phone"] = confcache.Get("site_phone")
	webconfig["site_email"] = confcache.Get("site_email")
	webconfig["site_address"] = confcache.Get("site_address")
	webconfig["site_icp"] = confcache.Get("site_icp")
	webconfig["order_multi"] = confcache.Get("order_multi")
	webconfig["site_pay_code"] = confcache.Get("site_pay_code")
	webconfig["toppic"] = confcache.Get("toppic")
	data["webconfig"] = webconfig
	//adsmodel := models.GetModel(names.ADS)
	//adscache := adsmodel.Cache()
	// data["webhomeads"] = adscache.Get("homeads")
	// data["swipe"] = adscache.Get("swipe")
	// data["notice"] = coredata.NoticeCache

	self.AjaxReturn(libs.SuccessCode, nil, data)
}
