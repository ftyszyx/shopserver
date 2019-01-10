package shop

import (
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type ShopTag struct {
	models.Model
}

type ShopTagData struct {
	Name     string `empty:"名称不能为空"`
	Pic      string
	Order_id int
}

func (self *ShopTag) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *ShopTag) GetModelStruct() interface{} {
	return ShopTagData{}
}
