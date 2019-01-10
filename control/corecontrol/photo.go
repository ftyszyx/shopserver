package corecontrol

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/qiniu/api.v7/storage"
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type PhotoController struct {
	base.BaseController
}

//检查数据正确性
func (self *PhotoController) checkData(dboper db.DBOperIO, data map[string]interface{}) error {
	albummodel := models.GetModel(names.ALBUM)
	if value, ok := data["album"]; ok {
		albumid := value.(string)
		if albummodel.CheckExit(dboper, "id", albumid) == false {
			return errors.New("相册不存在")
		}
	}
	return nil
}

func (self *PhotoController) BeforeSql(data map[string]interface{}) error {
	if self.GetMethod() == "Add" {
		picpath := self.GetPost()["path"]
		if self.GetModel().CheckExit(self.GetDb(), "path", picpath) {
			return errors.New("已经存在该图片")
		}
		return self.checkData(self.GetDb(), self.GetPost())
	} else if self.GetMethod() == "Edit" {
		return self.checkData(self.GetDb(), self.GetPost())
	}
	return nil
}
func (self *PhotoController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {

	if self.GetMethod() == "Del" {
		manger := libs.GetManger()
		self.CheckFieldExitAndReturn(self.GetPost(), "key", "key为空")
		key := self.GetPost()["key"].(string)
		bucket := beego.AppConfig.String("qiniu.bucket")
		err := manger.Delete(bucket, key)
		if err != nil {
			return err
		}
		self.AddLog(fmt.Sprintf("%+v", data))
	} else {
		self.Logcommon(data, oldinfo)
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

	ids, ok := self.GetPost()["id"].([]interface{})
	if ok == false || len(ids) == 0 {
		self.AjaxReturnError(errors.New("id空"))
		return
	}
	if len(ids) > 1000 {
		self.AjaxReturnError(errors.New("批量删除数量不能超过1000"))
		return
	}
	keys, ok := self.GetPost()["keys"].([]interface{})
	if ok == false || len(keys) == 0 {
		self.AjaxReturnError(errors.New("key为空"))
		return
	}
	// o := orm.NewOrm()
	idstr := db.SqlGetArrInfo(ids)
	_, err := self.GetDb().Raw(fmt.Sprintf("delete from %s where `id` in %s", self.GetModel().TableName(), idstr)).Exec()
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
	ids, ok := self.GetPost()["id"].([]interface{})
	if ok == false || len(ids) == 0 {
		self.AjaxReturnError(errors.New("id空"))
		return
	}
	album, ok := self.GetPost()["album"].(string)
	if ok == false {
		self.AjaxReturnError(errors.New("目标相册要填"))
	}
	albuminfo := self.GetModel().GetInfoById(self.GetDb(), album)
	if albuminfo == nil {
		self.AjaxReturnError(errors.New("目标相册不存在"))
	}
	idstr := db.SqlGetArrInfo(ids)
	// o := orm.NewOrm()
	_, err := self.GetDb().Raw(fmt.Sprintf("update %s set `album`=?  where `id` in %s", self.GetModel().TableName(), idstr), album).Exec()
	if err == nil {
		self.AddLog(fmt.Sprintf("%+v", ids))
		self.AjaxReturnSuccess("", nil)
	}
	self.AjaxReturnError(errors.WithStack(err))
}

func (self *PhotoController) GetUploadtoken() {
	key, _ := self.GetPost()["key"]
	keystr := key.(string)
	url := beego.AppConfig.String("qiniu.host")
	bucket := beego.AppConfig.String("qiniu.bucket")
	upToken := libs.GetUploadToken(keystr, bucket)
	var senddata = make(map[string]interface{})
	senddata["uptoken"] = upToken
	senddata["url"] = url
	self.AjaxReturnSuccess("", senddata)
}
