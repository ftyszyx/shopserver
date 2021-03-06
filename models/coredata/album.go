package coredata

import (
	"fmt"

	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type Album struct {
	models.Model
}

type AlbumData struct {
	Name     string `empty:"类型名不能为空"`
	Order_id int
}

func (self *Album) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *Album) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
	photoTableName := models.GetModel(names.PHOTO).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("photo") == true) {

		fieldstr += fmt.Sprintf("left join `%s` `photo` ON `album`.`cover_pic`=`photo`.`id`", photoTableName)
	}
	return sql.Alias("album").Join(fieldstr)
}
func (self *Album) InitField(sql db.SqlType) db.SqlType {
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
