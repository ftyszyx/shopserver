package coredata

import (
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type AdsPos struct {
	models.Model
}

type AdsPosData struct {
	Name      string `empty:"名称不能为空"`
	title     string
	title_pic string
}

func (self *AdsPos) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *AdsPos) GetModelStruct() interface{} {
	return AdsPosData{}
}
