package admin

import "github.com/zyx/shop_server/libs/db"

type ShopNoticeController struct {
	BaseController
}

func (self *ShopNoticeController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.logcommon(data, oldinfo)
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
