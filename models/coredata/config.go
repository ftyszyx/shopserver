package coredata

import (
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type Config struct {
	models.Model
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

func (self *Config) Init() {
	self.resetCache()
}

func (self *Config) resetCache() {
	var maps []db.Params
	cache := self.Cache()
	cache.ClearAll()
	num, err := db.NewOper().Raw(fmt.Sprintf("select * from %s", self.TableName())).Values(&maps)
	if err == nil && num > 0 {
		for _, v := range maps {
			logs.Info("config set:%s=%s", v["name"].(string), v["value"])
			cache.Put(v["name"].(string), v["value"], 0)
		}
	}

}

func (self *Config) ClearCache() {
	self.resetCache()
}
