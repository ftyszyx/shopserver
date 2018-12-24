package models

import "github.com/zyx/shop_server/libs/db"

type ShopItemType struct {
	Model
}

type ShopItemTypeData struct {
	Name       string `empty:"名称不能为空"  only:"名称重复"`
	Code       string `empty:"编码不能为空"  only:"编码重复"`
	Info       string
	Level      int
	Parent_id  string
	order_id   int
	Intro_text string
	pic        string
	Is_del     int
}

func (self *ShopItemType) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *ShopItemType) GetModelStruct() interface{} {
	return ShopItemTypeData{}
}
