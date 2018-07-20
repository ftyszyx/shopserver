package libs

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
)

//获取token
func GetToken(gettime interface{}, uid interface{}, password interface{}, groupid interface{}) string {
	str := fmt.Sprintf("%v%v%v%v", gettime, uid, groupid, password) //将[]byte转成16进制
	// logs.Info("token src str:%s", str)
	return GetStrMD5(str)
}

func GetFileMd5(path string) string {
	file, inerr := os.Open(path)
	defer file.Close()
	if inerr == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		return fmt.Sprintf("%x", md5h.Sum(nil))
	}
	return ""
}

func GetStrMD5(src string) string {
	return GetByteMD5([]byte(src))
}

func GetByteMD5(src []byte) string {
	cipherStr := md5.Sum(src)
	md5str1 := fmt.Sprintf("%x", cipherStr) //将[]byte转成16进制

	return md5str1
}

//发包
func AjaxReturn(self *beego.Controller, code int, msg interface{}, data interface{}) {
	out := make(map[string]interface{})
	if msg != nil {
		out["message"] = msg.(string)
	}
	out["code"] = code
	out["data"] = data
	self.Data["json"] = out
	self.ServeJSON()
	self.StopRun()
}

func AjaxReturnError(self *beego.Controller, msg interface{}) {
	AjaxReturn(self, ErrorCode, msg, nil)
}

func AjaxReturnSuccess(self *beego.Controller, msg interface{}, data interface{}) {
	AjaxReturn(self, SuccessCode, msg, data)
}

func StructToMapCmp(in interface{}, tag string, changemap map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tag == "" {
			out[fi.Name] = v.Field(i).Interface()
		} else {
			if tagv := fi.Tag.Get(tag); tagv != "" {
				// set key of map to value in struct field
				if tag == "edit" && changemap != nil {
					_, have := changemap[tagv]
					if have {
						out[tagv] = v.Field(i).Interface()
					}
				} else {
					out[tagv] = v.Field(i).Interface()
				}

			}
		}

	}
	return out, nil
}

func StructToMap(in interface{}, tag string) (map[string]interface{}, error) {
	return StructToMapCmp(in, tag, nil)
}

//按照结构体清除map无用字段,只有在结构体中有的字段，才会放进去
func ClearMapByStructTag(data map[string]interface{}, in interface{}, tag string) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	typ := v.Type()
	if data == nil {
		return data
	}
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		lowName := strings.ToLower(fi.Name)
		if tag == "" {
			_, have := data[lowName]
			if have {
				out[lowName] = data[lowName]
			}
		} else {
			if tagv := fi.Tag.Get(tag); tagv != "" {
				_, have := data[tagv]
				if have {
					out[tagv] = data[tagv]
				}

			}
		}

	}
	return out

}

func ClearMapByStruct(data map[string]interface{}, in interface{}) map[string]interface{} {
	return ClearMapByStructTag(data, in, "")
}
