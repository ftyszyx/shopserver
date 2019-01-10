package shop

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type ShopItemController struct {
	base.BaseController
}

func (self *ShopItemController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.GetMethod() == "ExportCsv" {
		self.AddLog(fmt.Sprintf("download data:%+v ", data))
	} else {
		self.Logcommon(data, oldinfo)
	}

	return nil
}
func (self *ShopItemController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *ShopItemController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *ShopItemController) Del() {
	self.AjaxReturnError(errors.New("不能删除商品"))
	self.DelCommonAndReturn(self)
}

func (self *ShopItemController) ExportCsv() {
	self.ExportCsvCommonAndReturn()
}

//上下架
func (self *ShopItemController) DownUpShop() {
	self.CheckFieldExitAndReturn(self.GetPost(), "id", "id为空")
	self.CheckFieldExitAndReturn(self.GetPost(), "is_onsale", "字段空")
	id := self.GetPost()["id"].(string)
	is_onsale := self.GetPost()["is_onsale"]
	// o := orm.NewOrm()
	_, err := self.GetDb().Raw(fmt.Sprintf("update %s set `is_onsale`=? where `id`=?", self.GetModel().TableName()), is_onsale, id).Exec()
	if err == nil {
		self.AddLog(fmt.Sprintf("data:%+v ", self.GetPost()))
		self.AjaxReturnSuccess("", nil)
		return
	}

	self.AjaxReturnError(errors.WithStack(err))
}
