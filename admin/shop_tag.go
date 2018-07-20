package admin

type ShopTagController struct {
	BaseController
}

func (self *ShopTagController) Add() {
	self.AddCommon(self)
}

func (self *ShopTagController) Edit() {
	self.EditCommon(self)
}

func (self *ShopTagController) Del() {
	self.DelCommon(self)
}
