package models

import (
	"fmt"

	"github.com/zyx/shop_server/libs"
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
	idnum          string
	order_vip_type int
	total_price    string
	client_provice string `empty:"收获地址省份空"`
	client_city    string `empty:"收获地址城市空"`
	client_area    string `empty:"收获地址地区空"`
}

func (self *ShopOrder) InitSqlField(sql libs.SqlType) libs.SqlType {

	return self.InitField(self.InitJoinString(sql, true))

}

func (self *ShopOrder) GetModelStruct() interface{} {
	return ShopOrderData{}
}

func (self *ShopOrder) ExportNameProcess(name string, value string) string {
	return value
}
func (self *ShopOrder) InitJoinString(sql libs.SqlType, allfield bool) libs.SqlType {
	memberTable := GetModel(USER).TableName()
	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("member") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `member` ON `member`.`id`=`order`.`user_id`", memberTable)
	}

	return sql.Alias("order").Join(fieldstr)
}
func (self *ShopOrder) InitField(sql libs.SqlType) libs.SqlType {
	return sql.Field(map[string]string{
		"order.id":             "id",
		"order.item_info":      "item_info",
		"order.order_time":     "order_time",
		"order.pay_time":       "pay_time",
		"order.status":         "status",
		"order.client_address": "client_address",
		"order.client_name":    "client_name",
		"order.user_id":        "user_id",
		"member.name":          "member_name",
		"order.idnum":          "idnum",
		"order.client_phone":   "client_phone",
		"order.refund_info":    "refund_info",
		"order.close_info":     "close_info",
		"order.close_type":     "close_type",
		"order.shipment_num":   "shipment_num",
		"order.client_info":    "client_info",
		"order.sell_info":      "sell_info",
		"order.close_time":     "close_time",
		"order.idnumpic1":      "idnumpic1",
		"order.idnumpic2":      "idnumpic2",
		"order.total_price":    "total_price",
		"order.client_provice": "client_provice",
		"order.client_city":    "client_city",
		"order.order_vip_type": "order_vip_type",
		"order.client_area":    "client_area",
		"order.pay_id":         "pay_id",
	})
}
