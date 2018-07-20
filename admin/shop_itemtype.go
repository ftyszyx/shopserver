package admin

type ShopItemTypeController struct {
	BaseController
}

func (self *ShopItemTypeController) Add() {
	self.AddCommon(self)
}

func (self *ShopItemTypeController) Edit() {
	self.EditCommon(self)
}

func (self *ShopItemTypeController) Del() {
	self.DelCommon(self)
}
