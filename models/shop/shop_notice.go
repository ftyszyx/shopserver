package shop

import (
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type ShopNotice struct {
	models.Model
}

type ShopNoticeData struct {
	Title    string `empty:"标题不能为空"`
	Content  string `empty:"内容不能为空"`
	Order_id int
}

func (self *ShopNotice) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *ShopNotice) GetModelStruct() interface{} {
	return ShopNoticeData{}
}

func (self *ShopNotice) Init() {
	self.resetCache()

}

func (self *ShopNotice) resetCache() {

	self.Cache().Put("all_valid_notice", self.GetInfoByField(db.NewOper(), "is_del", 0), 0)

}

func (self *ShopNotice) ClearCache() {
	self.resetCache()
}
