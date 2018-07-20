package admin

import (
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
)

type ShopController struct {
	BaseController
}

func (self *ShopController) GetInfo() {
	var data = make(map[string]interface{})

	var webconfig = make(map[string]interface{})

	webconfig["site_name"] = models.ConfigCache.Get("site_name")
	webconfig["logo"] = models.ConfigCache.Get("logo")
	webconfig["site_phone"] = models.ConfigCache.Get("site_phone")
	webconfig["site_email"] = models.ConfigCache.Get("site_email")
	webconfig["site_address"] = models.ConfigCache.Get("site_address")
	webconfig["site_icp"] = models.ConfigCache.Get("site_icp")
	webconfig["order_multi"] = models.ConfigCache.Get("order_multi")
	webconfig["site_pay_code"] = models.ConfigCache.Get("site_pay_code")
	webconfig["toppic"] = models.ConfigCache.Get("toppic")
	data["webconfig"] = webconfig
	data["webhomeads"] = models.AdsHomeCache.Get("homeads")
	data["swipe"] = models.AdsHomeCache.Get("swipe")
	data["notice"] = models.NoticeCache

	self.AjaxReturn(libs.SuccessCode, nil, data)
}
