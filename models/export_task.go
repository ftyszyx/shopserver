package models

import (
	"fmt"

	"github.com/zyx/shop_server/libs/db"
)

type ExportTask struct {
	Model
}

func (self *ExportTask) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}
func (self *ExportTask) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
	userTable := GetModel(USER).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("user") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `user` ON `user`.`id`=`exporttask`.`user_id`", userTable)
	}
	return sql.Alias("exporttask").Join(fieldstr)
}

func (self *ExportTask) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"exporttask.id":         "id",
		"exporttask.user_id":    "user_id",
		"user.name":             "user_name",
		"exporttask.build_time": "build_time",
		"exporttask.name":       "name",
		"exporttask.path":       "path",
	})
}
