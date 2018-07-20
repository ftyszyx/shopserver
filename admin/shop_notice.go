package admin

type ShopNoticeController struct {
	BaseController
}

func (self *ShopNoticeController) Add() {
	self.AddCommon(self)
}

func (self *ShopNoticeController) Edit() {
	self.EditCommon(self)
}

func (self *ShopNoticeController) Del() {
	self.DelCommon(self)
}
