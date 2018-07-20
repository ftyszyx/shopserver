package admin

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

type UserGroupController struct {
	BaseController
}

func (self *UserGroupController) BeforeSql(data map[string]interface{}) {
	logs.Info("before sql:%s", self.method)
	if self.method == "Del" {
		grouptype := data["group_type"].(string)
		if grouptype == strconv.Itoa(libs.UserSystem) {
			self.AjaxReturnError("系统用户组不可删")
		}
	} else if self.method == "Add" {
		data["group_type"] = libs.UserAdmin
	}
}
func (self *UserGroupController) AfterSql(data map[string]interface{}, oldinfo orm.Params) {

	if self.method == "Add" {
		self.AddLog(fmt.Sprintf("增加用户组:%+v", data))
	} else if self.method == "Edit" {
		id := self.postdata["id"].(string)
		self.model.ClearRowCache(id)
		self.AddLog(fmt.Sprintf("修改用户组:%+v", data))
	} else if self.method == "Del" {
		id := self.postdata["id"].(string)
		name := data["name"]
		self.model.ClearRowCache(id)
		self.AddLog(fmt.Sprintf("删除用户组:%s", name))
	} else {
		self.AddLog(fmt.Sprintf("%+v", data))
	}
}

func (self *UserGroupController) Add() {
	self.AddCommon(self)
}

func (self *UserGroupController) Edit() {
	self.EditCommon(self)
}

func (self *UserGroupController) Del() {
	self.DelCommon(self)
}
