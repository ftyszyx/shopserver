package corecontrol

import (
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type AdsPosController struct {
	base.BaseController
}

func (self *AdsPosController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.Logcommon(data, oldinfo)
	return nil
}

func (self *AdsPosController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *AdsPosController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *AdsPosController) Del() {
	self.DelCommonAndReturn(self)
}
