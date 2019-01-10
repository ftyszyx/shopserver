package coredata

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type Module struct {
	models.Model
}

func (self *Module) Init() {
	self.resetCache()
}

func (self *Module) resetCache() {
	var maps []db.Params
	num, err := db.NewOper().Raw(fmt.Sprintf("select * from %s", self.TableName())).Values(&maps)
	cache := self.Cache()
	cache.ClearAll()
	if err == nil && num > 0 {
		for _, v := range maps {
			v["name"] = ""
			control, ok := v["controller"].(string)
			if ok == false {
				control = ""
			}
			method, ok := v["method"].(string)
			if ok == false {
				method = ""
			}
			control = strings.ToLower(control)
			method = strings.ToLower(method)
			key := fmt.Sprintf("%v-%v", control, method)
			cache.Put(key, v, 0)
		}
		cache.Put("allmodel", maps, 0)
	}

}

//获取module
func (self *Module) GetModuleInfo(control string, method string) db.Params {
	control = strings.ToLower(control)
	method = strings.ToLower(method)
	cache := self.Cache()
	moduleinfo, ok := cache.Get(control + "-" + method).(db.Params)
	if ok {
		return moduleinfo
	}
	return nil
}

//判断模块要不要验证
func (self *Module) NeedAuth(control string, method string) bool {
	moduleinfo := self.GetModuleInfo(control, method)
	logs.Info("modleinfo:%+v", moduleinfo)
	if moduleinfo == nil || moduleinfo["need_auth"] == "0" {
		return false
	}
	return true
}

func (self *Module) ClearCache() {

	self.resetCache()
}
