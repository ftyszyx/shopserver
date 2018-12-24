package admin

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/zyx/shop_server/libs/db"
)

type ShopItemController struct {
	BaseController
}

func (self *ShopItemController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.method == "ExportCsv" {
		self.AddLog(fmt.Sprintf("download data:%+v ", data))
	} else {
		self.logcommon(data, oldinfo)
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
	self.CheckFieldExitAndReturn(self.postdata, "id", "id为空")
	self.CheckFieldExitAndReturn(self.postdata, "is_onsale", "字段空")
	id := self.postdata["id"].(string)
	is_onsale := self.postdata["is_onsale"]
	// o := orm.NewOrm()
	_, err := self.dboper.Raw(fmt.Sprintf("update %s set `is_onsale`=? where `id`=?", self.model.TableName()), is_onsale, id).Exec()
	if err == nil {
		self.AddLog(fmt.Sprintf("data:%+v ", self.postdata))
		self.AjaxReturnSuccess("", nil)
		return
	}

	self.AjaxReturnError(errors.WithStack(err))
}
