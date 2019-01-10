package base

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/coredata"
	"github.com/zyx/shop_server/models/names"
)

type BaseIO interface {
	GetModel() models.ModelInterface
	GetUid() string
	GetPost() map[string]interface{}
	GetControl() string
	GetMethod() string
}

type BaseController struct {
	beego.Controller
	control  string
	method   string
	uid      string      //角色id
	token    string      //token
	dboper   db.DBOperIO //数据库操作连接
	model    models.ModelInterface
	postdata map[string]interface{}
}

func (self *BaseController) GetModel() models.ModelInterface {
	return self.model
}

func (self *BaseController) GetUid() string {
	return self.uid
}
func (self *BaseController) SetUid(id string) {
	self.uid = id
}

func (self *BaseController) GetPost() map[string]interface{} {
	return self.postdata
}

func (self *BaseController) GetDb() db.DBOperIO {
	return self.dboper
}

func (self *BaseController) GetControl() string {
	return self.control
}

func (self *BaseController) GetMethod() string {
	return self.method
}

func (self *BaseController) Prepare() {
	self.control, self.method = self.GetControllerAndAction()
	self.control = strings.ToLower(self.control)
	self.control = strings.TrimSuffix(self.control, "controller")
	self.model = models.GetModel(self.control)
	module := models.GetModel(names.MODULE).(*coredata.Module)

	logs.Info("control:%s action:%s method:%s", self.control, self.method, self.Ctx.Input.Method())
	//不用登录
	//或者post数据
	if self.Ctx.Input.Method() == "POST" {
		if self.Ctx.Input.RequestBody != nil && len(self.Ctx.Input.RequestBody) > 0 {
			//logs.Info("postdata body:%v", self.Ctx.Input.RequestBody)
			err := json.Unmarshal(self.Ctx.Input.RequestBody, &self.postdata)
			if err != nil {
				self.AjaxReturn(libs.ErrorCode, err.Error(), nil)
			}
			logs.Info("postdata %+v", self.postdata)
		}
	}

	self.token = self.Ctx.Request.Header.Get("token")
	if self.token == "" {
		self.token = self.Input().Get("token")
	}
	self.uid = self.Ctx.Request.Header.Get("uid")
	if self.uid == "" {
		self.uid = self.Input().Get("uid")
	}

	logs.Info("uid:%s token:%s ", self.uid, self.token)
	self.dboper = db.NewOper()
	if module.NeedAuth(self.control, self.method) == false {
		logs.Info("not need auth ")
		return
	}
	if self.uid == "" {
		self.AjaxReturn(libs.AuthFail, "请先登录", nil)
	}
	if self.token == "" {
		self.AjaxReturn(libs.AuthFail, "请先登录", nil)
	}

	usermodel := models.GetModel(names.USER).(*coredata.User)
	ok, msg, code := usermodel.Auth(self.dboper, self.token, self.uid, self.control, self.method)
	if ok == false {
		self.AjaxReturn(code, msg, nil)
	}
}
func (self *BaseController) Finish() {
	logs.Info("finish")
}

//json 返回
func (self *BaseController) AjaxReturn(code int, msg interface{}, data interface{}) {
	self.dboper.RollbackIfNotNull()
	libs.AjaxReturn(&self.Controller, code, msg, data)
}

func (self *BaseController) AjaxReturnError(err error) {
	self.dboper.RollbackIfNotNull()
	logs.Error("err:%+v", err)
	libs.AjaxReturn(&self.Controller, libs.ErrorCode, err.Error(), nil)
}

func (self *BaseController) AjaxReturnErrorMsg(msg interface{}) {
	self.dboper.RollbackIfNotNull()
	libs.AjaxReturn(&self.Controller, libs.ErrorCode, msg, nil)
}

func (self *BaseController) AjaxReturnSuccess(msg interface{}, data interface{}) {
	self.dboper.RollbackIfNotNull()
	libs.AjaxReturn(&self.Controller, libs.SuccessCode, msg, data)
}

func (self *BaseController) AjaxReturnSuccessNull() {
	self.dboper.RollbackIfNotNull()
	libs.AjaxReturn(&self.Controller, libs.SuccessCode, "", nil)
}

// CheckExit 检查字段是否存在  checkExitvalue:true 只检查数据里有的字段  false：检查所有 返回 nil表示存在  非nil表示不存在
func (self *BaseController) CheckExit(stru interface{}, data map[string]interface{}, checkExitvalue bool) error {
	model := models.GetModel(self.control)
	v := reflect.ValueOf(stru)
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		field := strings.ToLower(fi.Name)
		value, have := data[field]
		//logs.Info("field:%s", field)
		//检查空
		if tagv := fi.Tag.Get("empty"); tagv != "" {
			//logs.Info("have:%t", have)
			if checkExitvalue {
				//只检查字段存在的字段
				if have {
					if self.checkEmpty(value) {
						return errors.New(tagv)
					}
				}
			} else {
				//检查是否存在
				if have == false {
					return errors.New(tagv)
				} else {
					if self.checkEmpty(value) {
						return errors.New(tagv)
					}
				}
			}
		}
		//检查数据是否唯一
		if tagv := fi.Tag.Get("only"); tagv != "" {
			if have {
				if model.CheckExit(self.dboper, field, value) {
					return errors.New(tagv)
				}
			}
		}
	}
	return nil

}

//CheckFieldExitAndReturn 检查字段是否存在
func (self *BaseController) CheckFieldExitAndReturn(data map[string]interface{}, field string, errtext string) {
	if self.CheckFieldExit(data, field) == false {
		self.AjaxReturn(libs.ErrorCode, errtext, nil)
	}
}

func (self *BaseController) CheckFieldExit(data map[string]interface{}, field string) bool {
	value, ok := data[field]
	if ok {
		if self.checkEmpty(value) {
			return false
		}
		return true
	}
	return false
}

func (self *BaseController) checkEmpty(value interface{}) bool {
	if valuestr, okstr := value.(string); okstr {
		if strings.TrimSpace(valuestr) == "" {
			return true
		}
	} else if valueint, okint := value.(int); okint {
		if valueint == 0 {
			return true
		}
	}
	return false
}

func (self *BaseController) saveFile() (string, string, int64, string, error) {

	tempfolder := beego.AppConfig.String("site.tempfolder")
	file, header, err := self.GetFile("upfile")
	if err != nil {
		return "", "", 0, "", err
	}
	md5h := md5.New()
	io.Copy(md5h, file)

	filemd5 := md5h.Sum(nil)

	md5str1 := fmt.Sprintf("%x", filemd5)

	namearr := strings.Split(header.Filename, ".")
	filetype := namearr[len(namearr)-1]
	err = self.SaveToFile("upfile", tempfolder+md5str1+"."+filetype)

	return md5str1, header.Filename, header.Size, filetype, err

}

func (self *BaseController) upload() (error, map[string]interface{}) {
	return self.uploadCommon(true)

}

func (self *BaseController) uploadCommon(uptoqiniu bool) (error, map[string]interface{}) {

	err, fileinfo := self.saveUploadFile()
	if err != nil {
		return errors.WithStack(err), nil
	}
	if uptoqiniu == true {
		bucket := beego.AppConfig.String("qiniu.bucket")
		_, err = libs.UploadFile(fileinfo["filename"].(string), fileinfo["filePath"].(string), bucket)
		if err != nil {
			return errors.WithStack(err), nil
		}
	}

	return nil, fileinfo

}

func (self *BaseController) saveUploadFile() (error, map[string]interface{}) {
	var fileinfo = make(map[string]interface{})
	fileName, filetitle, filesize, filetype, err := self.saveFile()
	tempfolder := beego.AppConfig.String("site.tempfolder")
	if err == nil {

		host := beego.AppConfig.String("qiniu.host")
		fileName = fileName + "." + filetype
		url := host + fileName
		filePath := tempfolder + fileName
		logs.Info("path:%s filetype:%s", filePath, filetype)
		fileinfo["filePath"] = filePath
		fileinfo["filename"] = fileName
		fileinfo["filetitle"] = filetitle
		fileinfo["filesize"] = filesize
		fileinfo["filetype"] = filetype
		fileinfo["url"] = url
		fileinfo["host"] = host
		return nil, fileinfo
	}
	return errors.WithStack(err), nil
}
