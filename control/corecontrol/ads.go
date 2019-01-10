package corecontrol

import (
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type AdsController struct {
	base.BaseController
}

func (self *AdsController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.Logcommon(data, oldinfo)
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
