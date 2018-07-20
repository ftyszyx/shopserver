package models

import (
	"fmt"

	"github.com/zyx/shop_server/libs"
)

type Log struct {
	Model
}

//LEFT JOIN `aq_sys_user` `check_user` ON `sell`.`check_user`=`check_user`.`id`
func (self *Log) InitSqlField(sql libs.SqlType) libs.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *Log) InitJoinString(sql libs.SqlType, allfield bool) libs.SqlType {
	moduleTableName := GetModel(MODULE).TableName()
	userTablename := GetModel(USER).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("user") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `user` ON `log`.`userid`=`user`.`id`", userTablename)
	}
	if (allfield == true) || (sql.NeedJointable("module") == true) {

		fieldstr += fmt.Sprintf(" left join `%s` `module` ON `log`.`method`=`module`.`method` and `log`.`controller`=`module`.`controller`", moduleTableName)
	}
	return sql.Alias("log").Join(fieldstr)
}

func (self *Log) InitField(sql libs.SqlType) libs.SqlType {
	return sql.Field(map[string]string{
		"user.name":   "user_name",
		"module.name": "module_name",
		"log.time":    "time",
		"log.link":    "link",
		"log.info":    "info",
	})
}
