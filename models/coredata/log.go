package coredata

import (
	"fmt"

	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type Log struct {
	models.Model
}

//LEFT JOIN `aq_sys_user` `check_user` ON `sell`.`check_user`=`check_user`.`id`
func (self *Log) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *Log) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
	moduleTableName := models.GetModel(names.MODULE).TableName()
	userTablename := models.GetModel(names.USER).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("user") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `user` ON `log`.`userid`=`user`.`id`", userTablename)
	}
	if (allfield == true) || (sql.NeedJointable("module") == true) {

		fieldstr += fmt.Sprintf(" left join `%s` `module` ON `log`.`method`=`module`.`method` and `log`.`controller`=`module`.`controller`", moduleTableName)
	}
	return sql.Alias("log").Join(fieldstr)
}

func (self *Log) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"user.name":   "user_name",
		"module.name": "module_name",
		"module.id":   "module_id",
		"log.time":    "time",
		"log.link":    "link",
		"log.info":    "info",
	})
}
