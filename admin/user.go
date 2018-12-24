package admin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/wechat"
)

type UserController struct {
	BaseController
}

//检查数据正确性
func (self *UserController) checkData(data map[string]interface{}) error {
	groupModel := models.GetModel(models.USERGROUP)
	logs.Info("checkData")
	if value, ok := data["user_group"]; ok {

		groupid := value.(string)
		groupinfo := groupModel.GetInfoById(self.dboper, groupid)
		if groupinfo == nil {

			return errors.New("用户组不对")
		}
		grouptype := groupinfo["group_type"].(string)
		if grouptype == strconv.Itoa(libs.UserSystem) {

			return errors.New("系统用户组不可设")
		}
	}
	return nil
}

func (self *UserController) BeforeSql(data map[string]interface{}) error {
	if self.method == "Add" {
		err := self.checkData(self.postdata)
		if err != nil {
			return err
		}
		defaultPass := beego.AppConfig.String("user.defaultPssword")
		logs.Info("default pass:%s", defaultPass)
		data["password"] = libs.GetStrMD5(defaultPass)
		data["reg_time"] = time.Now().Unix()
	} else if self.method == "Edit" {
		return self.checkData(self.postdata)
	} else if self.method == "ChangeValid" {

		if self.CheckFieldExit(self.postdata, "is_del") == false {
			return errors.New("数据空")
		}
		data["is_del"] = self.postdata["is_del"]
	}
	return nil
}
func (self *UserController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.method == "Add" {
		self.AddLog(fmt.Sprintf("增加角色:%+v", data))
	} else if self.method == "Edit" {
		id := self.postdata["id"].(string)
		self.model.ClearRowCache(id)
		self.AddLog(fmt.Sprintf("修改角色:%+v", data))
	} else if self.method == "Del" {
		id := self.postdata["id"].(string)
		name := data["name"]
		self.AddLog(fmt.Sprintf("删除角色:%s", name))
		self.model.ClearRowCache(id)
	} else if self.method == "ChangeValid" {
		id := self.postdata["id"].(string)
		username := oldinfo["name"]
		valid := self.postdata["is_valid"]
		self.AddLog(fmt.Sprintf("修改角色有效:%s %+v", username, valid))
		self.model.ClearRowCache(id)
	} else if self.method == "ChangePassword" {
		username := oldinfo["name"]
		self.AddLog(fmt.Sprintf("修改角色密码:%s", username))
		self.model.ClearRowCache(self.uid)
	} else if self.method == "UpdateCart" {
		self.AddLog(fmt.Sprintf("%+v", data))
		self.model.ClearRowCache(self.uid)
	} else if self.method == "UpdateAddress" {
		self.AddLog(fmt.Sprintf("%+v", data))
		self.model.ClearRowCache(self.uid)
	} else if self.method == "RefreshToken" {
		self.AddLog(fmt.Sprintf("%+v", data))
		senddata := make(map[string]interface{})
		senddata["user_token"] = data["user_token"]
		senddata["token_expire"] = data["token_expire"]
		self.AjaxReturnSuccess("", senddata)
	} else if self.method == "ResetPassword" {
		self.AddLog(fmt.Sprintf("重置密码：%s", oldinfo["name"]))
		self.model.ClearRowCache(self.uid)
	} else {
		self.AddLog(fmt.Sprintf("%+v", data))
	}
	return nil
}

func (self *UserController) ChangeValid() {
	self.EditCommonAndReturn(self)
}

func (self *UserController) ChangePassword() {

	self.CheckFieldExitAndReturn(self.postdata, "password", "密码不能为空")
	changedata := make(map[string]interface{})
	changedata["password"] = libs.GetStrMD5(self.postdata["password"].(string))
	self.updateSqlByIdAndReturn(self, changedata, self.uid)
}

func (self *UserController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *UserController) Edit() {

	self.EditCommonAndReturn(self)
}

func (self *UserController) Del() {
	self.AjaxReturnError(errors.New("不能删除用户"))
	self.DelCommonAndReturn(self)
}

func (self *UserController) UpdateName() {
	self.CheckFieldExitAndReturn(self.postdata, "name", "姓名不能为空")
	changedata := make(map[string]interface{})
	changedata["name"] = self.postdata["name"]
	self.updateSqlByIdAndReturn(self, changedata, self.uid)
}

func (self *UserController) UpdateAccount() {
	self.CheckFieldExitAndReturn(self.postdata, "account", "账号不能为空")
	if self.model.CheckExit(self.dboper, "account", self.postdata["account"]) == true {
		self.AjaxReturnError(errors.New("账号名已存在,修改失败"))
	}
	changedata := make(map[string]interface{})
	changedata["account"] = self.postdata["account"]
	self.updateSqlByIdAndReturn(self, changedata, self.uid)
}

func (self *UserController) UpdateHead() {
	self.CheckFieldExitAndReturn(self.postdata, "head", "头像不能为空")
	changedata := make(map[string]interface{})
	changedata["head"] = self.postdata["head"]
	self.updateSqlByIdAndReturn(self, changedata, self.uid)
}

func (self *UserController) UpdatePhone() {
	self.CheckFieldExitAndReturn(self.postdata, "phone", "手机号不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "code", "验证码不能为空")
	phone := self.postdata["phone"].(string)
	code := self.postdata["code"].(string)

	codestr, ok := models.PhoneCodeCache.Get(phone).(string)
	if ok == false || codestr == "" || codestr != code {
		self.AjaxReturnError(errors.New("验证码不对"))
	}
	changedata := make(map[string]interface{})
	changedata["phone"] = phone
	self.updateSqlByIdAndReturn(self, changedata, self.uid)

}

func (self *UserController) ResetPassword() {

	self.CheckFieldExitAndReturn(self.postdata, "id", "id不能为空")

	changedata := make(map[string]interface{})
	// password := string(utils.RandomCreateBytes(6))
	password := beego.AppConfig.String("user.defaultPssword")
	changedata["password"] = libs.GetStrMD5(password)

	//self.updateSqlByIdAndReturn(self.dboper,self,changedata, self.postdata["id"])

	err := self.updateSqlCommon(self, changedata, "id", self.postdata["id"])
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}

	self.AjaxReturnSuccess("", map[string]interface{}{"newpass": password})
}

//添加购物车
func (self *UserController) UpdateCart() {
	if self.uid == "" {
		self.AjaxReturn(libs.AuthFail, "请先登录", nil)
	}
	changedata := make(map[string]interface{})
	changedata["shop_cart"] = self.postdata["shop_cart"]
	self.updateSqlByIdAndReturn(self, changedata, self.uid)
}

func (self *UserController) UpdateAddress() {
	if self.uid == "" {
		self.AjaxReturn(libs.AuthFail, "请先登录", nil)
	}
	changedata := make(map[string]interface{})
	changedata["address"] = self.postdata["address"]
	self.updateSqlByIdAndReturn(self, changedata, self.uid)
}

//获取角色信息
func (self *UserController) GetUserInfo() {
	usermodel := models.GetModel(models.USER)
	groupmodel := models.GetModel(models.USERGROUP)
	modulemodel := models.GetModel(models.MODULE)
	id := self.Ctx.Request.Header.Get("uid")
	logs.Info("get user info id:%v", id)
	if id == "" {
		self.AjaxReturn(libs.AuthFail, "uid空", nil)
	}
	userInfo := usermodel.GetInfoAndCache(self.dboper, id, false) //更新缓存
	if userInfo == nil {
		self.AjaxReturn(libs.ErrorCode, "无效", nil)
	}
	var data = make(map[string]interface{})
	groupinfo := groupmodel.GetInfoAndCache(self.dboper, userInfo["user_group"].(string), false)
	if groupinfo == nil {
		self.AjaxReturn(libs.ErrorCode, "用户组不存在", nil)
	}
	data["groupinfo"] = groupinfo
	data["head"] = userInfo["head"]

	grouptype := groupinfo["group_type"].(string)
	data["limit_show_order"] = groupinfo["limit_show_order"]
	if grouptype == strconv.Itoa(libs.UserAdmin) {
		var modules []interface{}
		modulelist := modulemodel.Cache()
		for _, vlaue := range modulelist {
			//logs.Info("value:%v", vlaue)
			if vlaue["need_auth"].(string) == "1" {
				//logs.Info("add")
				var item = make(map[string]interface{})
				item["id"] = vlaue["id"]
				item["controller"] = vlaue["controller"]
				item["method"] = vlaue["method"]
				modules = append(modules, item)
			}
		}
		data["modules"] = modules
	}

	data["name"] = userInfo["name"]
	expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	usermodel.(*models.User).ExtendExpireTime(self.dboper, id, expiretime) //延长时间
	self.AjaxReturn(libs.SuccessCode, nil, data)
}

//刷新token
func (self *UserController) RefreshToken() {
	userGroupModel := models.GetModel(models.USERGROUP)
	self.CheckFieldExitAndReturn(self.postdata, "id", "操作对象空")
	changedata := make(map[string]interface{})
	uid := self.postdata["id"].(string)
	userinfo := self.model.GetInfoAndCache(self.dboper, uid, false)
	if userinfo == nil {
		self.AjaxReturnError(errors.New("角色不存在"))
	}
	curtime := time.Now().Unix()
	groupid := userinfo["user_group"].(string)
	groupinfo := userGroupModel.GetInfoAndCache(self.dboper, groupid, false)
	expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	usertoken := libs.GetToken(curtime, uid, userinfo["password"], groupid)
	changedata["user_token"] = usertoken
	changedata["token_get_time"] = curtime
	changedata["token_expire"] = curtime + int64(expiretime)
	changedata["last_login_time"] = curtime
	self.model.ClearRowCache(uid)
	self.updateSqlByIdAndReturn(self, changedata, uid)

}

//商城获取用户信息
func (self *UserController) GetShopUserInfo() {
	usermodel := models.GetModel(models.USER)
	ordermodel := models.GetModel(models.SHOP_ORDER)
	userGroupModel := models.GetModel(models.USERGROUP)
	id := self.Ctx.Request.Header.Get("uid")
	userInfo := usermodel.GetInfoAndCache(self.dboper, id, false) //更新缓存

	var data = make(map[string]interface{})
	data["name"] = userInfo["name"]
	data["account"] = userInfo["account"]
	data["mail"] = userInfo["mail"]
	data["phone"] = userInfo["phone"]
	data["shop_cart"] = userInfo["shop_cart"]
	data["head"] = userInfo["head"]
	data["groupid"] = userInfo["user_group"]
	data["address"] = userInfo["address"]
	data["groupid"] = userInfo["user_group"]
	data["wchat_openid"] = userInfo["wchat_openid"]
	order := self.Input().Get("order")
	if order != "" {
		//要获取订单信息
		data["order_waitpay"] = ordermodel.GetNumByField(self.dboper, map[string]interface{}{"user_id": self.uid, "status": libs.OrderStatusWaitPay})
		data["order_pay"] = ordermodel.GetNumByField(self.dboper, map[string]interface{}{"user_id": self.uid, "status": libs.OrderStatusWaitcheck})
		data["order_send"] = ordermodel.GetNumByField(self.dboper, map[string]interface{}{"user_id": self.uid, "status": libs.OrderStatusSend})
		data["order_refund"] = ordermodel.GetNumByField(self.dboper, map[string]interface{}{"user_id": self.uid, "status": libs.OrderStatusRefund})
	}
	update := self.Input().Get("update")
	if update != "" {
		groupinfo := userGroupModel.GetInfoAndCache(self.dboper, userInfo["user_group"].(string), false)
		expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
		if err != nil {
			self.AjaxReturnError(errors.WithStack(err))
		}
		usermodel.(*models.User).ExtendExpireTime(self.dboper, id, expiretime) //延长时间
	}
	self.AjaxReturn(libs.SuccessCode, nil, data)
}

func (self *UserController) GetWchatJsConf() {
	url := self.postdata["url"]
	if url == nil {
		self.AjaxReturnError(errors.New("参数错误"))
	}
	config, err := wechat.JsdkInstance.GetConfig(url.(string))
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	senddata := make(map[string]interface{})
	senddata["appid"] = config.AppID
	senddata["timestamp"] = config.Timestamp
	senddata["nonceStr"] = config.NonceStr
	senddata["signature"] = config.Signature
	self.AjaxReturnSuccess("", config)
}
