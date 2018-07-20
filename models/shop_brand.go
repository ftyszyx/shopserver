package models

import "github.com/zyx/shop_server/libs"

type ShopBrand struct {
	Model
}

type ShopBrandData struct {
	Name     string `empty:"名称不能为空"`
	Order_id int
	Pic      string
}

func (self *ShopBrand) InitSqlField(sql libs.SqlType) libs.SqlType {
	return sql
}

func (self *ShopBrand) GetModelStruct() interface{} {
	return ShopBrandData{}
}
