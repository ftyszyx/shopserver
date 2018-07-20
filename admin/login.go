package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/wechat"
)

type LoginController struct {
	BaseController
}

//登录成功
func (self *LoginController) loginSucesss(userinfo orm.Params, changeinfo map[string]interface{}, noreturn bool) {
	logs.Info("login ok")
	userModel := models.GetModel(models.USER)
	userGroupModel := models.GetModel(models.USERGROUP)
	curtime := time.Now().Unix()
	logintime := time.Now().Unix()
	id := userinfo["id"].(string)
	senddata := make(map[string]interface{})
	expire_time, _ := strconv.ParseInt(userinfo["token_expire"].(string), 10, 64)
	if expire_time < curtime {
		//会过期
		logs.Info("user token is expire")
		usertoken := libs.GetToken(curtime, id, userinfo["password"], userinfo["user_group"].(string))
		changeinfo["user_token"] = usertoken
		changeinfo["token_get_time"] = curtime
		senddata["token"] = usertoken
	} else {
		//没过期
		logs.Info("user token is ok")
		changeinfo["user_token"] = userinfo["user_token"]
		senddata["token"] = userinfo["user_token"]
	}
	groupinfo := userGroupModel.GetInfoAndCache(userinfo["user_group"].(string), false)
	expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	changeinfo["token_expire"] = curtime + int64(expiretime)
	changeinfo["last_login_time"] = logintime

	senddata["uid"] = id
	o := orm.NewOrm()
	_, err = o.Raw(fmt.Sprintf(`update  %s set %s  where id=?`, userModel.TableName(), libs.SqlGetKeyValue(changeinfo, "=")), id).Exec()
	if err == nil {
		userModel.GetInfoAndCache(id, true) //更新缓存
		if noreturn == false {
			self.AjaxReturnSuccess("登录成功", senddata)
		}
	} else {
		self.AjaxReturnError(err.Error())
	}

}

//新增用户
func (self *LoginController) addUser(adddata map[string]interface{}, fieldkey string, fieldvalue string, noreturn bool) orm.Params {
	logs.Info("add user")
	var res []orm.Params
	userGroupModel := models.GetModel(models.USERGROUP)
	userModel := models.GetModel(models.USER)
	o := orm.NewOrm()

	num, err := o.Raw(fmt.Sprintf(`select * from %s where group_type=? limit 1`, userGroupModel.TableName()), libs.UserMember).Values(&res)
	if err != nil && num == 0 {
		self.AjaxReturnError("用户组不存在或错误")
	}
	groupinfo := res[0]
	password := string(utils.RandomCreateBytes(13))
	curtime := time.Now().Unix()
	adddata["password"] = libs.GetStrMD5(password)
	adddata["reg_time"] = curtime
	adddata["user_group"] = groupinfo["id"]
	adddata["token_get_time"] = curtime
	expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	adddata["token_expire"] = curtime + int64(expiretime)

	adddata["user_token"] = libs.GetToken(adddata["token_expire"], groupinfo["id"], adddata["password"], adddata["user_group"].(string))
	adddata["last_login_time"] = curtime

	keys, values := libs.SqlGetInsertInfo(adddata)
	_, err = o.Raw(fmt.Sprintf("insert into %s (%s) values (%s)", userModel.TableName(), keys, values)).Exec()
	if err != nil {
		self.AjaxReturnError(err.Error())
	}

	num, err = o.Raw(fmt.Sprintf(`select * from %s where %s=? limit 1`, userModel.TableName(), fieldkey), fieldvalue).Values(&res)
	if err == nil && num > 0 {
		// userModel.ClearRowCache(res[0]["id"].(string))
		if noreturn == false {
			data := make(map[string]interface{})
			data["token"] = adddata["user_token"]
			data["uid"] = res[0]["id"]

			self.AjaxReturnSuccess("登录成功", data)
		}
		return res[0]
	} else {
		self.AjaxReturnError(err.Error())

	}
	return nil

}

type loginPost struct {
	Username string
	Password string
}

func (self *LoginController) Login() {
	var data = new(loginPost)
	json.Unmarshal(self.Ctx.Input.RequestBody, data)
	logs.Info("login username:%s pass:%s ", data.Username, data.Password)
	o := orm.NewOrm()
	var res []orm.Params
	passMd5 := libs.GetStrMD5(data.Password)
	userModel := models.GetModel(models.USER)
	num, err := o.Raw(fmt.Sprintf(`select * from %s where account=? and password=? limit 1`, userModel.TableName()), data.Username, passMd5).Values(&res)
	if err == nil && num > 0 {
		self.loginSucesss(res[0], make(map[string]interface{}), false)
	}
	libs.AjaxReturn(&self.Controller, libs.ErrorCode, "账号或密码错误", nil)

}

func (self *LoginController) LoginOut() {
	logs.Info("loginout")
	o := orm.NewOrm()
	userModel := models.GetModel(models.USER)
	_, err := o.Raw(fmt.Sprintf(`update %s set user_token='%s',token_expire='%d' where id='%s'`, userModel.TableName(), "", 0, self.uid)).Exec()
	if err == nil {
		userModel.GetInfoAndCache(self.uid, true) //更新缓存
		libs.AjaxReturn(&self.Controller, libs.SuccessCode, "登出成功", nil)
	}

	libs.AjaxReturn(&self.Controller, libs.ErrorCode, "登录失败", nil)
}

//获取验证码
func (self *LoginController) GetCaptchaCode() {
	codeid, err := models.CaptchaCode.CreateCaptcha()
	senddata := make(map[string]interface{})
	if err == nil {
		senddata["codeid"] = codeid
		libs.AjaxReturnSuccess(&self.Controller, "", senddata)
	}

	libs.AjaxReturnError(&self.Controller, err.Error())
}

type captchaData struct {
	Captcha_id string
	Captcha    string
	Phone      string
}

//获取手机验证码
func (self *LoginController) GetPhoneCode() {
	var data = new(captchaData)
	json.Unmarshal(self.Ctx.Input.RequestBody, data)
	if data.Phone == "" {
		self.AjaxReturnError("手机号不能为空")
	}
	logs.Info("get code:%+v", data)
	if models.CaptchaCode.Verify(data.Captcha_id, data.Captcha) == false {
		self.AjaxReturnError("验证码错误")
	}
	codestr := string(utils.RandomCreateBytes(6))
	//发送验证码
	codestr = "1111"
	models.PhoneCodeCache.Put(data.Phone, codestr, 120*time.Second)

	err := libs.SendQQMsg(data.Phone, codestr)
	if err == nil {
		self.AjaxReturnSuccess("验证码发送成功", nil)
	}
	self.AjaxReturnError(err.Error())

}

type PhoneLoginData struct {
	Phone string
	Code  string
}

//手机登录
func (self *LoginController) LoginWithPhone() {
	var data = new(PhoneLoginData)
	json.Unmarshal(self.Ctx.Input.RequestBody, data)
	if data.Phone == "" {
		self.AjaxReturnError("手机号不能为空")
	}
	if data.Code == "" {
		self.AjaxReturnError("验证码不能为空")
	}

	codestr, ok := models.PhoneCodeCache.Get(data.Phone).(string)
	logs.Info("get code:%+v getstr:%s", data, codestr)
	if ok == false || codestr == "" || codestr != data.Code {
		self.AjaxReturnError("验证码不对")
	}

	userModel := models.GetModel(models.USER)
	var res []orm.Params
	o := orm.NewOrm()
	num, err := o.Raw(fmt.Sprintf(`select * from %s where phone=? limit 1`, userModel.TableName()), data.Phone).Values(&res)
	if err == nil {
		//是老用户
		if num > 0 {
			logs.Info("old user")
			self.loginSucesss(res[0], make(map[string]interface{}), false)
		} else {
			logs.Info("new user")
			adddata := make(map[string]interface{})
			adddata["phone"] = data.Phone
			adddata["account"] = data.Phone
			self.addUser(adddata, "phone", data.Phone, false)
		}
	}
	self.AjaxReturnError(err.Error())

}

//wchat登录
func (self *LoginController) LoginWithWchat() {
	callbackurl := beego.AppConfig.String("wechat.logincallback")
	url, err := wechat.OauthInstance.GetRedirectURL(callbackurl, "snsapi_userinfo", "test")
	if err != nil {
		logs.Error(err.Error())
		self.AjaxReturnError(err.Error())
	}
	self.AjaxReturnSuccess("", map[string]interface{}{"url": url})
}

//回调
func (self *LoginController) WchatLoginCallback() {
	logs.Info("WchatLoginCallback")
	code := self.Input().Get("code")
	resToken, err := wechat.OauthInstance.GetUserAccessToken(code)
	if err != nil {
		logs.Error(err.Error())
		//self.AjaxReturnError(err.Error())
		return
	}

	userInfo, err := wechat.OauthInstance.GetUserInfo(resToken.AccessToken, resToken.OpenID)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	logs.Info("userinfo:%v", userInfo)

	userModel := models.GetModel(models.USER)
	var res []orm.Params
	o := orm.NewOrm()
	num, err := o.Raw(fmt.Sprintf(`select * from %s where wchat_openid=? limit 1`, userModel.TableName()), userInfo.OpenID).Values(&res)
	var uid string
	var token string
	changedata := make(map[string]interface{})

	if err == nil {
		//是老用户
		if num > 0 {
			logs.Info("old user")

			changedata["country"] = userInfo.Country
			changedata["sex"] = userInfo.Sex
			changedata["province"] = userInfo.Province
			changedata["city"] = userInfo.City
			self.loginSucesss(res[0], changedata, true)
			token = changedata["user_token"].(string)
			uid = res[0]["id"].(string)
		} else {
			logs.Info("new user")
			changedata["account"] = userInfo.OpenID
			changedata["name"] = userInfo.Nickname
			changedata["country"] = userInfo.Country
			changedata["sex"] = userInfo.Sex
			changedata["province"] = userInfo.Province
			changedata["city"] = userInfo.City
			changedata["head"] = userInfo.HeadImgURL
			changedata["wchat_unionid"] = userInfo.Unionid
			changedata["wchat_openid"] = userInfo.OpenID

			resinfo := self.addUser(changedata, "wchat_openid", userInfo.OpenID, true)
			uid = resinfo["id"].(string)
			token = resinfo["user_token"].(string)
		}
		newlocation := fmt.Sprintf("%s?uid=%s&token=%s", beego.AppConfig.String("wechat.loginokurl"), uid, token)
		logs.Info("goto new location:%s", newlocation)
		http.Redirect(self.Ctx.ResponseWriter, self.Ctx.Request, newlocation, 302)
	} else {
		self.AjaxReturnError(err.Error())
	}
}

// func (self *LoginController) CheckWchatLogin() {
// 	token := self.Input().Get("token")
// 	uid := self.Input().Get("uid")
// 	data := make(map[string]interface{})
// 	data["token"] = token
// 	data["uid"] = uid
// 	self.AjaxReturn(libs.WchatLoginOk, "登录成功", data)
// }
