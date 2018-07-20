package models

import (
	"fmt"

	"github.com/zyx/shop_server/libs"
)

type ShopItem struct {
	Model
}

type ShopItemData struct {
	Code      string `empty:"编码不能为空" only:"编码重复"`
	Name      string `empty:"名称不能为空" only:"名称重复"`
	Item_type string `empty:"类型不能为空"`
	Weight    float64
	is_onsale int
	Order_id  int
	Pics      string
	icon      string
	Price     float64
	Store_num int
	Brand     string `empty:"品牌不能为空"`
	Spec      string
	Tag       string
	Desc      string
}

func (self *ShopItem) InitSqlField(sql libs.SqlType) libs.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *ShopItem) GetModelStruct() interface{} {
	return ShopItemData{}
}
func (self *ShopItem) InitJoinString(sql libs.SqlType, allfield bool) libs.SqlType {
	itemtypename := GetModel(SHOP_ITEMTYPE).TableName()
	brandname := GetModel(SHOP_BRAND).TableName()
	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("item_type") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `item_type` ON `item_type`.`id`=`item`.`item_type`", itemtypename)
	}
	if (allfield == true) || (sql.NeedJointable("brand") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `brand` ON `brand`.`id`=`item`.`brand`", brandname)
	}
	return sql.Alias("item").Join(fieldstr)
}

func (self *ShopItem) InitField(sql libs.SqlType) libs.SqlType {
	return sql.Field(map[string]string{
		"item.id":        "id",
		"item.name":      "name",
		"item.pics":      "pics",
		"item.icon":      "icon",
		"item.price":     "price",
		"item.store_num": "store_num",
		"item.sell_num":  "sell_num",
		"item.is_onsale": "is_onsale",
		"item.order_id":  "order_id",
		"item.item_type": "item_type",
		"item_type.name": "item_type_name",
		"item.brand":     "brand",
		"brand.name":     "brand_name",
		"item.spec":      "spec",
		"item.desc":      "desc",
		"item.code":      "code",
		"item.tag":       "tag",
	})
}
