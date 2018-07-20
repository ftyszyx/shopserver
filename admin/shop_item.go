package admin

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type ShopItemController struct {
	BaseController
}

func (self *ShopItemController) Add() {
	self.AddCommon(self)
}

func (self *ShopItemController) Edit() {
	self.EditCommon(self)
}

func (self *ShopItemController) Del() {
	self.DelCommon(self)
}

//上下架
func (self *ShopItemController) DownUpShop() {
	self.CheckFieldExit(self.postdata, "id", "id为空")
	self.CheckFieldExit(self.postdata, "is_onsale", "字段空")
	id := self.postdata["id"].(string)
	is_onsale := self.postdata["is_onsale"]
	o := orm.NewOrm()
	_, err := o.Raw(fmt.Sprintf("update %s set `is_onsale`=? where `id`=?", self.model.TableName()), is_onsale, id).Exec()
	if err == nil {
		self.AjaxReturnSuccess("", nil)
		return
	}
	self.AjaxReturnError(err.Error())
}
