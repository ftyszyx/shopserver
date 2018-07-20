package admin

type ExportController struct {
	BaseController
}

func (self *ExportController) Add() {
	self.AddCommon(self)
}

func (self *ExportController) Edit() {
	self.EditCommon(self)
}

func (self *ExportController) Del() {
	self.DelCommon(self)
}
