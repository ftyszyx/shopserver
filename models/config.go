package models

import (
	"fmt"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

type Config struct {
	Model
}

type ConfigData struct {
	value string
}

func (self *Config) InitSqlField(sql libs.SqlType) libs.SqlType {
	return sql
}

func (self *Config) GetModelStruct() interface{} {
	return ConfigData{}
}

var ConfigCache cache.Cache

func (self *Config) Init() {
	ConfigCache, _ = cache.NewCache("memory", `{"interval":0}`) //不过期
	self.resetCache()

}

func (self *Config) resetCache() {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw(fmt.Sprintf("select * from %s", self.tablename)).Values(&maps)
	if err == nil && num > 0 {
		for _, v := range maps {
			ConfigCache.Put(v["name"].(string), v["value"], 0)
		}
	}
}

func (self *Config) UpdateCache() {
	ConfigCache.ClearAll()
	self.resetCache()
}

func (self *Config) ClearCache() {
	self.cache = make(map[string]orm.Params)
	self.UpdateCache()
}
