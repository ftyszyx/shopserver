package admin

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
)

type ConfigController struct {
	BaseController
}

func (self *ConfigController) Edit() {

	err := self.dboper.Begin()

	modelcheck := self.model.GetModelStruct()
	changelist := self.postdata["list"].([]interface{})
	for _, value := range changelist {
		mapvalue := value.(map[string]interface{})
		id := mapvalue["id"].(string)
		changedata := libs.ClearMapByStruct(mapvalue, modelcheck)
		_, err = self.dboper.Raw(fmt.Sprintf("update %s set %s where `id`=?", self.model.TableName(), db.SqlGetKeyValue(changedata, "=")), id).Exec()
	}

	if err == nil {
		err = self.dboper.Commit()
		self.AddLog(fmt.Sprintf("postdata:%+v ", self.postdata))
		self.AjaxReturnSuccess("", nil)
		return
	} else {
		err = self.dboper.Rollback()
		self.AjaxReturnError(errors.WithStack(err))
	}

}
