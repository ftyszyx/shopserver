package shop

import (
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type ShopNoticeController struct {
	base.BaseController
}

func (self *ShopNoticeController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.Logcommon(data, oldinfo)
	return nil
}
func (self *ShopNoticeController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *ShopNoticeController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *ShopNoticeController) Del() {
	self.DelCommonAndReturn(self)
}
