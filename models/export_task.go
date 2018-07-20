package models

import (
	"fmt"

	"github.com/zyx/shop_server/libs"
)

type ExportTask struct {
	Model
}

func (self *ExportTask) InitSqlField(sql libs.SqlType) libs.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}
func (self *ExportTask) InitJoinString(sql libs.SqlType, allfield bool) libs.SqlType {
	userTable := GetModel(USER).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("user") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `user` ON `user`.`id`=`exporttask`.`user_id`", userTable)
	}
	return sql.Alias("exporttask").Join(fieldstr)
}

func (self *ExportTask) InitField(sql libs.SqlType) libs.SqlType {
	return sql.Field(map[string]string{
		"exporttask.id":         "id",
		"exporttask.user_id":    "user_id",
		"user.name":             "user_name",
		"exporttask.build_time": "build_time",
		"exporttask.name":       "name",
		"exporttask.path":       "path",
	})
}
