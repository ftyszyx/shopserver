package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type Module struct {
	Model
}

//初始化所有module
func (self *Module) Init() {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw(fmt.Sprintf("select * from %s", self.tablename)).Values(&maps)
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
			self.cache[key] = v
		}
	}

}

//获取module
func (self *Module) GetModuleInfo(control string, method string) orm.Params {
	control = strings.ToLower(control)
	method = strings.ToLower(method)
	moduleinfo, ok := self.cache[control+"-"+method]
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
