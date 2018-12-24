package models

import (
	"fmt"

	"github.com/zyx/shop_server/libs/db"
)

type Post struct {
	Model
}

type PostData struct {
	Title    string `empty:"标题不能为空"`
	Content  string `empty:"内容不能为空"`
	summary  string `empty:"摘要不能为空"`
	Type     string `empty:"文章类型不能为空"`
	Order_id int
	pic      string
}

//LEFT JOIN `aq_sys_user` `check_user` ON `sell`.`check_user`=`check_user`.`id`
func (self *Post) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *Post) GetModelStruct() interface{} {
	return PostData{}
}
func (self *Post) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
	posttypeTableName := GetModel(POSTTYPE).TableName()
	userTablename := GetModel(USER).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("user") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `user` ON `post`.`build_user`=`user`.`id`", userTablename)
	}
	if (allfield == true) || (sql.NeedJointable("posttype") == true) {

		fieldstr += fmt.Sprintf(" left join `%s` `posttype` ON `posttype`.`id`=`post`.`type`", posttypeTableName)
	}

	return sql.Alias("post").Join(fieldstr)
}
func (self *Post) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"user.name":       "build_user",
		"posttype.name":   "typename",
		"post.type":       "type",
		"post.id":         "id",
		"post.title":      "title",
		"post.build_time": "build_time",
		"post.content":    "content",
		"post.summary":    "summary",
		"post.pic":        "pic",
		"post.is_del":     "is_del",
		"post.order_id":   "order_id",
	})
}
