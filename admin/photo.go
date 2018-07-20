package admin

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/qiniu/api.v7/storage"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
)

type PhotoController struct {
	BaseController
}

//检查数据正确性
func (self *PhotoController) checkData(data map[string]interface{}) {
	albummodel := models.GetModel(models.ALBUM)
	if value, ok := data["album"]; ok {
		albumid := value.(string)
		if albummodel.CheckExit("id", albumid) == false {
			self.AjaxReturnError("相册不存在")
		}
	}
}

func (self *PhotoController) BeforeSql(data map[string]interface{}) {
	if self.method == "Add" {
		picpath := self.postdata["path"]
		if self.model.CheckExit("path", picpath) {
			self.AjaxReturnSuccess("已经存在", nil)
		}
		self.checkData(self.postdata)
	} else if self.method == "Edit" {
		self.checkData(self.postdata)
	}
}
func (self *PhotoController) AfterSql(data map[string]interface{}, oldinfo orm.Params) {

	if self.GetMethod() == "Del" {
		manger := libs.GetManger()
		self.CheckFieldExit(self.postdata, "key", "key为空")
		key := self.postdata["key"].(string)
		bucket := beego.AppConfig.String("qiniu.bucket")
		err := manger.Delete(bucket, key)
		if err != nil {
			self.AjaxReturnError(err.Error())
			return
		}
		self.AddLog(fmt.Sprintf("%+v", data))
	} else {
		self.AddLog(fmt.Sprintf("%+v", data))
	}
}

func (self *PhotoController) Add() {
	self.AddCommon(self)
}

func (self *PhotoController) Edit() {
	self.EditCommon(self)
}

//删除
func (self *PhotoController) Del() {

	ids, ok := self.postdata["id"].([]interface{})
	if ok == false || len(ids) == 0 {
		self.AjaxReturnError("id空")
		return
	}
	if len(ids) > 1000 {
		self.AjaxReturnError("批量删除数量不能超过1000")
		return
	}
	keys, ok := self.postdata["keys"].([]interface{})
	if ok == false || len(keys) == 0 {
		self.AjaxReturnError("key为空")
		return
	}
	o := orm.NewOrm()
	idstr := libs.SqlGetArrInfo(ids)
	_, err := o.Raw(fmt.Sprintf("delete from %s where `id` in %s", self.model.TableName(), idstr)).Exec()
	if err == nil {
		deleteOps := make([]string, 0, len(keys))
		bucket := beego.AppConfig.String("qiniu.bucket")
		for _, key := range keys {
			deleteOps = append(deleteOps, storage.URIDelete(bucket, key.(string)))
		}
		manger := libs.GetManger()
		_, err := manger.Batch(deleteOps)
		if err != nil {
			self.AjaxReturnError(err.Error())
		}
		self.AjaxReturnSuccess("", nil)
	}
	self.AjaxReturnError(err.Error())
}

//移动
func (self *PhotoController) MoveMulti() {
	ids, ok := self.postdata["id"].([]interface{})
	if ok == false || len(ids) == 0 {
		self.AjaxReturnError("id空")
		return
	}
	album, ok := self.postdata["album"].(string)
	if ok == false {
		self.AjaxReturnError("目标相册要填")
	}
	albuminfo := self.model.GetInfoById(album)
	if albuminfo == nil {
		self.AjaxReturnError("目标相册不存在")
	}
	idstr := libs.SqlGetArrInfo(ids)
	o := orm.NewOrm()
	_, err := o.Raw(fmt.Sprintf("update %s set `album`=?  where `id` in %s", self.model.TableName(), idstr), album).Exec()
	if err == nil {
		self.AjaxReturnSuccess("", nil)
	}
	self.AjaxReturnError(err.Error())
}

func (self *PhotoController) GetUploadtoken() {
	key, _ := self.postdata["key"]
	keystr := key.(string)
	url := beego.AppConfig.String("qiniu.host")
	bucket := beego.AppConfig.String("qiniu.bucket")
	upToken := libs.GetUploadToken(keystr, bucket)
	var senddata = make(map[string]interface{})
	senddata["uptoken"] = upToken
	senddata["url"] = url
	self.AjaxReturnSuccess("", senddata)
}
