package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs/db"
)

type Module struct {
	Model
}

var ModuleCache cache.Cache //主页广告缓存

//初始化所有module
func (self *Module) Init() {
	ModuleCache, _ = cache.NewCache("memory", `{"interval":0}`) //不过期
	self.InitModuleCache()
}

func (self *Module) InitModuleCache() {
	// o := orm.NewOrm()
	var maps []db.Params
	num, err := db.NewOper().Raw(fmt.Sprintf("select * from %s", self.tablename)).Values(&maps)
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
			//self.cache[key] = v
			ModuleCache.Put(key, v, 0)
		}
	}

}

//获取module
func (self *Module) GetModuleInfo(control string, method string) db.Params {
	control = strings.ToLower(control)
	method = strings.ToLower(method)
	moduleinfo, ok := ModuleCache.Get(control + "-" + method).(db.Params)
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

func (self *Module) updateModuleCache() {
	ModuleCache.ClearAll()
	self.InitModuleCache()
}

func (self *Module) ClearCache() {
	self.cache = make(map[string]db.Params)
	self.updateModuleCache()
}
