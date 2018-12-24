package libs

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

func getRandom() string {
	rand.Seed(int64(time.Now().UnixNano()))
	return strconv.Itoa(rand.Intn(899999) + 100000)
}

func GetSign(moblie string, reqtime int64, random string) string {
	var t = strconv.FormatInt(reqtime, 10)
	appkey := beego.AppConfig.String("qqmsg.appkey")
	var sigContent = "appkey=" + appkey + "&random=" + random + "&time=" + t

	if len(moblie) > 0 {
		sigContent += "&mobile=" + moblie
	}
	logs.Info("sign info:%s", sigContent)
	h := sha256.New()
	h.Write([]byte(sigContent))

	return fmt.Sprintf("%x", h.Sum(nil))
}

type SMSResult struct {
	Result uint   `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Sid    string `json:"sid,omitempty"`
	Fee    uint   `json:"fee,omitempty"`
}

func SendQQMsg(phone string, code string) error {
	appid := beego.AppConfig.String("qqmsg.appid")
	codemsg := beego.AppConfig.String("qqmsg.codemsg")
	random := getRandom()
	urlstr := beego.AppConfig.String("qqmsg.url") + "?sdkappid=" + appid + "&random=" + random

	req := httplib.Post(urlstr)
	sendata := make(map[string]interface{})
	reqtime := time.Now().Unix()
	sendata["time"] = reqtime
	sendata["type"] = 0
	sendata["sig"] = GetSign(phone, reqtime, random)
	sendata["msg"] = fmt.Sprintf(codemsg, code)
	sendata["tel"] = map[string]string{"mobile": phone, "nationcode": "86"}

	reqbuf, err := json.Marshal(sendata)
	if err != nil {
		return err
	}

	logs.Info("send data:%s", string(reqbuf))
	logs.Info("send url:%s", urlstr)

	req.Body(string(reqbuf))
	req.Header("Content-Type", "application/json")

	respdata, err := req.Bytes()
	if err != nil {
		return err
	}
	getData := new(SMSResult)
	json.Unmarshal(respdata, getData)

	if getData.Result == 0 {
		return nil
	}
	return errors.New(getData.Errmsg + ":" + strconv.Itoa(int(getData.Result)))

}
