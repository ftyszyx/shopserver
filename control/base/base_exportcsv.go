package base

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

func (self *BaseController) ExportCsvCommonAndReturn() {
	err, _ := self.ExportCsvCommon()
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	self.AjaxReturnSuccessNull()
}

func (self *BaseController) ExportCsvCommonSearch(search map[string]interface{}) (error, map[string]interface{}) {
	headlist := libs.GetStrArr(self.GetPost()["headlist"].([]interface{}))
	filename := self.GetPost()["filename"].(string)
	namelist := libs.GetStrArr(self.GetPost()["namelist"].([]interface{}))
	tempfolder := beego.AppConfig.String("site.tempfolder")
	if headlist == nil || len(headlist) == 0 || namelist == nil || len(namelist) == 0 {

		return errors.New("需要导出的字段为空"), nil
	}
	var limitpagenum = 1000000
	var reqdata = models.AllReqData{Search: search, Rownum: limitpagenum, And: true}

	_, num, _ := self.AllExcCommon(reqdata, libs.GetAll_type_num)

	if num == 0 {
		return errors.New("没有数据可导出"), nil
	}
	filepath := tempfolder + filename + ".csv"

	fileio, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.WithStack(err), nil
	}
	defer fileio.Close()

	totalpagenum := num/limitpagenum + 1
	logs.Info("begin write header: num:%d limitpagenum:%d totalpage:%d", num, limitpagenum, totalpagenum)
	if line, err := libs.UTF82GBK(strings.Join(headlist, ",")); err == nil { // 写入一行
		fileio.WriteString(line + "\n")
	}
	logs.Info("begin getdata")
	for curpage := 1; curpage <= totalpagenum; curpage++ {
		reqdata.Page = curpage
		//logs.Info("get page:%d", curpage)
		err, _, list := self.AllExcCommon(reqdata, libs.GetAll_type)
		if err != nil {
			return err, nil
		}

		for _, rowdata := range list {
			//每一行
			var rowstrarr []string
			for _, name := range namelist {
				datastr, err := self.model.ExportNameProcess(name, rowdata[name], rowdata)
				if err != nil {
					return errors.WithStack(err), nil
				}
				datastr = strings.Replace(datastr, `"`, `""`, -1)
				var datastrnew = ""
				datastrnew, err = libs.UTF82GBK(datastr)
				if err != nil { // 写入一行
					logs.Error("write line err:", err.Error())
					logs.Error(fmt.Sprintf("问题行id:%+v 问题字段名:%s 问题字段内容:%s 错误信息:%s ", rowdata["id"], name, datastr, err.Error()))

					return errors.New(fmt.Sprintf("问题行id:%+v 问题字段名:%s 问题字段内容:%s 错误信息:%s ", rowdata["id"], name, datastr, err.Error())), nil
				} else {
					rowstrarr = append(rowstrarr, `"`+datastrnew+`"`)
				}

			}
			rowstr := strings.Join(rowstrarr, ",")
			fileio.WriteString(rowstr + "\n")
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

		return errors.WithStack(err), nil
	}

	//增加一项
	adddata := make(map[string]interface{})
	adddata["user_id"] = self.uid
	adddata["build_time"] = time.Now().Unix()
	adddata["path"] = url
	adddata["name"] = filename
	logs.Info("begin save task", url)
	exporttaskTable := models.GetModel(names.EXPORT_TASK).TableName()
	return self.AddCommonTable(self, adddata, exporttaskTable), adddata
}

//导出csv表
func (self *BaseController) ExportCsvCommon() (error, map[string]interface{}) {
	search := self.GetPost()["search"].(map[string]interface{})
	return self.ExportCsvCommonSearch(search)

}
