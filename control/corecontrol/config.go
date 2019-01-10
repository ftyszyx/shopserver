package corecontrol

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
)

type ConfigController struct {
	base.BaseController
}

func (self *ConfigController) Edit() {

	err := self.GetDb().Begin()

	modelcheck := self.GetModel().GetModelStruct()
	changelist := self.GetPost()["list"].([]interface{})
	for _, value := range changelist {
		mapvalue := value.(map[string]interface{})
		id := mapvalue["id"].(string)
		changedata := libs.ClearMapByStruct(mapvalue, modelcheck)
		_, err = self.GetDb().Raw(fmt.Sprintf("update %s set %s where `id`=?", self.GetModel().TableName(), db.SqlGetKeyValue(changedata, "=")), id).Exec()
	}

	if err == nil {
		err = self.GetDb().Commit()
		self.AddLog(fmt.Sprintf("postdata:%+v ", self.GetPost()))
		self.AjaxReturnSuccess("", nil)
		return
	} else {
		err = self.GetDb().Rollback()
		self.AjaxReturnError(errors.WithStack(err))
	}

}
