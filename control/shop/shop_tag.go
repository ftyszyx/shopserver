package shop

import (
	"github.com/pkg/errors"
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type ShopTagController struct {
	base.BaseController
}

func (self *ShopTagController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.Logcommon(data, oldinfo)
	return nil
}

func (self *ShopTagController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *ShopTagController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *ShopTagController) Del() {
	self.AjaxReturnError(errors.New("不能删除标签"))
	self.DelCommonAndReturn(self)
}
