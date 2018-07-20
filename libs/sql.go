package libs

import (
	"bytes"
	"fmt"
	"strings"
)

/*
order":{"sort_id":"asc"}
"search":{"name":["LIKE","tt"]}
"search":{"name":["LIKE","tt"],"is_del":0}}
{"pay_time":[[">",1525190400],["<",1525276800],"and"]}}

SELECT COUNT(*) AS tp_count FROM `aq_sell` `sell` LEFT JOIN `aq_store` `store` ON `sell`.`store_id`=`store`.`id`
LEFT JOIN `aq_sys_user` `check_user` ON `sell`.`check_user`=`check_user`.`id`
LEFT JOIN `aq_sys_user` `build_user` ON `sell`.`build_user`=`build_user`.`id` LEFT JOIN `aq_shop` `shop` ON `sell`.`shop_id`=`shop`.`id`
LEFT JOIN `aq_item` `item` ON `sell`.`item_id`=`item`.`id`
LEFT JOIN `aq_item_type` `item_type` ON `item`.`type`=`item_type`.`id`
LEFT JOIN `aq_store_item` `store_item` ON `sell`.`store_id`=store_item.store_id and sell.item_id=store_item.itemid
WHERE  `sell`.`id` = 'saas'  AND ( `pay_time` > 1525276800 and `pay_time` < 1528214400 )  AND `item`.`short_name` LIKE '%dsafd%' LIMIT 1

SELECT `sell`.`id` AS `id`,trim(sell.id) AS `id_like`,`sell`.`status` AS `status`,`sell`.`refund_status` AS `refund_status`,`sell`.`pay_status` AS `pay_status`,`sell`.`del_info` AS `del_info`,`sell`.`shop_id` AS `shop_id`,`shop`.`name` AS `shop_name`,`sell`.`shop_order` AS `shop_order`,`sell`.`customer_name` AS `customer_name`,`sell`.`customer_addr` AS `customer_addr`,`sell`.`info` AS `info`,`sell`.`user_info` AS `user_info`,`sell`.`pay_time` AS `pay_time`,`sell`.`discount` AS `discount`,`sell`.`user_phone` AS `user_phone`,`sell`.`pay_money` AS `pay_money`,`sell`.`user_id_number` AS `user_id_number`,`sell`.`customer_account` AS `customer_account`,`sell`.`store_id` AS `store_id`,`store`.`name` AS `store_name`,`sell`.`build_time` AS `build_time`,`sell`.`build_user` AS `build_user`,`build_user`.`name` AS `build_user_name`,`sell`.`check_user` AS `check_user`,`sell`.`check_time` AS `check_time`,`check_user`.`name` AS `check_user_name`,`sell`.`assign_order` AS `assign_order`,`sell`.`item_id` AS `item_id`,`sell`.`num` AS `num`,`item`.`name` AS `item_name`,`item`.`short_name` AS `item_short_name`,`item`.`sort_id` AS `item_sort_id`,`item_type`.`sort_id` AS `item_type_sort_id`,`sell`.`sell_type` AS `sell_type`,`sell`.`logistics` AS `logistics`,`sell`.`track_man` AS `track_man`,`sell`.`sell_vip_type` AS `sell_vip_type`,`sell`.`sell_vip_info` AS `sell_vip_info`,`sell`.`logistics_merge` AS `logistics_merge`,`item`.`milk_period` AS `item_milk_period`,`sell`.`order_time` AS `order_time`,`sell`.`customer_province` AS `customer_province`,`sell`.`customer_city` AS `customer_city`,`sell`.`customer_area` AS `customer_area`,`sell`.`send_user_name` AS `send_user_name`,`sell`.`send_user_phone` AS `send_user_phone`,`sell`.`unit_price` AS `unit_price`,`sell`.`freight_price` AS `freight_price`,`sell`.`service_price` AS `service_price`,`sell`.`freight_unit_price` AS `freight_unit_price`,`sell`.`service_unit_price` AS `service_unit_price`,`store_item`.`in_store` AS `in_store_num`
 FROM `aq_sell` `sell`
LEFT JOIN `aq_store` `store` ON `sell`.`store_id`=`store`.`id`
 LEFT JOIN `aq_sys_user` `check_user` ON `sell`.`check_user`=`check_user`.`id`
 LEFT JOIN `aq_sys_user` `build_user` ON `sell`.`build_user`=`build_user`.`id` LEFT JOIN `aq_shop` `shop` ON `sell`.`shop_id`=`shop`.`id`
 LEFT JOIN `aq_item` `item` ON `sell`.`item_id`=`item`.`id`
LEFT JOIN `aq_item_type` `item_type` ON `item`.`type`=`item_type`.`id`
LEFT JOIN `aq_store_item` `store_item` ON `sell`.`store_id`=store_item.store_id and sell.item_id=store_item.itemid
WHERE  `sell`.`id` = 'saas'  AND ( `pay_time` > 1525276800 and `pay_time` < 1528214400 )  AND `item`.`short_name` LIKE '%dsafd%'
ORDER BY `item_sort_id`  asc LIMIT 0,50
*/

type SqlType interface {
	getSql() string
	Select() string
	Find() string
	Count() string
	Name(string) SqlType
	Where(data map[string]interface{}) SqlType
	WhereOr(data map[string]interface{}) SqlType
	Field(map[string]string) SqlType
	Alias(string) SqlType
	Order(map[string]interface{}) SqlType
	Limit([]int) SqlType
	Join(string) SqlType
	GetJoinStr() string
	GetAlias() string

	NeedJointable(tablename string) bool
}

type SqlBuild struct {
	TableName string
	field     map[string]string      //要取的字段
	where     map[string]interface{} //查找语句
	alias     string                 //表名别名
	getType   string                 // find  count select 类型
	limit     []int                  //limit  [3,100]  start num
	andOr     bool                   // true->and false ->or
	order     map[string]interface{} //排序
	joinStr   string                 //
}

func (self *SqlBuild) GetJoinStr() string {
	return self.joinStr
}

func (self *SqlBuild) GetAlias() string {
	return self.alias
}

//是否要内联table
func (self *SqlBuild) NeedJointable(tablename string) bool {
	for key, _ := range self.where {
		strarr := strings.Split(key, ".")
		if len(strarr) > 1 && strarr[0] == tablename {
			return true
		}
	}
	for key, _ := range self.order {
		strarr := strings.Split(key, ".")
		if len(strarr) > 1 && strarr[0] == tablename {
			return true
		}
	}
	return false
}

func (self *SqlBuild) getSql() string {
	var connectstr = "And"
	var whereStr = ""
	var orderStr = ""
	var fieldStr = " * "
	var tableAlias = ""
	var limitStr = ""

	if self.andOr == false {
		connectstr = "OR"
	}
	if self.where != nil {

		whereStr = strings.TrimSpace(SqlGetSearch(self.where, connectstr))
		if whereStr != "" {
			whereStr = "where " + whereStr
		}
	}
	if self.order != nil {
		orderStr = strings.TrimSpace(SqlGetKeys(self.order, " "))
		if orderStr != "" {
			orderStr = "ORDER BY " + orderStr
		}
	}
	if self.field != nil {
		fieldStr = strings.TrimSpace(SqlGetField(self.field))
	}
	if self.alias != "" {
		tableAlias = "`" + self.alias + "`"
	}
	if self.getType == "find" {
		limitStr = " limit 1 "
	} else if self.getType == "select" {
		if len(self.limit) == 2 {
			limitStr = fmt.Sprintf(" limit %d,%d ", self.limit[0], self.limit[1])
		}
	}
	if self.getType == "count" {
		countfield := "count(*) as " + SQL_COUNT_NAME
		return strings.TrimSpace(fmt.Sprintf("select %s from `%s` %s %s %s ", countfield, self.TableName, tableAlias, self.joinStr, whereStr))
	} else if self.getType == "select" || self.getType == "find" {
		return strings.TrimSpace(fmt.Sprintf("select %s from `%s` %s %s %s %s %s", fieldStr, self.TableName, tableAlias, self.joinStr, whereStr, orderStr, limitStr))
	}
	return ""
}

func (self *SqlBuild) Select() string {
	self.getType = "select"
	return self.getSql()
}

func (self *SqlBuild) Find() string {
	self.getType = "find"
	return self.getSql()
}

func (self *SqlBuild) Count() string {
	self.getType = "count"
	return self.getSql()
}

func (self *SqlBuild) Where(data map[string]interface{}) SqlType {
	if data == nil {
		return self
	}
	if self.where == nil {
		self.where = make(map[string]interface{})
	}
	for key, value := range data {
		self.where[key] = value
	}
	self.andOr = true
	return self
}

func (self *SqlBuild) WhereOr(data map[string]interface{}) SqlType {
	if data == nil {
		return self
	}
	if self.where == nil {
		self.where = make(map[string]interface{})
	}
	for key, value := range data {
		self.where[key] = value
	}
	self.andOr = false
	return self
}

func (self *SqlBuild) Field(data map[string]string) SqlType {
	self.field = data
	return self
}

func (self *SqlBuild) Alias(name string) SqlType {
	self.alias = name
	return self
}

func (self *SqlBuild) Name(name string) SqlType {
	self.TableName = name
	return self
}

func (self *SqlBuild) Order(data map[string]interface{}) SqlType {
	self.order = data
	return self
}

func (self *SqlBuild) Limit(data []int) SqlType {
	self.limit = data
	return self
}

func (self *SqlBuild) Join(data string) SqlType {
	self.joinStr = data
	return self
}

func SqlGetKeyValue(arr map[string]interface{}, connect string) string {
	var buffer bytes.Buffer
	for key, value := range arr {

		buffer.WriteString(" ")
		buffer.WriteString(SqlGetKey(key))
		buffer.WriteString(connect)
		buffer.WriteString(" ")
		buffer.WriteString(SqlGetString(value))
		buffer.WriteString(",")
	}
	return strings.Trim(buffer.String(), ",")
}

func SqlGetKeys(arr map[string]interface{}, connect string) string {
	var buffer bytes.Buffer
	for key, value := range arr {

		buffer.WriteString(" ")
		buffer.WriteString(SqlGetKey(key))
		buffer.WriteString(connect)
		buffer.WriteString(" ")
		buffer.WriteString(SqlEscap(value.(string)))
		buffer.WriteString(",")

	}
	return strings.Trim(buffer.String(), ",")
}

func SqlGetString(value interface{}) string {
	// logs.Info("type:%s", reflect.TypeOf(value))
	// logs.Info("value:%v", value)

	if temp, ok := value.(string); ok {
		return "'" + SqlEscap(temp) + "'"
	} else {
		return "'" + fmt.Sprintf("%v", value) + "'"
	}
	//return ""
}
func SqlGetSearch(search map[string]interface{}, andstr string) string {
	var buffer bytes.Buffer
	for key, value := range search {
		buffer.WriteString(" ")
		if valuetemp, ok := value.([]interface{}); ok {
			if len(valuetemp) == 2 {
				//( like aa)
				buffer.WriteString(SqlGetKey(key))
				buffer.WriteString(" ")
				buffer.WriteString(SqlEscap(valuetemp[0].(string)))
				buffer.WriteString(" ")
				buffer.WriteString(SqlGetString(valuetemp[1]))
			} else {
				//( `pay_time` > 1525276800 and `pay_time` < 1528214400 )

				buffer.WriteString("(")
				var condarrstr []string
				var condconect = "and"
				for _, cond := range valuetemp {
					strtemp, ok := cond.(string)
					if ok == false {
						condarr := cond.([]interface{})
						var buffertemp bytes.Buffer
						buffertemp.WriteString(SqlGetKey(key))
						buffertemp.WriteString(" ")
						buffertemp.WriteString(SqlEscap(condarr[0].(string)))
						buffertemp.WriteString(" ")
						buffertemp.WriteString(SqlGetString(condarr[1]))
						condarrstr = append(condarrstr, buffertemp.String())
					} else {
						condconect = strtemp
					}
				}
				//logs.Info("arr %v %d", condarrstr, len(condarrstr))
				buffer.WriteString(strings.Join(condarrstr, " "+condconect+" "))

				buffer.WriteString(")")
			}
		} else {
			buffer.WriteString(SqlGetKey(key))
			buffer.WriteString("=")
			buffer.WriteString(SqlGetString(value))
		}
		buffer.WriteString(" ")
		buffer.WriteString(andstr)
	}
	return strings.Trim(buffer.String(), andstr)
}

func SqlGetKey(key string) string {
	var buffer bytes.Buffer
	if strings.Index(key, "(") > 0 { //处理trim(sell.id)
		buffer.WriteString(SqlEscap(key))
	} else {
		arr := strings.Split(key, ".")
		for i := 0; i < len(arr); i++ {
			if i > 0 {
				buffer.WriteString(".")
			}
			buffer.WriteString("`")
			buffer.WriteString(SqlEscap(arr[i]))
			buffer.WriteString("`")
		}
	}
	return buffer.String()
}

//获取要取的字段
func SqlGetField(fields map[string]string) string {
	var buffer bytes.Buffer
	for key, value := range fields {
		//写as左边
		buffer.WriteString(SqlGetKey(key))
		if strings.TrimSpace(value) != "" {
			buffer.WriteString(" AS ")
			//as 右边
			buffer.WriteString("`")
			buffer.WriteString(SqlEscap(value))
			buffer.WriteString("`")
		}
		buffer.WriteString(",")
	}
	//去掉最后一个,
	return strings.Trim(buffer.String(), ",")
}

func SqlGetInsertInfo(arr map[string]interface{}) (string, string) {
	var bufferkey bytes.Buffer
	var buffervalue bytes.Buffer
	for key, value := range arr {
		bufferkey.WriteString(SqlGetKey(key))
		bufferkey.WriteString(",")
		buffervalue.WriteString(SqlGetString(value))
		buffervalue.WriteString(",")
	}
	return strings.Trim(bufferkey.String(), ","), strings.Trim(buffervalue.String(), ",")
}

func SqlGetArrInfo(arr []interface{}) string {
	var buffervalue bytes.Buffer
	for _, value := range arr {
		buffervalue.WriteString(SqlGetString(value))
		buffervalue.WriteString(",")
	}
	return "(" + strings.Trim(buffervalue.String(), ",") + ")"
}

func SqlEscap(src string) string {
	return strings.Replace(src, "'", "\\'", -1)
}
