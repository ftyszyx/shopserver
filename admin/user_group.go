package admin

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
)

type UserGroupController struct {
	BaseController
}

func (self *UserGroupController) BeforeSql(data map[string]interface{}) error {
	logs.Info("before sql:%s", self.method)
	if self.method == "Del" {
		grouptype := data["group_type"].(string)
		if grouptype == strconv.Itoa(libs.UserSystem) {

			return errors.New("系统用户组不可删")
		}
	} else if self.method == "Add" {
		if data["group_type"] == nil {
			data["group_type"] = libs.UserAdmin
		}

	}
	return nil
}
func (self *UserGroupController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {

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
	return nil
}

func (self *UserGroupController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *UserGroupController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *UserGroupController) Del() {
	self.AjaxReturnError(errors.New("不能删除用户组"))
	self.DelCommonAndReturn(self)
}
