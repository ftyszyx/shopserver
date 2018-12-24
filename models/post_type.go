package models

import "github.com/zyx/shop_server/libs/db"

type PostType struct {
	Model
}
type PostTypeData struct {
	Name      string `empty:"类型名不能为空"`
	Order_id  int
	Parent_id string
	Level     int `empty:"层级不能为空"`
}

//LEFT JOIN `aq_sys_user` `check_user` ON `sell`.`check_user`=`check_user`.`id`
func (self *PostType) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *PostType) GetModelStruct() interface{} {
	return PostTypeData{}
}
