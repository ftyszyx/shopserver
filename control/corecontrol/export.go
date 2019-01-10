package corecontrol

import (
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type ExportController struct {
	base.BaseController
}

func (self *ExportController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {

	self.Logcommon(data, oldinfo)

	return nil
}
func (self *ExportController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *ExportController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *ExportController) Del() {
	self.DelCommonAndReturn(self)
}
