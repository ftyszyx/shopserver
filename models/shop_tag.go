package models

import "github.com/zyx/shop_server/libs"

type ShopTag struct {
	Model
}

type ShopTagData struct {
	Name     string `empty:"名称不能为空"`
	Pic      string
	Order_id int
}

func (self *ShopTag) InitSqlField(sql libs.SqlType) libs.SqlType {
	return sql
}

func (self *ShopTag) GetModelStruct() interface{} {
	return ShopTagData{}
}
