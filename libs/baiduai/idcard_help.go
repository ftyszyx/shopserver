package baiduai

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
)

//识别身份证
type IDDataPos struct {
	Left   uint32
	Top    uint32
	Width  uint32
	Height uint32
}

type IDDataResult struct {
	location IDDataPos
	Words    string
}

type IDDataResp struct {
	Direction        int32
	Image_status     string
	Risk_type        string
	Edit_tool        string
	Log_id           uint64
	Words_result     map[string]IDDataResult
	Words_result_num uint32
}

var IDErrMesage = map[string]string{
	"reversed_side":   "身份证正反面颠倒",
	"non_idcard":      "上传的图片中不包含身份证",
	"blurred":         "身份证模糊",
	"other_type_card": "其他类型证照",
	"over_exposure":   "身份证关键字段反光或过曝",
	"over_dark":       "身份证欠曝（亮度过低）",
	"unknown":         "未知错误",
}

var IDRishErrMessage = map[string]string{
	"copy":      "复印件",
	"temporary": "临时身份证",
	"screen":    "翻拍",
	"unknown":   "未知情况",
}

//获取身份证结果
func (self *BaiduAiType) GetIdResByPath(filepath string, cardside string) (*IDDataResp, error) {
	ff, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return self.GetIdResByData(ff, cardside)

}

func (self *BaiduAiType) GetIdResByData(imagdata []byte, cardside string) (*IDDataResp, error) {
	bufstore := make([]byte, base64.StdEncoding.EncodedLen(len(imagdata)))
	base64.StdEncoding.Encode(bufstore, imagdata)
	var sendata = make(map[string]interface{})
	sendata["detect_direction"] = false
	sendata["id_card_side"] = cardside
	sendata["image"] = string(bufstore)
	sendata["detect_risk"] = true
	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/idcard"

	resp, err := self.SendRequest(url, sendata, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var getData IDDataResp
	err = json.Unmarshal(resp, &getData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if getData.Image_status != "normal" {
		//错误
		logs.Error(IDErrMesage[getData.Image_status])
		if getData.Image_status != "unknown" {
			return nil, errors.New(IDErrMesage[getData.Image_status])
		}
	}

	if getData.Risk_type != "normal" {
		//return nil, errors.New(IDRishErrMessage[getData.Risk_type])
		logs.Error(IDRishErrMessage[getData.Risk_type])
		if getData.Risk_type == "copy" {
			return nil, errors.New(IDRishErrMessage[getData.Risk_type])
		}
	}

	if getData.Edit_tool != "" {
		logs.Error(getData.Edit_tool + "编辑过")
	}
	return &getData, nil

}

//检查url地址的身份证信息
func (self *BaiduAiType) GetIdResByUrl(url string, cardside string) (*IDDataResp, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return GlobalaiIDCarddata.GetIdResByData(data, cardside)

}
