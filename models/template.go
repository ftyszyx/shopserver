package models

import "github.com/zyx/shop_server/libs/db"

type Template struct {
	Model
}

type TemplateData struct {
	Name string `empty:"名称不能为空"`
}

func (self *Template) InitSqlField(sql db.SqlType) db.SqlType {
	//return self.InitField(self.InitJoinString(sql, true))
	return sql
}

func (self *Template) GetModelStruct() interface{} {
	return TemplateData{}
}

// func (self *Template) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {

// }
