package admin

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
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
	uid      string //角色id
	token    string //token
	model    models.ModelInterface
	postdata map[string]interface{}
}

func (self *BaseController) GetModel() models.ModelInterface {
	return self.model
}

func (self *BaseController) GetUid() string {
	return self.uid
}

func (self *BaseController) GetPost() map[string]interface{} {
	return self.postdata
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
	module := models.GetModel(models.MODULE).(*models.Module)

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

	if module.NeedAuth(self.control, self.method) == false {
		logs.Info("not auth ")
		return
	}
	if self.uid == "" {
		self.AjaxReturn(libs.AuthFail, "token空", nil)
	}
	if self.token == "" {
		self.AjaxReturn(libs.AuthFail, "token空", nil)
	}

	usermodel := models.GetModel(models.USER).(*models.User)
	ok, msg, code := usermodel.Auth(self.token, self.uid, self.control, self.method)
	if ok == false {
		self.AjaxReturn(code, msg, nil)
	}
}

//json 返回
func (self *BaseController) AjaxReturn(code int, msg interface{}, data interface{}) {
	libs.AjaxReturn(&self.Controller, code, msg, data)
}

func (self *BaseController) AjaxReturnError(msg interface{}) {
	libs.AjaxReturn(&self.Controller, libs.ErrorCode, msg, nil)
}

func (self *BaseController) AjaxReturnSuccess(msg interface{}, data interface{}) {
	libs.AjaxReturn(&self.Controller, libs.SuccessCode, msg, data)
}

//获取所有时请求数据
type AllReqData struct {
	Page   int
	Rownum int
	Order  map[string]interface{}
	And    bool
	Search map[string]interface{}
}

//通用的查询列表
func (self *BaseController) All() {
	var data = AllReqData{And: true}
	err := json.Unmarshal(self.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.Info(err.Error())
		self.AjaxReturn(libs.ErrorCode, err.Error(), nil)
		return
	}
	// logs.Info("data :%v body:%s  ", data, string(self.Ctx.Input.RequestBody))
	// logs.Info("alldata :%+v", data)
	self.AllExc(data)
}
func (self *BaseController) AllExc(data AllReqData) {
	err, num, datalist := self.AllExcCommon(data, libs.GetAll_type)
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	var senddata = make(map[string]interface{})
	senddata["num"] = num
	senddata["list"] = datalist
	self.AjaxReturn(libs.SuccessCode, "", senddata)
}

func (self *BaseController) AllExcCommon(data AllReqData, gettype int) (error, int, []orm.Params) {
	o := orm.NewOrm()
	var totalnum = 0
	var dataList []orm.Params
	var sqltext libs.SqlType
	model := models.GetModel(self.control)
	sqltext = &libs.SqlBuild{}
	sqltext = sqltext.Name(model.TableName())

	if data.And {
		sqltext = sqltext.Where(data.Search)
	} else {
		sqltext = sqltext.WhereOr(data.Search)
	}
	sqltext = model.InitJoinString(sqltext, false)
	num, err := o.Raw(sqltext.Count()).Values(&dataList)
	if err == nil && num > 0 {
		totalnum, err = strconv.Atoi(dataList[0][libs.SQL_COUNT_NAME].(string))
		if err != nil {
			return err, 0, nil
		}
		if gettype == libs.GetAll_type_num {
			return nil, totalnum, nil
		}
		sqltext = sqltext.Order(data.Order)

		if data.Page == 0 {
			//不用分页
			sqltext = model.InitJoinString(model.InitField(sqltext), true)

			num, err = o.Raw(sqltext.Select()).Values(&dataList)
			if err == nil {
				return nil, totalnum, dataList
			} else {
				return err, 0, nil
			}
		} else {
			//用分页
			var start = (data.Page - 1) * data.Rownum
			if totalnum > 1000 {
				//总数很多
				tablealias := sqltext.GetAlias()
				selfidname := "id"
				if tablealias != "" {
					selfidname = tablealias + ".id"
				}
				subsql := sqltext.Limit([]int{start, data.Rownum}).Field(map[string]string{selfidname: "id"}).Select()
				var newsqltext libs.SqlType
				newsqltext = &libs.SqlBuild{}
				newsqltext = newsqltext.Name(model.TableName()).Order(data.Order)
				newsqltext = model.InitJoinString(model.InitField(newsqltext), true)
				oldjoinstr := newsqltext.GetJoinStr()
				newsqltext.Join(oldjoinstr + fmt.Sprintf(" INNER join (%s) a ON `a`.`id`=%s ", subsql, libs.SqlGetKey(selfidname)))
				num, err = o.Raw(newsqltext.Select()).Values(&dataList)
				if err == nil {
					return nil, totalnum, dataList
				} else {
					self.AjaxReturnError(err.Error())
					return err, 0, nil
				}

			} else {
				//按正常方式
				sqltext = model.InitJoinString(model.InitField(sqltext), true)
				num, err = o.Raw(sqltext.Limit([]int{start, data.Rownum}).Select()).Values(&dataList)
				if err == nil {
					return nil, totalnum, dataList
				} else {
					self.AjaxReturnError(err.Error())
					return err, 0, nil
				}
			}

		}
	} else {
		return err, 0, nil
	}
}

//检查字段是否存在  checkExitvalue:true 只检查数据里有的字段  false：检查所有
func (self *BaseController) CheckExit(stru interface{}, data map[string]interface{}, checkExitvalue bool) bool {
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
					return self.isEmpty(value, tagv)
				}
			} else {
				//检查是否存在

				if have == false {
					self.AjaxReturn(libs.ErrorCode, tagv, nil)
					return false
				} else {
					if self.isEmpty(value, tagv) {
						return false
					}
				}
			}
		}
		//检查数据是否唯一
		if tagv := fi.Tag.Get("only"); tagv != "" {
			if model.CheckExit(field, value) {
				self.AjaxReturn(libs.ErrorCode, tagv, nil)
				return false
			}
		}
	}
	return true

}

//检查字段是否存在
func (self *BaseController) CheckFieldExit(data map[string]interface{}, field string, errtext string) bool {
	value, ok := data[field]
	if ok {
		if self.isEmpty(value, errtext) {
			return false
		}
		return true
	}
	self.AjaxReturn(libs.ErrorCode, errtext, nil)
	return false
}

func (self *BaseController) isEmpty(value interface{}, errtext string) bool {
	if valuestr, okstr := value.(string); okstr {
		if strings.TrimSpace(valuestr) == "" {
			self.AjaxReturn(libs.ErrorCode, errtext, nil)
			return true
		}
	} else if valueint, okint := value.(int); okint {
		if valueint == 0 {
			self.AjaxReturn(libs.ErrorCode, errtext, nil)
			return true
		}
	}
	return false
}

func (self *BaseController) BeforeSql(data map[string]interface{}) {

}

func (self *BaseController) AfterSql(data map[string]interface{}, oldinfo orm.Params) {
	self.AddLog(fmt.Sprintf("%+v", data))
}

func (self *BaseController) AddCommon(sqlcall SqlIO) {
	datacheck := self.model.GetModelStruct()
	self.CheckExit(datacheck, self.postdata, false)
	adddata := libs.ClearMapByStruct(self.postdata, datacheck)
	self.AddCommonExe(sqlcall, adddata)
}

func (self *BaseController) AddCommonExe(sqlcall SqlIO, adddata map[string]interface{}) {
	self.AddCommonTable(sqlcall, adddata, self.model.TableName())
}

func (self *BaseController) AddCommonTable(sqlcall SqlIO, adddata map[string]interface{}, table string) {
	o := orm.NewOrm()
	sqlcall.BeforeSql(adddata)
	keys, values := libs.SqlGetInsertInfo(adddata)
	_, err := o.Raw(fmt.Sprintf("insert into %s (%s) values (%s)", table, keys, values)).Exec()
	if err == nil {
		sqlcall.AfterSql(adddata, nil)
		self.AjaxReturn(libs.SuccessCode, "", nil)
		return
	}
	self.AjaxReturn(libs.ErrorCode, err.Error(), nil)
}

func (self *BaseController) EditCommon(sqlcall SqlIO) {
	modelcheck := self.model.GetModelStruct()
	self.CheckFieldExit(self.postdata, "id", "id为空")
	id := self.postdata["id"].(string)
	self.CheckExit(modelcheck, self.postdata, true)
	changedata := libs.ClearMapByStruct(self.postdata, modelcheck)
	if len(changedata) == 0 {
		self.AjaxReturnError("没有修改")
	}
	self.updateSqlCommon(sqlcall, changedata, "id", id)
}

func (self *BaseController) updateSqlById(sqlcall SqlIO, changedata map[string]interface{}, id interface{}) {
	self.updateSqlCommon(sqlcall, changedata, "id", id)
}

func (self *BaseController) updateSqlCommon(sqlcall SqlIO, changedata map[string]interface{}, field string, id interface{}) {
	o := orm.NewOrm()
	sqlcall.BeforeSql(changedata)
	oldinfo := self.model.GetInfoByField(field, id)
	if oldinfo == nil {
		self.AjaxReturnError("没找到")
	}
	_, err := o.Raw(fmt.Sprintf("update %s set %s where `%s`=?", self.model.TableName(), libs.SqlGetKeyValue(changedata, "="), field), id).Exec()
	if err == nil {
		sqlcall.AfterSql(changedata, oldinfo[0])
		self.AjaxReturnSuccess("", nil)
		return
	}
	self.AjaxReturnError(err.Error())
}

func (self *BaseController) DelCommon(sqlcall SqlIO) {
	self.CheckFieldExit(self.postdata, "id", "id为空")
	id := self.postdata["id"].(string)
	oldinfo := self.model.GetInfoById(id)
	if oldinfo == nil {
		self.AjaxReturnError("id 不存在")
	}
	sqlcall.BeforeSql(oldinfo)
	o := orm.NewOrm()
	_, err := o.Raw(fmt.Sprintf("delete from %s where `id`=?", self.model.TableName()), id).Exec()
	if err == nil {
		sqlcall.AfterSql(oldinfo, nil)
		self.AjaxReturn(libs.SuccessCode, "", nil)
		return
	}
	self.AjaxReturn(libs.ErrorCode, err.Error(), nil)
}

//增加日志
func (self *BaseController) AddLog(info string) {
	models.AddLog(self.uid, info, self.control, self.method)
}

func (self *BaseController) AddLogLink(info string, link string) {
	models.AddLogLink(self.uid, info, self.control, self.method, link)
}

func (self *BaseController) saveFile() (string, string, int64, string, error) {

	tempfolder := beego.AppConfig.String("site.tempfolder")
	file, header, err := self.GetFile("upfile")
	if err != nil {

		return "", "", 0, "", err
	}
	//logs.Info(" name:%s  content:%s", header.Filename, header.Header.Get("lastModifiedDate"))
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
	var fileinfo = make(map[string]interface{})
	fileName, filetitle, filesize, filetype, err := self.saveFile()
	tempfolder := beego.AppConfig.String("site.tempfolder")
	if err == nil {
		bucket := beego.AppConfig.String("qiniu.bucket")
		host := beego.AppConfig.String("qiniu.host")
		url := host + fileName
		filePath := tempfolder + fileName + "." + filetype
		_, err = libs.UploadFile(fileName, filePath, bucket)
		fileinfo["filePath"] = filePath
		fileinfo["filename"] = fileName
		fileinfo["filetitle"] = filetitle
		fileinfo["filesize"] = filesize
		fileinfo["filetype"] = filetype
		fileinfo["url"] = url
		fileinfo["host"] = host
		return err, fileinfo
	}
	return err, nil
}

//导出csv表
func (self *BaseController) ExportCsvCommon() {
	search := self.postdata["search"].(map[string]interface{})
	headlist := libs.GetStrArr(self.postdata["headlist"].([]interface{}))
	filename := self.postdata["filename"].(string)
	namelist := libs.GetStrArr(self.postdata["namelist"].([]interface{}))
	tempfolder := beego.AppConfig.String("site.tempfolder")
	if headlist == nil || len(headlist) == 0 || namelist == nil || len(namelist) == 0 {
		self.AjaxReturnError("需要导出的字段为空")
	}
	var limitpagenum = 1000000
	var reqdata = AllReqData{Search: search, Rownum: limitpagenum}

	_, num, _ := self.AllExcCommon(reqdata, libs.GetAll_type_num)

	if num == 0 {
		self.AjaxReturnError("没有数据可导出")
	}
	filepath := tempfolder + filename + ".csv"

	fileio, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	defer fileio.Close()
	logs.Info("begin write header")
	totalpagenum := num/limitpagenum + 1
	if line, err := libs.UTF82GBK(strings.Join(headlist, ",")); err == nil { // 写入一行
		fileio.WriteString(line + "\n")
	}
	logs.Info("begin getdata")
	for curpage := 1; curpage <= totalpagenum; curpage++ {
		reqdata.Page = curpage

		err, _, list := self.AllExcCommon(reqdata, libs.GetAll_type)
		if err != nil {
			self.AjaxReturnError(err.Error())
		}
		for _, rowdata := range list {
			//每一行
			var rowstrarr []string
			for _, name := range namelist {
				rowstrarr = append(rowstrarr, rowdata[name].(string))
			}
			if line, err := libs.UTF82GBK(strings.Join(rowstrarr, ",")); err == nil { // 写入一行
				fileio.WriteString(line + "\n")
			}
		}

	}
	fileio.Close()

	logs.Info("begin upload")

	filemd5str := libs.GetFileMd5(filepath)

	bucket := beego.AppConfig.String("qiniu.bucket")
	host := beego.AppConfig.String("qiniu.host")
	url := host + filemd5str + ".csv"
	logs.Info("update filepath:%s", filemd5str)
	_, err = libs.UploadFile(filemd5str+".csv", filepath, bucket)
	os.Remove(filepath)
	if err != nil {
		self.AjaxReturnError("upload file err:" + err.Error())
	}

	//增加一项
	adddata := make(map[string]interface{})
	adddata["user_id"] = self.uid
	adddata["build_time"] = time.Now().Unix()
	adddata["path"] = url
	adddata["name"] = filename
	logs.Info("begin save task")
	exporttaskTable := models.GetModel(models.EXPORT_TASK).TableName()
	self.AddCommonTable(self, adddata, exporttaskTable)
}
