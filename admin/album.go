package admin

type AlbumController struct {
	BaseController
}

func (self *AlbumController) BeforeSql(data map[string]interface{}) {

}

func (self *AlbumController) Add() {
	self.AddCommon(self)
}

func (self *AlbumController) Edit() {
	self.EditCommon(self)
}

func (self *AlbumController) Del() {
	self.DelCommon(self)
}

func (self *AlbumController) ChangeDefault() {
	self.CheckFieldExit(self.postdata, "default", "图片不能为空")
	self.CheckFieldExit(self.postdata, "id", "要修改的相册不能为空")
	changedata := make(map[string]interface{})
	changedata["default"] = self.postdata["default"]
	self.updateSqlById(self, changedata, self.postdata["id"])
}

func (self *AlbumController) ChangeCover() {

	self.CheckFieldExit(self.postdata, "cover_pic", "图片不能为空")
	self.CheckFieldExit(self.postdata, "id", "要修改的相册不能为空")
	changedata := make(map[string]interface{})
	changedata["cover_pic"] = self.postdata["cover_pic"]
	self.updateSqlById(self, changedata, self.postdata["id"])
}
