package models

import "github.com/zyx/shop_server/libs/db"

type Export struct {
	Model
}

type ExportData struct {
	Name  string `empty:"名称不能为空"`
	Value string `empty:"内容不能为空"`
	Model string `empty:"模块不能为空"`
}

func (self *Export) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *Export) GetModelStruct() interface{} {
	return ExportData{}
}
