package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

type User struct {
	Model
}

type UserData struct {
	Account    string `empty:"账号名不能为空" only:"账号名重复"`
	Name       string `empty:"姓名不能为空" only:"用户名重复"`
	Mail       string `empty:"邮箱不能为空" only:"邮箱重复"`
	Phone      string `empty:"手机不能为空" only:"手机重复"`
	User_group string `empty:"用户组不能为空"`
}

func (self *User) GetModelStruct() interface{} {
	return UserData{}
}

func (self *User) GetFieldName(name string) string {
	return "user." + name
}

func (self *User) InitSqlField(sql libs.SqlType) libs.SqlType {

	return self.InitField(self.InitJoinString(sql, true))
}
func (self *User) InitJoinString(sql libs.SqlType, allfield bool) libs.SqlType {
	groupname := GetModel(USERGROUP).TableName()

	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("group") == true) {
		fieldstr += fmt.Sprintf("left join `%s` `group` ON `group`.`id`=`user`.`user_group`", groupname)
	}
	return sql.Alias("user").Join(fieldstr)
}

func (self *User) InitField(sql libs.SqlType) libs.SqlType {
	return sql.Field(map[string]string{
		"user.id":              "id",
		"user.name":            "name",
		"user.head":            "head",
		"user.account":         "account",
		"user.mail":            "mail",
		"user.reg_time":        "reg_time",
		"user.last_login_time": "last_login_time",
		"user.wchat_openid":    "wchat_openid",
		"user.phone":           "phone",
		"user.password":        "password",
		"user.is_del":          "is_del",
		"user.user_group":      "user_group",
		"user.shop_cart":       "shop_cart",
		"user.token_expire":    "token_expire",
		"group.name":           "user_group_name",
		"group.group_type":     "user_group_type",
	})
}

//验证接口
func (self *User) Auth(token string, uid string, control string, method string) (bool, string, int) {
	info := self.GetInfoAndCache(uid, false)
	module := GetModel("module").(*Module)
	usergroup := GetModel(USERGROUP).(*UserGroup)
	if info != nil {
		//logs.Info("info:%v", info)
		if info["is_del"].(string) == "1" {
			return false, "账号禁用", libs.AuthFail
		}

		expire_time, _ := strconv.ParseInt(info["token_expire"].(string), 10, 64)
		if expire_time < time.Now().Unix() {
			return false, "token失效", libs.AuthFail
		} else {
			usertoken := info["user_token"].(string)
			if usertoken != token {
				logs.Info(fmt.Sprintf("user token not same :%s sqltoken:%s", token, usertoken))
				return false, "token错误", libs.AuthFail
			} else {
				moduleinfo := module.GetModuleInfo(control, method)
				moduleid := moduleinfo["id"].(string)
				// logs.Info("module id:%s", moduleid)
				group := usergroup.GetInfoAndCache(info["user_group"].(string), false)
				if group == nil {
					logs.Info(fmt.Sprintf("user have no group"))
					return false, "group查不到", libs.NoAccessRight
				} else {
					if group["group_type"].(string) == strconv.Itoa(libs.UserSystem) {
						logs.Info("is system")
						return true, "", libs.SuccessCode
					} else {
						moduleids, ok := group["module_ids"].(string)
						if ok {
							var idarr []string
							err := json.Unmarshal([]byte(moduleids), &idarr)
							if err != nil {
								return false, "无权限", libs.NoAccessRight
							}
							// logs.Info("module id:%s,idarr id:%+v", moduleid, idarr)
							for _, v := range idarr {
								// logs.Info("v:%s moduleid:%s %v", v, moduleid, v == moduleid)
								if v == moduleid {
									// logs.Info("find:%s", moduleid)
									return true, "", libs.SuccessCode
								}
							}
						}
						// logs.Info("not find:%s", moduleid)
						return false, "无权限", libs.NoAccessRight

					}
				}
			}
		}

	}
	return false, "账号不存在", libs.AuthFail
}

//延长token过期时间
func (self *User) ExtendExpireTime(uid string, expiretime int) bool {

	o := orm.NewOrm()
	tokenTime := time.Now().Unix()
	expirtTime := tokenTime + int64(expiretime)
	_, err := o.Raw(fmt.Sprintf(`update %s set token_expire=%d where id='%s'`, self.tablename, expirtTime, uid)).Exec()
	if err == nil {
		self.ClearRowCache(uid) //更新缓存
		logs.Error("extend expire ok")
		return true
	}
	logs.Error("extend expire error", err)
	return false
}
