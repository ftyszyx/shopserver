package admin

import (
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
)

type SystemController struct {
	BaseController
}

func (self *SystemController) Refresh() {

	models.RefrshAllCache()
	self.AjaxReturn(libs.SuccessCode, "成功", nil)
}
