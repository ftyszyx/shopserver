package admin

import "github.com/zyx/shop_server/libs/db"

type AdsController struct {
	BaseController
}

func (self *AdsController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.logcommon(data, oldinfo)
	return nil
}

func (self *AdsController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *AdsController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *AdsController) Del() {
	self.DelCommonAndReturn(self)
}
