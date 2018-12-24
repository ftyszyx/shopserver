package admin

import "github.com/zyx/shop_server/libs/db"

type LogisticsTaskController struct {
	BaseController
}

func (self *LogisticsTaskController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.logcommon(data, oldinfo)
	return nil
}
func (self *LogisticsTaskController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *LogisticsTaskController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *LogisticsTaskController) Del() {
	self.DelCommonAndReturn(self)
}
