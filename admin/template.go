package admin

type TemplateController struct {
	BaseController
}

func (self *TemplateController) Add() {
	self.AddCommon(self)
}

func (self *TemplateController) Edit() {
	self.EditCommon(self)
}

func (self *TemplateController) Del() {
	self.DelCommon(self)
}
