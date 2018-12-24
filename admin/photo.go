package admin

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/qiniu/api.v7/storage"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type PhotoController struct {
	BaseController
}

//检查数据正确性
func (self *PhotoController) checkData(dboper db.DBOperIO, data map[string]interface{}) error {
	albummodel := models.GetModel(models.ALBUM)
	if value, ok := data["album"]; ok {
		albumid := value.(string)
		if albummodel.CheckExit(dboper, "id", albumid) == false {
			return errors.New("相册不存在")
		}
	}
	return nil
}

func (self *PhotoController) BeforeSql(data map[string]interface{}) error {
	if self.method == "Add" {
		picpath := self.postdata["path"]
		if self.model.CheckExit(self.dboper, "path", picpath) {
			return errors.New("已经存在该图片")
		}
		return self.checkData(self.dboper, self.postdata)
	} else if self.method == "Edit" {
		return self.checkData(self.dboper, self.postdata)
	}
	return nil
}
func (self *PhotoController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {

	if self.GetMethod() == "Del" {
		manger := libs.GetManger()
		self.CheckFieldExitAndReturn(self.postdata, "key", "key为空")
		key := self.postdata["key"].(string)
		bucket := beego.AppConfig.String("qiniu.bucket")
		err := manger.Delete(bucket, key)
		if err != nil {
			return err
		}
		self.AddLog(fmt.Sprintf("%+v", data))
	} else {
		self.logcommon(data, oldinfo)
	}
	return nil
}

func (self *PhotoController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *PhotoController) Edit() {
	self.EditCommonAndReturn(self)
}

//删除
func (self *PhotoController) Del() {

	ids, ok := self.postdata["id"].([]interface{})
	if ok == false || len(ids) == 0 {
		self.AjaxReturnError(errors.New("id空"))
		return
	}
	if len(ids) > 1000 {
		self.AjaxReturnError(errors.New("批量删除数量不能超过1000"))
		return
	}
	keys, ok := self.postdata["keys"].([]interface{})
	if ok == false || len(keys) == 0 {
		self.AjaxReturnError(errors.New("key为空"))
		return
	}
	// o := orm.NewOrm()
	idstr := db.SqlGetArrInfo(ids)
	_, err := self.dboper.Raw(fmt.Sprintf("delete from %s where `id` in %s", self.model.TableName(), idstr)).Exec()
	if err == nil {
		deleteOps := make([]string, 0, len(keys))
		bucket := beego.AppConfig.String("qiniu.bucket")
		for _, key := range keys {
			deleteOps = append(deleteOps, storage.URIDelete(bucket, key.(string)))
		}
		manger := libs.GetManger()
		_, err := manger.Batch(deleteOps)
		if err != nil {
			self.AjaxReturnError(errors.WithStack(err))
		}
		self.AjaxReturnSuccess("", nil)
	}
	self.AjaxReturnError(errors.WithStack(err))
}

//移动
func (self *PhotoController) MoveMulti() {
	ids, ok := self.postdata["id"].([]interface{})
	if ok == false || len(ids) == 0 {
		self.AjaxReturnError(errors.New("id空"))
		return
	}
	album, ok := self.postdata["album"].(string)
	if ok == false {
		self.AjaxReturnError(errors.New("目标相册要填"))
	}
	albuminfo := self.model.GetInfoById(self.dboper, album)
	if albuminfo == nil {
		self.AjaxReturnError(errors.New("目标相册不存在"))
	}
	idstr := db.SqlGetArrInfo(ids)
	// o := orm.NewOrm()
	_, err := self.dboper.Raw(fmt.Sprintf("update %s set `album`=?  where `id` in %s", self.model.TableName(), idstr), album).Exec()
	if err == nil {
		self.AddLog(fmt.Sprintf("%+v", ids))
		self.AjaxReturnSuccess("", nil)
	}
	self.AjaxReturnError(errors.WithStack(err))
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
