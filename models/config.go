package models

import (
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/cache"
	"github.com/zyx/shop_server/libs/db"
)

type Config struct {
	Model
}

type ConfigData struct {
	value string
}

func (self *Config) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *Config) GetModelStruct() interface{} {
	return ConfigData{}
}

var ConfigCache cache.Cache

func (self *Config) Init() {
	logs.Info("config init")
	ConfigCache, _ = cache.NewCache("memory", `{"interval":0}`) //不过期
	self.resetCache()

}

func (self *Config) resetCache() {
	// o := orm.NewOrm()
	var maps []db.Params
	num, err := db.NewOper().Raw(fmt.Sprintf("select * from %s", self.tablename)).Values(&maps)
	if err == nil && num > 0 {
		for _, v := range maps {
			logs.Info("config set:%s=%s", v["name"].(string), v["value"])
			ConfigCache.Put(v["name"].(string), v["value"], 0)
		}
	}

}

func (self *Config) UpdateCache() {
	ConfigCache.ClearAll()
	self.resetCache()
}

func (self *Config) ClearCache() {
	self.cache = make(map[string]db.Params)
	self.UpdateCache()
}
