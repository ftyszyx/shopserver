package shop

import (
	"github.com/pkg/errors"
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type ShopBrandController struct {
	base.BaseController
}

func (self *ShopBrandController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.Logcommon(data, oldinfo)
	return nil
}
func (self *ShopBrandController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *ShopBrandController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *ShopBrandController) Del() {
	self.AjaxReturnError(errors.New("不能删除品牌"))
	self.DelCommonAndReturn(self)
}
