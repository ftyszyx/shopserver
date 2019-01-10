package logistics

import (
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type LogisticsTask struct {
	models.Model
}

type LogisticsTaskData struct {
	Name     string `empty:"名称不能为空"`
	Tasklist string `empty:"进度项不能为空"`
}

func (self *LogisticsTask) InitSqlField(sql db.SqlType) db.SqlType {
	//return self.InitField(self.InitJoinString(sql, true))
	return sql
}

func (self *LogisticsTask) GetModelStruct() interface{} {
	return LogisticsTaskData{}
}

// func (self *LogisticsTask) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {

// }
