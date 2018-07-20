package admin

type ShopBrandController struct {
	BaseController
}

func (self *ShopBrandController) Add() {
	self.AddCommon(self)
}

func (self *ShopBrandController) Edit() {
	self.EditCommon(self)
}

func (self *ShopBrandController) Del() {
	self.DelCommon(self)
}
