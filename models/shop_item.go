package models

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs/db"
)

type ShopItem struct {
	Model
}

type ShopItemData struct {
	Code            string `empty:"编码不能为空"`
	Name            string `empty:"名称不能为空" only:"名称重复"`
	Item_type       string `empty:"类型不能为空"`
	Weight          float64
	is_onsale       int
	Order_id        int
	Pics            string
	icon            string
	Price           float64
	Store_num       int
	basenum         int
	Brand           string `empty:"品牌不能为空"`
	Spec            string
	Tag             string
	Group_price     string
	Desc            string
	no_service      string
	Item_unit       string
	Item_shelf_life string
	Is_sync_shipnum int
	Idnum_need      string
	Supply_source   string
	Min_num         int
}

//商品规格结构
type ItemSpecTagData struct {
	Id   int
	Name string
}

type ItemSpecItemData struct {
	Name string
	Id   int
	List []ItemSpecTagData
}

type ItemSpecDetailIdData struct {
	Specid int
	Tagid  int
}

type ItemPriceGroupinfo struct {
	Groupid string
	Price   string
}

type ItemSpecDetailData struct {
	Namearr     []ItemSpecDetailIdData
	Price       string
	Store_num   string
	Code        string
	Pic         string
	Group_price []ItemPriceGroupinfo
}

type ItemSpecData struct {
	DetailList []ItemSpecDetailData
	SpecList   []ItemSpecItemData
}

func (self *ShopItem) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *ShopItem) GetModelStruct() interface{} {
	return ShopItemData{}
}
func (self *ShopItem) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
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

func (self *ShopItem) GetItemPrice(groupid string, iteminfo db.Params, code string) (string, error) {
	specinfo := iteminfo["spec"]
	var specdata *ItemSpecData
	//看规格
	if specinfo != nil {
		specinfostr := specinfo.(string)
		if specinfostr != "" {
			specdata = new(ItemSpecData)
			logs.Info("specstr:%s", specinfostr)
			err := json.Unmarshal([]byte(specinfostr), specdata)
			if err != nil {
				return "", err
			}
			for _, detailitem := range specdata.DetailList {
				if detailitem.Code == code {
					for _, groupitem := range detailitem.Group_price {
						if groupitem.Groupid == groupid {
							return groupitem.Price, nil
						}
					}
					return detailitem.Price, nil
				}
			}
		}

	}

	defaultprice := iteminfo["price"]
	groupinfo := iteminfo["group_price"]
	if groupinfo != nil {
		groupstr := groupinfo.(string)
		if groupstr != "" {
			defaultGroupinfo := new([]ItemPriceGroupinfo)
			err := json.Unmarshal([]byte(groupstr), defaultGroupinfo)
			if err != nil {
				return "", err
			}
			for _, groupitem := range *defaultGroupinfo {
				if groupitem.Groupid == groupid {
					return groupitem.Price, nil
				}
			}
		}
	}
	return defaultprice.(string), nil
}

func (self *ShopItem) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"item.id":              "id",
		"item.name":            "name",
		"item.pics":            "pics",
		"item.icon":            "icon",
		"item.price":           "price",
		"item.store_num":       "store_num",
		"item.item_unit":       "item_unit",
		"item.item_shelf_life": "item_shelf_life",
		"item.sell_num":        "sell_num",
		"item.is_onsale":       "is_onsale",
		"item.order_id":        "order_id",
		"item.item_type":       "item_type",
		"item_type.name":       "item_type_name",
		"item.brand":           "brand",
		"brand.name":           "brand_name",
		"item.spec":            "spec",
		"item.is_sync_shipnum": "is_sync_shipnum",
		"item.basenum":         "basenum",
		"item.desc":            "desc",
		"item.group_price":     "group_price",
		"item.code":            "code",
		"item.no_service":      "no_service",
		"item.tag":             "tag",
		"item.supply_source":   "supply_source",
		"item.idnum_need":      "idnum_need",
		"item.min_num":         "min_num",
	})
}
