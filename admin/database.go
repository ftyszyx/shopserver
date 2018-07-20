package admin

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

//数据库备份还原
type DatabaseController struct {
	BaseController
}

func (self *DatabaseController) BeforeSql(data map[string]interface{}) {

}
func (self *DatabaseController) AfterSql(data map[string]interface{}, oldinfo orm.Params) {
	if self.method == "Del" {
		manger := libs.GetManger()
		host := beego.AppConfig.String("qiniu.host")
		bucket := beego.AppConfig.String("qiniu.bucket")
		key := strings.TrimPrefix(data["path"].(string), host)
		logs.Info("del key:%s", key)
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

var savefilename = "databackup"

// var fileMaxSize = 10485760 //最大的文件
var insert_max_num = 500

func (self *DatabaseController) Add() {
	self.CheckFieldExit(self.postdata, "name", "名字不能为空")
	path, err := SaveDatabase()
	if err != nil {
		logs.Error(err.Error())
		self.AjaxReturnError(err.Error())
	}
	adddata := make(map[string]interface{})
	adddata["name"] = self.postdata["name"]
	adddata["user_id"] = self.uid
	adddata["build_time"] = time.Now().Unix()
	adddata["path"] = path
	self.AddCommonExe(self, adddata)
}

func (self *DatabaseController) Edit() {
	self.CheckFieldExit(self.postdata, "name", "名字不能为空")
	self.CheckFieldExit(self.postdata, "id", "id空")
	changedata := make(map[string]interface{})
	changedata["name"] = self.postdata["name"]
	self.updateSqlById(self, changedata, self.postdata["id"])
}

func (self *DatabaseController) Del() {
	self.DelCommon(self)
}

//还原
func (self *DatabaseController) Restore() {
	self.CheckFieldExit(self.postdata, "id", "id空")
	info := self.model.GetInfoById(self.postdata["id"])
	if info == nil {
		self.AjaxReturnError("找不到")
	}
	err := RestoreDatabase(info["path"].(string))
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	self.AjaxReturnSuccess("", nil)
}

func getSqlfilePath() string {
	tempfolder := beego.AppConfig.String("site.tempfolder")
	return tempfolder + savefilename + ".sql"
}

//保存数据库
func SaveDatabase() (string, error) {
	var filepatharr []string
	// var fileindex = 1
	filepath := getSqlfilePath()
	tempfolder := beego.AppConfig.String("site.tempfolder")
	fileio, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	defer fileio.Close()
	filepatharr = append(filepatharr, filepath)
	tablelist, err := libs.GetTableList()
	if err != nil {
		return "", err
	}
	// var sqlstring bytes.Buffer

	fileio.WriteString("SET FOREIGN_KEY_CHECKS=0;")
	fileio.WriteString(libs.SqlLineEnd)
	for _, tableinfo := range tablelist {

		tablename := tableinfo["Name"].(string)
		logs.Info("get tableinfo:%s", tablename)

		//忽略的表
		if tablename == "aq_log" {
			continue
		}
		if tablename == "aq_database" {
			continue
		}
		fileio.WriteString("-- ----------------------------")
		fileio.WriteString(libs.SqlLineEnd)
		fileio.WriteString("-- Table structure for " + tablename)
		fileio.WriteString(libs.SqlLineEnd)
		fileio.WriteString("-- ----------------------------")
		fileio.WriteString(libs.SqlLineEnd)
		fileio.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", tablename))
		fileio.WriteString(libs.SqlLineEnd)
		tablestr, err := libs.GetTableString(tablename)
		if err != nil {
			return "", err
		}
		fileio.WriteString(tablestr)
		fileio.WriteString(libs.SqlLineEnd)

		//行
		fileio.WriteString("-- ----------------------------")
		fileio.WriteString(libs.SqlLineEnd)
		fileio.WriteString("-- Records of " + tablename)
		fileio.WriteString(libs.SqlLineEnd)
		fileio.WriteString("-- ----------------------------")
		fileio.WriteString(libs.SqlLineEnd)

		var dataList []orm.Params
		db := orm.NewOrm()
		_, err = db.Raw("select count(*) as countnum from " + tablename).Values(&dataList)
		if err != nil {
			return "", err
		}

		totalnumstr := dataList[0]["countnum"].(string)
		totalrow, err := strconv.Atoi(totalnumstr)
		if err != nil {
			return "", err
		}
		pagenum := (totalrow / insert_max_num) + 1
		for curpage := 0; curpage < pagenum; curpage++ {
			startrow := insert_max_num * curpage
			rowstr, err := libs.GetInsertSql(tablename, startrow, insert_max_num)
			if err != nil {
				return "", err
			}
			fileio.WriteString(rowstr)
			// if sqlstring.Len() > fileMaxSize {
			// 	logs.Info("new file")
			// 	_, err := fileio.WriteString(sqlstring.String())
			// 	if err != nil {
			// 		return "", err
			// 	}
			// 	sqlstring.Reset()
			// 	fileio.Close()
			// 	//新生成一个文件
			// 	filepath := getSqlfilePath(fileindex)
			// 	fileindex++
			// 	fileio, err = os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			// 	if err != nil {
			// 		return "", err
			// 	}
			// 	defer fileio.Close()
			// 	filepatharr = append(filepatharr, filepath)
			// }

		}
	}
	// _, err = fileio.WriteString(sqlstring.String())
	// if err != nil {
	// 	return "", err
	// }
	// sqlstring.Reset()
	fileio.Close()

	// logs.Info("begin zip file")
	zippath := tempfolder + "databackall.zip"
	err = libs.Compress(filepatharr, zippath)
	if err != nil {
		return "", err
	}

	for _, filepath := range filepatharr {
		os.Remove(filepath)
	}

	filemd5str := libs.GetFileMd5(zippath)
	bucket := beego.AppConfig.String("qiniu.bucket")
	host := beego.AppConfig.String("qiniu.host")
	url := host + filemd5str + ".zip"
	logs.Info("update filepath:%s", zippath)
	_, err = libs.UploadFile(filemd5str+".zip", zippath, bucket)
	os.Remove(zippath)
	if err != nil {
		return "", err
	}

	return url, nil
}

func RestoreDatabase(path string) error {
	tempfolder := beego.AppConfig.String("site.tempfolder")
	outpath := tempfolder + "download_tablebase.zip"
	releasepath := tempfolder + "download_tablebaseout/"
	res, err := http.Get(path)
	if err != nil {
		return err
	}
	fileio, err := os.Create(outpath)
	defer fileio.Close()
	if err != nil {
		return err
	}
	io.Copy(fileio, res.Body)
	err = libs.DeCompress(outpath, releasepath)
	if err != nil {
		return err
	}
	filetemp, err := os.Open(releasepath)
	if err != nil {
		return err
	}
	defer filetemp.Close()
	fileInfos, err := filetemp.Readdir(-1)
	if err != nil {
		return err
	}
	//读sql文件
	db := orm.NewOrm()
	db.Begin()
	for _, fi := range fileInfos {
		f, err := os.Open(releasepath + fi.Name())
		if err != nil {
			return err
		}
		defer f.Close()
		br := bufio.NewReader(f)
		var buffersql bytes.Buffer
		for {
			line, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			if len(line) <= 2 {
				continue
			}
			if (line[0] == '-') && (line[1] == '-') {
				continue
			}
			// logs.Info("write:%s", line)
			buffersql.Write(line)
			linestr := strings.TrimSpace(string(line))
			if linestr[len(linestr)-1] == ';' {
				//结束
				// logs.Error("get one")
				sqltext := buffersql.String()
				buffersql.Reset()
				_, err := db.Raw(sqltext).Exec()
				if err != nil {
					db.Rollback()
					return err
				}

			}
		}
	}
	db.Commit()
	return nil
}
