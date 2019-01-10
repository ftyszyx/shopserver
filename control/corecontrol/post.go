package corecontrol

import (
	"fmt"
	"time"

	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs/db"
)

type PostController struct {
	base.BaseController
}

func (self *PostController) BeforeSql(data map[string]interface{}) error {
	if self.GetMethod() == "Abandon" {
		data["is_del"] = self.GetPost()["is_del"]
	} else if self.GetMethod() == "Add" {
		data["is_del"] = 0
		data["build_user"] = self.GetUid()
		data["build_time"] = time.Now().Unix()
	}
	return nil
}
func (self *PostController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.GetMethod() == "Add" {
		self.AddLog(fmt.Sprintf("增加文章:%+v", data))
	} else if self.GetMethod() == "Edit" {
		self.AddLog(fmt.Sprintf("修改文章:%+v", data))
	} else if self.GetMethod() == "Del" {
		name := data["title"]
		self.AddLog(fmt.Sprintf("删除文章:%s", name))
	} else if self.GetMethod() == "Abandon" {
		name := data["title"]
		self.AddLog(fmt.Sprintf("废弃:%s", name))
	}
	return nil
}

//废弃
func (self *PostController) Abandon() {
	self.EditCommonAndReturn(self)
}

func (self *PostController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *PostController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *PostController) Del() {
	self.DelCommonAndReturn(self)
}
