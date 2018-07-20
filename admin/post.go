package admin

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type PostController struct {
	BaseController
}

func (self *PostController) BeforeSql(data map[string]interface{}) {
	if self.method == "Abandon" {
		data["is_del"] = self.postdata["is_del"]
	} else if self.method == "Add" {
		data["is_del"] = 0
		data["build_user"] = self.GetUid()
		data["build_time"] = time.Now().Unix()
	}
}
func (self *PostController) AfterSql(data map[string]interface{}, oldinfo orm.Params) {
	if self.method == "Add" {
		self.AddLog(fmt.Sprintf("增加文章:%+v", data))
	} else if self.method == "Edit" {
		self.AddLog(fmt.Sprintf("修改文章:%+v", data))
	} else if self.method == "Del" {
		name := data["title"]
		self.AddLog(fmt.Sprintf("删除文章:%s", name))
	} else if self.method == "Abandon" {
		name := data["title"]
		self.AddLog(fmt.Sprintf("废弃:%s", name))
	} else {
		self.AddLog(fmt.Sprintf("%+v", data))
	}
}

//废弃
func (self *PostController) Abandon() {
	self.EditCommon(self)
}

func (self *PostController) Add() {
	self.AddCommon(self)
}

func (self *PostController) Edit() {
	self.EditCommon(self)
}

func (self *PostController) Del() {
	self.DelCommon(self)
}
