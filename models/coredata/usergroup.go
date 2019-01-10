package coredata

import "github.com/zyx/shop_server/models"

type UserGroup struct {
	models.Model
}

type UserGroupData struct {
	Name             string `empty:"组名不能为空"`
	Module_ids       string
	expire_time      int
	Group_type       int
	Limit_show_order string
}

func (self *UserGroup) GetModelStruct() interface{} {
	return UserGroupData{}
}
