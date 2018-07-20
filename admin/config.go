package admin

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

type ConfigController struct {
	BaseController
}

func (self *ConfigController) Edit() {
	o := orm.NewOrm()
	err := o.Begin()

	modelcheck := self.model.GetModelStruct()
	changelist := self.postdata["list"].([]interface{})
	for _, value := range changelist {
		mapvalue := value.(map[string]interface{})
		id := mapvalue["id"].(string)
		changedata := libs.ClearMapByStruct(mapvalue, modelcheck)
		_, err = o.Raw(fmt.Sprintf("update %s set %s where `id`=?", self.model.TableName(), libs.SqlGetKeyValue(changedata, "=")), id).Exec()
	}
	if err == nil {
		err = o.Commit()
		self.AjaxReturnSuccess("", nil)
		return
	} else {
		err = o.Rollback()
		self.AjaxReturnError(err.Error())
	}

}
