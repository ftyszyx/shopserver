package admin

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/baiduai"

	"github.com/astaxie/beego"
)

type UploadController struct {
	BaseController
}

func (self *UploadController) Get() {
	action := self.Input().Get("action")
	switch action {
	case "config": //这里是conf/config.json
		file, err := os.Open("conf/config.json")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer file.Close()
		fd, err := ioutil.ReadAll(file)
		src := string(fd)
		re, _ := regexp.Compile("\\/\\*[\\S\\s]+?\\*\\/") //参考php的$CONFIG = json_decode(preg_replace("/\/\*[\s\S]+?\*\//", "", file_get_contents("config.json")), true);

		src = re.ReplaceAllString(src, "")
		tt := []byte(src)
		var r interface{}
		json.Unmarshal(tt, &r) //这个byte要解码
		self.Data["jsonp"] = r
		self.ServeJSONP()
		self.StopRun()
	case "listimage":
	case "listfile":
	case "catchimage":
		sourcearr := self.Input().Get("source")
		logs.Info("srouce:%+v", sourcearr)
	}
}

func (self *UploadController) Post() {
	op := self.Input().Get("action")
	logs.Info("Post:%s", op)
	switch op {
	case "uploadimage":
		self.editorUpload()
	case "uploadvideo":
		self.editorUpload()
	case "uploadfile":
		self.editorUpload()
	}
}

func (self *UploadController) PicUpload() {
	logs.Info("picupload")
	err, fileinfo := self.upload()
	logs.Info("picupload:%+v", fileinfo)
	if err == nil {
		err = os.Remove(fileinfo["filePath"].(string))
		self.AjaxReturnSuccess("", fileinfo)
	} else {
		logs.Info("PicUpload error:%s", err.Error())
		self.AjaxReturnError(errors.WithStack(err))
	}
}

//上传身份证
func (self *UploadController) UploadIDNum() {
	// logs.Info("uploadIDNum")
	side := self.GetString("side", "")
	logs.Info("uploadIDNum:%s", side)

	if side != "front" && side != "back" {
		self.AjaxReturnError(errors.New("格式不对"))
	}
	//保存
	err, fileinfo := self.saveUploadFile()
	if err != nil {
		logs.Info("savefile err %+v", err)
		self.AjaxReturnError(errors.WithStack(err))
	}
	logs.Info("save ok")
	//识别
	cardres, err := baiduai.GlobalaiIDCarddata.GetIdResByPath(fileinfo["filePath"].(string), side)
	if err != nil {
		logs.Info("Globalaidata err %+v", err)
		self.AjaxReturnError(errors.WithStack(err))
	}

	bucket := beego.AppConfig.String("qiniu.bucket")
	_, err = libs.UploadFile(fileinfo["filename"].(string), fileinfo["filePath"].(string), bucket)

	if err == nil {
		err = os.Remove(fileinfo["filePath"].(string))
		fileinfo["result"] = cardres.Words_result
		fileinfo["side"] = side
		fileinfo["res_num"] = cardres.Words_result_num
		self.AjaxReturnSuccess("", fileinfo)
	} else {
		logs.Info("PicUpload error:%+v", err.Error())
		self.AjaxReturnError(errors.WithStack(err))
	}
}

func (self *UploadController) editorUpload() {
	err, fileinfo := self.upload()
	if err == nil {
		err = os.Remove(fileinfo["filePath"].(string))
		self.Data["json"] = map[string]interface{}{
			"state":    "SUCCESS",
			"url":      fileinfo["url"],
			"title":    fileinfo["filetitle"],
			"original": fileinfo["filetitle"],
			"size":     fileinfo["filesize"],
			"type":     fileinfo["filetype"],
		}
		logs.Info("send %+v ", self.Data["json"])
		self.ServeJSON()
		self.StopRun()
	} else {
		logs.Info("error:%s", err.Error())
		self.Data["json"] = map[string]interface{}{"state": "error"}
		self.ServeJSON()
		self.StopRun()
	}
}

func (self *UploadController) uploadScrawl() {
	tempfolder := beego.AppConfig.String("site.tempfolder")
	ww := self.Input().Get("upfile")
	ddd, _ := base64.StdEncoding.DecodeString(ww)
	newname := strconv.FormatInt(time.Now().Unix(), 10)           // + "_" + filename
	err := ioutil.WriteFile(tempfolder+newname+".jpg", ddd, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		beego.Error(err)
	}
}
