package admin

type AdsPosController struct {
	BaseController
}

func (self *AdsPosController) Add() {
	self.AddCommon(self)
}

func (self *AdsPosController) Edit() {
	self.EditCommon(self)
}

func (self *AdsPosController) Del() {
	self.DelCommon(self)
}
