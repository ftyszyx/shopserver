package libs

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/axgle/mahonia"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func CheckPhone(phone string) bool {
	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

func GbkToUtf8(src []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(src), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func UTF82GBK(src string) (string, error) {
	// reader := transform.NewReader(strings.NewReader(src), simplifiedchinese.GBK.NewEncoder())
	// if buf, err := ioutil.ReadAll(reader); err != nil {
	// 	return "", err
	// } else {
	// 	return string(buf), nil
	// }
	enc := mahonia.NewEncoder("gbk")
	//converts a  string from UTF-8 to gbk encoding.
	return enc.ConvertString(src), nil
}

func GetStrArr(src []interface{}) []string {
	var strarr []string
	for _, item := range src {
		strarr = append(strarr, item.(string))
	}
	return strarr
}

func FormatTableTime(stamp string) (datestring string) {
	if stamp == "" || stamp == "0" {
		return ""
	}
	return FormatTime(stamp, "Y-m-d H:i:s")
}

func FormatTime(stamp string, format string) (datestring string) {
	timeint, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		logs.Error("err:%s", err.Error())
		return
	}
	timeinfo := time.Unix(timeint, 0)
	datestring = beego.Date(timeinfo, format)
	return
}

var (
	idnum_coefficient []int32 = []int32{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	idnum_code        []byte  = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
)

func CheckIdNum(idCardNo string) bool {
	if len(idCardNo) != 18 {
		return false
	}

	idByte := []byte(strings.ToUpper(idCardNo))

	sum := int32(0)
	for i := 0; i < 17; i++ {
		sum += int32(byte(idByte[i])-byte('0')) * idnum_coefficient[i]
	}
	return idnum_code[sum%11] == idByte[17]
}

func addZero(timestring string) string {
	if len([]rune(timestring)) == 1 {
		return "0" + timestring
	}
	return timestring
}
func ParseTime(timestr string) (error, time.Time) {
	re, err := regexp.Compile(`\s+`)
	if err != nil {
		return errors.WithStack(err), time.Time{}
	}
	areaarr := re.Split(timestr, 2)
	if len(areaarr) != 2 {
		// logs.Info("areaarr:%+v", areaarr)
		return errors.New("fomart err"), time.Time{}
	}
	datearr := strings.Split(areaarr[0], "/")
	if len(datearr) != 3 {
		return errors.New("date fomart err"), time.Time{}
	}

	datearr[1] = addZero(datearr[1])
	datearr[2] = addZero(datearr[2])

	timearr := strings.Split(areaarr[1], ":")
	if len(timearr) == 2 {
		timearr = strings.Split(areaarr[1]+":00", ":")
	}
	if len(timearr) != 3 {
		return errors.New("time fomart err"), time.Time{}
	}
	timearr[0] = addZero(timearr[0])
	timearr[1] = addZero(timearr[1])
	timearr[2] = addZero(timearr[2])
	newstr := strings.Join(datearr, "/") + " " + strings.Join(timearr, ":")
	starttime, err := beego.DateParse(newstr, "Y/m/d H:i:s")
	if err != nil {
		return errors.WithStack(err), time.Time{}
	}

	return nil, starttime
}
