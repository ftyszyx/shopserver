package admin

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/zyx/shop_server/libs/db"
)

type PostTypeController struct {
	BaseController
}

func (self *PostTypeController) BeforeSql(data map[string]interface{}) error {
	if self.method == "Abandon" {
		data["is_del"] = self.postdata["is_del"]
	} else if self.method == "Add" {
		data["is_del"] = 0
	} else if self.method == "Edit" {
		//logs.Info("data:%#v  parent:%+v", data, data["parent_id"])
		if data["parent_id"] != nil {
			if data["parent_id"].(string) == self.postdata["id"].(string) {
				return errors.New("父节点不能是自己")
			}
		}

	}
	return nil
}
func (self *PostTypeController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
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
	return nil
}

//废弃
func (self *PostTypeController) Abandon() {
	self.EditCommonAndReturn(self)
}

func (self *PostTypeController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *PostTypeController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *PostTypeController) Del() {
	self.DelCommonAndReturn(self)
}
