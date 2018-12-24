package models

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
)

type ShopOrder struct {
	Model
}

type ShopOrderData struct {
	// Item_info      string `empty:"商品信息空"`
	client_name    string `empty:"收货人名空"`
	client_address string `empty:"收获地址空"`
	client_phone   string `empty:"收件手机空"`
	client_info    string
	idnumpic1      string
	idnumpic2      string
	// shipment_num   string
	idnum           string
	order_vip_type  int
	total_price     string
	send_user_name  string
	send_user_phone string
	pay_type        int

	client_provice string `empty:"收获地址省份空"`
	client_city    string `empty:"收获地址城市空"`
	client_area    string `empty:"收获地址地区空"`
}

type ErpSelldata struct {
	Logistics string
}

func (self *ShopOrder) InitSqlField(sql db.SqlType) db.SqlType {

	return self.InitField(self.InitJoinString(sql, true))

}

func (self *ShopOrder) GetModelStruct() interface{} {
	return ShopOrderData{}
}

func (self *ShopOrder) ExportNameProcess(name string, celldata interface{}, row db.Params) (string, error) {

	if celldata == nil {
		// logs.Info("field %s is nil", name)
		return "", nil
	}
	value, ok := celldata.(string)
	if ok == false {
		return "", errors.New("upload file err:" + name + " not exit")
	}

	if name == "pay_time" || name == "order_time" || name == "close_time" {
		return libs.FormatTableTime(value), nil
	} else if name == "client_phone" || name == "idnum" {
		return "\t" + value, nil
	} else if name == "status" {
		valuenum, err := strconv.Atoi(value)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return libs.OrderStatusArr[valuenum], nil
	} else if name == "order_vip_type" {
		valuenum, err := strconv.Atoi(value)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return libs.OrderVipTypeArr[valuenum], nil
	} else if name == "close_type" {
		valuenum, err := strconv.Atoi(value)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return libs.OrderCloseTypeArr[valuenum], nil
	}
	return value, nil
}
func (self *ShopOrder) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
	memberTable := GetModel(USER).TableName()
	itemTable := GetModel(SHOP_ITEM).TableName()
	paycodetable := GetModel(PAYCODE).TableName()
	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("member") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `member` ON `member`.`id`=`order`.`user_id`", memberTable)
		if (allfield == true) || (sql.NeedJointable("track_admin_member") == true) {
			fieldstr += fmt.Sprintf("left join `%s` `track_admin_member` ON `track_admin_member`.`id`=`member`.`track_admin`", memberTable)
		}
	}

	if (allfield == true) || (sql.NeedJointable("item") == true) {
		fieldstr += fmt.Sprintf("left join `%s` `item` ON `item`.`id`=`order`.`itemid`", itemTable)
	}
	if (allfield == true) || (sql.NeedJointable("paycode") == true) {
		fieldstr += fmt.Sprintf("left join `%s` `paycode` ON `paycode`.`id`=`order`.`pay_id`", paycodetable)
	}

	return sql.Alias("order").Join(fieldstr)
}
func (self *ShopOrder) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"order.id": "id",
		// "order.item_info":       "item_info",
		"order.itemid":         "itemid",
		"order.num":            "num",
		"order.specname":       "specname",
		"order.itemcode":       "itemcode",
		"item.name":            "itemname",
		"item.is_sync_shipnum": "is_sync_shipnum",
		"order.itempic":        "itempic",
		"order.unitprice":      "unitprice",
		"paycode.money":        "paycode_money",

		"order.order_time":        "order_time",
		"order.pay_time":          "pay_time",
		"order.status":            "status",
		"order.client_address":    "client_address",
		"order.client_name":       "client_name",
		"order.user_id":           "user_id",
		"member.name":             "member_name",
		"member.account":          "member_account",
		"member.track_admin":      "member_track_admin",
		"track_admin_member.name": "member_track_admin_name",
		"order.idnum":             "idnum",
		"order.pay_type":          "pay_type",
		"order.pay_check_info":    "pay_check_info",
		"order.client_phone":      "client_phone",
		"order.refund_info":       "refund_info",
		"order.close_info":        "close_info",
		"order.close_type":        "close_type",
		"order.shipment_num":      "shipment_num",
		"order.client_info":       "client_info",
		"order.sell_info":         "sell_info",
		"order.close_time":        "close_time",
		"order.idnumpic1":         "idnumpic1",
		"order.idnumpic2":         "idnumpic2",
		"order.total_price":       "total_price",
		"order.client_provice":    "client_provice",
		"order.client_city":       "client_city",
		"order.order_vip_type":    "order_vip_type",
		"order.client_area":       "client_area",
		"order.send_user_name":    "send_user_name",
		"order.send_user_phone":   "send_user_phone",
		"order.freight_price":     "freight_price",
		"order.service_price":     "service_price",
		"order.supply_source":     "supply_source",
		"order.pay_id":            "pay_id",
	})
}
