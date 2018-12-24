package models

import (
	"github.com/zyx/shop_server/libs/db"
)

type ShopNotice struct {
	Model
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

var NoticeCache []db.Params //主页广告缓存

func (self *ShopNotice) Init() {
	self.resetCache()

}

func (self *ShopNotice) resetCache() {

	NoticeCache = self.GetInfoByField(db.NewOper(), "is_del", 0)

}

func (self *ShopNotice) UpdateCache() {
	self.resetCache()
}
