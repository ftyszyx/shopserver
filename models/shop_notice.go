package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

type ShopNotice struct {
	Model
}

type ShopNoticeData struct {
	Title    string `empty:"标题不能为空"`
	Content  string `empty:"内容不能为空"`
	Order_id int
}

func (self *ShopNotice) InitSqlField(sql libs.SqlType) libs.SqlType {
	return sql
}

func (self *ShopNotice) GetModelStruct() interface{} {
	return ShopNoticeData{}
}

var NoticeCache []orm.Params //主页广告缓存

func (self *ShopNotice) Init() {
	self.resetCache()

}

func (self *ShopNotice) resetCache() {

	NoticeCache = self.GetInfoByField("is_del", 0)

}

func (self *ShopNotice) UpdateCache() {
	self.resetCache()
}
