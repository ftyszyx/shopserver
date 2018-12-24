package admin

import "github.com/zyx/shop_server/libs/db"

type TemplateController struct {
	BaseController
}

func (self *TemplateController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.logcommon(data, oldinfo)
	return nil
}
func (self *TemplateController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *TemplateController) Edit() {
	self.EditCommonAndReturn(self)

}

func (self *TemplateController) Del() {
	self.DelCommonAndReturn(self)
}
