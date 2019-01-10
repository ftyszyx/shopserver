package base

import (
	"fmt"

	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

func (self *BaseController) Logcommon(data map[string]interface{}, oldinfo db.Params) error {
	if self.method == "Add" {
		self.AddLog(fmt.Sprintf("add data:%+v ", data))
	} else if self.method == "Edit" {
		self.AddLog(fmt.Sprintf("change data:%+v ", data))
	} else if self.method == "Del" {
		self.AddLog(fmt.Sprintf(" oldinfo:%+v ", data))
	}
	return nil
}

//增加日志
func (self *BaseController) AddLog(info string) {
	models.AddLog(self.dboper, self.uid, info, self.control, self.method)
}

func (self *BaseController) AddLogLink(info string, link string) {
	models.AddLogLink(self.dboper, self.uid, info, self.control, self.method, link)
}
