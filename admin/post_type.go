package admin

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

type PostTypeController struct {
	BaseController
}

func (self *PostTypeController) BeforeSql(data map[string]interface{}) {
	if self.method == "Abandon" {
		data["is_del"] = self.postdata["is_del"]
	} else if self.method == "Add" {
		data["is_del"] = 0
	} else if self.method == "Edit" {
		if data["parent_id"] != nil && data["parent_id"].(string) == data["id"].(string) {
			self.AjaxReturn(libs.ErrorCode, "父节点不能是自己", nil)
		}
	}
}
func (self *PostTypeController) AfterSql(data map[string]interface{}, oldinfo orm.Params) {
	if self.method == "Add" {
		self.AddLog(fmt.Sprintf("增加文章类型:%+v", data))
	} else if self.method == "Edit" {
		self.AddLog(fmt.Sprintf("修改文章类型:%+v", data))
	} else if self.method == "Del" {
		name := data["title"]
		self.AddLog(fmt.Sprintf("删除文章类型:%s", name))
	} else if self.method == "Abandon" {
		name := data["title"]
		self.AddLog(fmt.Sprintf("废弃文章:%s", name))
	} else {
		self.AddLog(fmt.Sprintf("%+v", data))
	}
}

//废弃
func (self *PostTypeController) Abandon() {
	self.EditCommon(self)
}

func (self *PostTypeController) Add() {
	self.AddCommon(self)
}

func (self *PostTypeController) Edit() {
	self.EditCommon(self)
}

func (self *PostTypeController) Del() {
	self.DelCommon(self)
}
