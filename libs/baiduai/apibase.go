package baiduai

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/httplib"
	"github.com/zyx/shop_server/libs"
)

var accessTokenUrl = "https://aip.baidubce.com/oauth/2.0/token"
var reportUrl = "https://aip.baidubce.com/rpc/2.0/feedback/v1/report"

var apiCache cache.Cache

type BaiduAiType struct {
	Appid       string
	Apikey      string
	Secretkey   string
	IsClouduser bool
	Scope       string
}

type AuthResp struct {
	Refresh_token  string
	Expires_in     int64
	Scope          string
	Session_key    string
	Access_token   string
	Session_secret string
	Ispermission   bool
	GetTime        int64
}

type DataRespErr struct {
	Error_code int
	Error_msg  string
}

type TokenRespErr struct {
	Error             string
	Error_description string
}

var GlobalaiIDCarddata *BaiduAiType

func InitBaiduAiIDcard() {
	logs.Info("init idcard baiduai")
	appid := beego.AppConfig.String("baiduai.appid")
	appkey := beego.AppConfig.String("baiduai.appkey")
	secretkey := beego.AppConfig.String("baiduai.secretkey")
	GlobalaiIDCarddata = NewIO(appid, appkey, secretkey, "brain_ocr_idcard")
}

func init() {
	apiCache, _ = cache.NewCache("memory", `{"interval":0}`) //不过期
}

func NewIO(appid string, apikey string, secretkey string, scope string) *BaiduAiType {
	o := new(BaiduAiType)
	o.Apikey = apikey
	o.Appid = appid
	o.Secretkey = secretkey
	o.Scope = scope
	return o
}

type BaiduAiIO interface {
	SendRequest(string, map[string]interface{}, bool) ([]byte, error)
}

func (self *BaiduAiType) SendRequest(url string, data map[string]interface{}, tokenfresh bool) ([]byte, error) {

	authdata, err := self.Auth(tokenfresh)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	urlstr := url + "?access_token=" + authdata.Access_token
	logs.Info("send url:%s", urlstr)
	req := httplib.Post(urlstr)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	req.Header("Content-Type", "application/x-www-form-urlencoded")

	for key, value := range data {
		req.Param(key, fmt.Sprintf("%+v", value))
		if key != "image" {
			logs.Info("key:%s value:%+v", key, value)
		}

	}

	respdata, err := req.Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	logs.Info("get res:%s", string(respdata))
	var getData map[string]interface{}
	err = json.Unmarshal(respdata, &getData)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if getData["error_code"] == nil {
		return respdata, nil
	} else {
		var errinfo DataRespErr
		err = json.Unmarshal(respdata, &errinfo)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if tokenfresh == false {
			if errinfo.Error_code == 110 {
				//token失效
				return self.SendRequest(url, data, true)
			}
		}
		return nil, errors.Errorf("errcode:%d error_msg:%s", errinfo.Error_code, errinfo.Error_msg)

	}
}

func (self *BaiduAiType) getAuthObj() *AuthResp {

	res := apiCache.Get(self.Apikey)
	if res != nil {
		resdata, ok := res.(AuthResp)
		if ok == false {
			return nil
		}
		if resdata.GetTime+resdata.Expires_in > time.Now().Unix() {
			//还没过期
			if self.isPermission(&resdata) == false {
				logs.Error("scope not permission")
				return nil
			}
			return &resdata
		} else {
			logs.Error("baiduai is expire")
		}

	}
	return nil
}

func (self *BaiduAiType) Auth(refresh bool) (*AuthResp, error) {

	if refresh == false {
		authdata := self.getAuthObj()
		if authdata != nil {
			return authdata, nil
		}
	}

	urlstr := accessTokenUrl + "?grant_type=client_credentials" + "&client_id=" + self.Apikey + "&client_secret=" + self.Secretkey
	logs.Info("Auth url:%s", urlstr)
	req := httplib.Post(urlstr)
	req.Body("")
	req.Header("Content-Type", "application/json")

	respdata, err := req.Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var datamap map[string]interface{}
	logs.Info("respse:%s", string(respdata))
	err = json.Unmarshal(respdata, &datamap)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if datamap["error"] != nil {
		return nil, errors.Errorf("get token err:%+v", datamap)
	}
	var getData AuthResp
	err = json.Unmarshal(respdata, &getData)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	getData.Ispermission = self.isPermission(&getData)
	getData.GetTime = time.Now().Unix()

	if getData.Ispermission == false {
		return nil, errors.New("not permission")
	} else {
		apiCache.Put(self.Apikey, getData, 0)
	}
	return &getData, nil

}

func (self *BaiduAiType) isPermission(data *AuthResp) bool {
	if data == nil || data.Scope == "" {
		return false
	}

	scopelist := strings.Split(data.Scope, " ")
	getindex := libs.SliceIndex(scopelist, func(item string) bool {
		return item == self.Scope
	})
	if getindex == -1 {
		return false
	}
	return true
}
