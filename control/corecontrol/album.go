package corecontrol

import (
	"fmt"

	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type AlbumController struct {
	base.BaseController
}

func (self *AlbumController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.GetMethod() == "ChangeDefault" {
		self.AddLog(fmt.Sprintf("postdata:%+v ", self.GetPost()))
	} else if self.GetMethod() == "ChangeCover" {
		self.AddLog(fmt.Sprintf("postdata:%+v ", self.GetPost()))
	} else {
		self.Logcommon(data, oldinfo)
	}
	return nil
}

func (self *AlbumController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *AlbumController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *AlbumController) Del() {
	self.DelCommonAndReturn(self)
}

func (self *AlbumController) ChangeDefault() {
	self.CheckFieldExitAndReturn(self.GetPost(), "default", "图片不能为空")
	self.CheckFieldExitAndReturn(self.GetPost(), "id", "要修改的相册不能为空")
	changedata := make(map[string]interface{})
	changedata["default"] = self.GetPost()["default"]
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetPost()["id"])
}

func (self *AlbumController) ChangeCover() {

	self.CheckFieldExitAndReturn(self.GetPost(), "cover_pic", "图片不能为空")
	self.CheckFieldExitAndReturn(self.GetPost(), "id", "要修改的相册不能为空")
	changedata := make(map[string]interface{})
	changedata["cover_pic"] = self.GetPost()["cover_pic"]
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetPost()["id"])
}
