package libs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
)

type Kuai100Data struct {
	Context string
	Time    string
	Ftime   string
}

type Kuai100Resp struct {
	Message   string
	State     string
	Status    string
	Condition string
	Ischeck   string
	Com       string
	Nu        string
	Data      []Kuai100Data
}

type Kuai00SendParam struct {
	Com      string `json:"com"`
	Num      string `json:"num"`
	From     string `json:"from"`
	To       string `json:"to"`
	Resultv2 int    `json:"resultv2"`
}

func GetKuai100Info(costomerid string, key string, company string, num string) (*Kuai100Resp, error) {
	var param Kuai00SendParam
	param.Com = company
	param.Num = num
	parambuf, err := json.Marshal(param)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	paramstr := string(parambuf)

	var sigContent = paramstr + key + costomerid
	signstr := strings.ToUpper(GetStrMD5(sigContent))

	url := fmt.Sprintf("https://poll.kuaidi100.com/poll/query.do?customer=%s&sign=%s&param=%s", costomerid, signstr, paramstr)
	req := httplib.Get(url)
	var respdata []byte
	logs.Info("url:%s", url)
	respdata, err = req.Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	getData := new(Kuai100Resp)
	logs.Info("get data:%s", string(respdata))
	err = json.Unmarshal(respdata, getData)
	if err != nil {
		logs.Info("parse data err")
		return nil, errors.WithStack(err)
	}
	// logs.Info("get data:%v", getData)
	if getData.Status != "200" {
		return nil, errors.New(getData.Message)
	}
	return getData, nil
}
