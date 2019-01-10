package coredata

import (
	"fmt"

	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type DataBase struct {
	models.Model
}

type DataBaseData struct {
	Name string `empty:"名称不能为空"`
}

func (self *DataBase) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *DataBase) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
	userTable := models.GetModel(names.USER).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("user") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `user` ON `user`.`id`=`database`.`user_id`", userTable)
	}
	return sql.Alias("database").Join(fieldstr)
}
func (self *DataBase) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"database.id":         "id",
		"database.user_id":    "user_id",
		"user.name":           "user_name",
		"database.build_time": "build_time",
		"database.name":       "name",
		"database.path":       "path",
	})
}

func (self *DataBase) GetModelStruct() interface{} {
	return DataBaseData{}
}
