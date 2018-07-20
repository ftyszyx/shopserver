package models

import (
	"fmt"

	"github.com/zyx/shop_server/libs"
)

type Album struct {
	Model
}

type AlbumData struct {
	Name     string `empty:"类型名不能为空"`
	Order_id int
}

func (self *Album) InitSqlField(sql libs.SqlType) libs.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *Album) InitJoinString(sql libs.SqlType, allfield bool) libs.SqlType {
	photoTableName := GetModel(PHOTO).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("photo") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `photo` ON `album`.`cover_pic`=`photo`.`id`", photoTableName)
	}
	return sql.Alias("album").Join(fieldstr)
}
func (self *Album) InitField(sql libs.SqlType) libs.SqlType {
	return sql.Field(map[string]string{
		"photo.path":      "cover_pic_path",
		"album.id":        "id",
		"album.name":      "name",
		"album.default":   "default",
		"album.cover_pic": "cover_pic",
		"album.order_id":  "order_id",
		"album.is_del":    "is_del",
	})
}

func (self *Album) GetModelStruct() interface{} {
	return AlbumData{}
}
