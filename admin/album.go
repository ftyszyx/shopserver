package admin

import (
	"fmt"

	"github.com/zyx/shop_server/libs/db"
)

type AlbumController struct {
	BaseController
}

func (self *AlbumController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.method == "ChangeDefault" {
		self.AddLog(fmt.Sprintf("postdata:%+v ", self.postdata))
	} else if self.method == "ChangeCover" {
		self.AddLog(fmt.Sprintf("postdata:%+v ", self.postdata))
	} else {
		self.logcommon(data, oldinfo)
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
	self.CheckFieldExitAndReturn(self.postdata, "default", "图片不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "id", "要修改的相册不能为空")
	changedata := make(map[string]interface{})
	changedata["default"] = self.postdata["default"]
	self.updateSqlByIdAndReturn(self, changedata, self.postdata["id"])
}

func (self *AlbumController) ChangeCover() {

	self.CheckFieldExitAndReturn(self.postdata, "cover_pic", "图片不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "id", "要修改的相册不能为空")
	changedata := make(map[string]interface{})
	changedata["cover_pic"] = self.postdata["cover_pic"]
	self.updateSqlByIdAndReturn(self, changedata, self.postdata["id"])
}
