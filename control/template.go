package control

//模板
import "github.com/zyx/shop_server/libs/db"
import "github.com/zyx/shop_server/control/base"

type TemplateController struct {
	base.BaseController
}

func (self *TemplateController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.Logcommon(data, oldinfo)
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
