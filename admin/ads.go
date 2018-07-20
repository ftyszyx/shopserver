package admin

type AdsController struct {
	BaseController
}

func (self *AdsController) Add() {
	self.AddCommon(self)
}

func (self *AdsController) Edit() {
	self.EditCommon(self)
}

func (self *AdsController) Del() {
	self.DelCommon(self)
}
