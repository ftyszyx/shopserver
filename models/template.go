package models

import "github.com/zyx/shop_server/libs"

type Template struct {
	Model
}

type TemplateData struct {
	Name string `empty:"名称不能为空"`
}

func (self *Template) InitSqlField(sql libs.SqlType) libs.SqlType {
	return sql
}

func (self *Template) GetModelStruct() interface{} {
	return TemplateData{}
}
