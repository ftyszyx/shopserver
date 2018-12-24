package models

import "github.com/zyx/shop_server/libs/db"

type ShopBrand struct {
	Model
}

type ShopBrandData struct {
	Name     string `empty:"名称不能为空"`
	Order_id int
	Pic      string
}

func (self *ShopBrand) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *ShopBrand) GetModelStruct() interface{} {
	return ShopBrandData{}
}
