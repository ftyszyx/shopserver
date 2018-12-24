package admin

import (
	"github.com/pkg/errors"
	"github.com/zyx/shop_server/libs/db"
)

type ShopItemTypeController struct {
	BaseController
}

func (self *ShopItemTypeController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	self.logcommon(data, oldinfo)
	return nil
}
func (self *ShopItemTypeController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *ShopItemTypeController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *ShopItemTypeController) Del() {
	self.AjaxReturnError(errors.New("不能删除商品类型"))
	self.DelCommonAndReturn(self)
}
