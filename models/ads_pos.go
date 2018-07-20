package models

import "github.com/zyx/shop_server/libs"

type AdsPos struct {
	Model
}

type AdsPosData struct {
	Name      string `empty:"名称不能为空"`
	title     string
	title_pic string
}

func (self *AdsPos) InitSqlField(sql libs.SqlType) libs.SqlType {
	return sql
}

func (self *AdsPos) GetModelStruct() interface{} {
	return AdsPosData{}
}
