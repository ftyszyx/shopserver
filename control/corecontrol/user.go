package corecontrol

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/coredata"
	"github.com/zyx/shop_server/models/names"
	"github.com/zyx/shop_server/wechat"
)

type UserController struct {
	base.BaseController
}

//检查数据正确性
func (self *UserController) checkData(data map[string]interface{}) error {
	groupModel := models.GetModel(names.USERGROUP)
	logs.Info("checkData")
	if value, ok := data["user_group"]; ok {

		groupid := value.(string)
		groupinfo := groupModel.GetInfoById(self.GetDb(), groupid)
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
	if self.GetMethod() == "Add" {
		err := self.checkData(self.GetPost())
		if err != nil {
			return err
		}
		defaultPass := beego.AppConfig.String("user.defaultPssword")
		logs.Info("default pass:%s", defaultPass)
		data["password"] = libs.GetStrMD5(defaultPass)
		data["reg_time"] = time.Now().Unix()
	} else if self.GetMethod() == "Edit" {
		return self.checkData(self.GetPost())
	} else if self.GetMethod() == "ChangeValid" {

		if self.CheckFieldExit(self.GetPost(), "is_del") == false {
			return errors.New("数据空")
		}
		data["is_del"] = self.GetPost()["is_del"]
	}
	return nil
}
func (self *UserController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.GetMethod() == "Add" {
		self.AddLog(fmt.Sprintf("增加角色:%+v", data))
	} else if self.GetMethod() == "Edit" {
		id := self.GetPost()["id"].(string)
		self.GetModel().ClearRowCache(id)
		self.AddLog(fmt.Sprintf("修改角色:%+v", data))
	} else if self.GetMethod() == "Del" {
		id := self.GetPost()["id"].(string)
		name := data["name"]
		self.AddLog(fmt.Sprintf("删除角色:%s", name))
		self.GetModel().ClearRowCache(id)
	} else if self.GetMethod() == "ChangeValid" {
		id := self.GetPost()["id"].(string)
		username := oldinfo["name"]
		valid := self.GetPost()["is_valid"]
		self.AddLog(fmt.Sprintf("修改角色有效:%s %+v", username, valid))
		self.GetModel().ClearRowCache(id)
	} else if self.GetMethod() == "ChangePassword" {
		username := oldinfo["name"]
		self.AddLog(fmt.Sprintf("修改角色密码:%s", username))
		self.GetModel().ClearRowCache(self.GetUid())
	} else if self.GetMethod() == "UpdateCart" {
		self.AddLog(fmt.Sprintf("%+v", data))
		self.GetModel().ClearRowCache(self.GetUid())
	} else if self.GetMethod() == "UpdateAddress" {
		self.AddLog(fmt.Sprintf("%+v", data))
		self.GetModel().ClearRowCache(self.GetUid())
	} else if self.GetMethod() == "RefreshToken" {
		self.AddLog(fmt.Sprintf("%+v", data))
		senddata := make(map[string]interface{})
		senddata["user_token"] = data["user_token"]
		senddata["token_expire"] = data["token_expire"]
		self.AjaxReturnSuccess("", senddata)
	} else if self.GetMethod() == "ResetPassword" {
		self.AddLog(fmt.Sprintf("重置密码：%s", oldinfo["name"]))
		self.GetModel().ClearRowCache(self.GetUid())
	} else {
		self.AddLog(fmt.Sprintf("%+v", data))
	}
	return nil
}

func (self *UserController) ChangeValid() {
	self.EditCommonAndReturn(self)
}

func (self *UserController) ChangePassword() {

	self.CheckFieldExitAndReturn(self.GetPost(), "password", "密码不能为空")
	changedata := make(map[string]interface{})
	changedata["password"] = libs.GetStrMD5(self.GetPost()["password"].(string))
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetUid())
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
	self.CheckFieldExitAndReturn(self.GetPost(), "name", "姓名不能为空")
	changedata := make(map[string]interface{})
	changedata["name"] = self.GetPost()["name"]
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetUid())
}

func (self *UserController) UpdateAccount() {
	self.CheckFieldExitAndReturn(self.GetPost(), "account", "账号不能为空")
	if self.GetModel().CheckExit(self.GetDb(), "account", self.GetPost()["account"]) == true {
		self.AjaxReturnError(errors.New("账号名已存在,修改失败"))
	}
	changedata := make(map[string]interface{})
	changedata["account"] = self.GetPost()["account"]
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetUid())
}

func (self *UserController) UpdateHead() {
	self.CheckFieldExitAndReturn(self.GetPost(), "head", "头像不能为空")
	changedata := make(map[string]interface{})
	changedata["head"] = self.GetPost()["head"]
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetUid())
}

func (self *UserController) UpdatePhone() {
	self.CheckFieldExitAndReturn(self.GetPost(), "phone", "手机号不能为空")
	self.CheckFieldExitAndReturn(self.GetPost(), "code", "验证码不能为空")
	phone := self.GetPost()["phone"].(string)
	code := self.GetPost()["code"].(string)

	codestr, ok := models.PhoneCodeCache.Get(phone).(string)
	if ok == false || codestr == "" || codestr != code {
		self.AjaxReturnError(errors.New("验证码不对"))
	}
	changedata := make(map[string]interface{})
	changedata["phone"] = phone
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetUid())

}

func (self *UserController) ResetPassword() {

	self.CheckFieldExitAndReturn(self.GetPost(), "id", "id不能为空")

	changedata := make(map[string]interface{})
	// password := string(utils.RandomCreateBytes(6))
	password := beego.AppConfig.String("user.defaultPssword")
	changedata["password"] = libs.GetStrMD5(password)

	//self.UpdateSqlByIdAndReturn(self.GetDb(),self,changedata, self.GetPost()["id"])

	err := self.UpdateSqlCommon(self, changedata, "id", self.GetPost()["id"])
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}

	self.AjaxReturnSuccess("", map[string]interface{}{"newpass": password})
}

//添加购物车
func (self *UserController) UpdateCart() {
	if self.GetUid() == "" {
		self.AjaxReturn(libs.AuthFail, "请先登录", nil)
	}
	changedata := make(map[string]interface{})
	changedata["shop_cart"] = self.GetPost()["shop_cart"]
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetUid())
}

func (self *UserController) UpdateAddress() {
	if self.GetUid() == "" {
		self.AjaxReturn(libs.AuthFail, "请先登录", nil)
	}
	changedata := make(map[string]interface{})
	changedata["address"] = self.GetPost()["address"]
	self.UpdateSqlByIdAndReturn(self, changedata, self.GetUid())
}

//获取角色信息
func (self *UserController) GetUserInfo() {
	usermodel := models.GetModel(names.USER)
	groupmodel := models.GetModel(names.USERGROUP)
	modulemodel := models.GetModel(names.MODULE)
	id := self.Ctx.Request.Header.Get("uid")
	logs.Info("get user info id:%v", id)
	if id == "" {
		self.AjaxReturn(libs.AuthFail, "uid空", nil)
	}
	userInfo := usermodel.GetInfoAndCache(self.GetDb(), id, false) //更新缓存
	if userInfo == nil {
		self.AjaxReturn(libs.ErrorCode, "无效", nil)
	}
	var data = make(map[string]interface{})
	groupinfo := groupmodel.GetInfoAndCache(self.GetDb(), userInfo["user_group"].(string), false)
	if groupinfo == nil {
		self.AjaxReturn(libs.ErrorCode, "用户组不存在", nil)
	}
	data["groupinfo"] = groupinfo
	data["head"] = userInfo["head"]

	grouptype := groupinfo["group_type"].(string)
	data["limit_show_order"] = groupinfo["limit_show_order"]
	if grouptype == strconv.Itoa(libs.UserAdmin) {
		var modules []interface{}
		modulelist := modulemodel.Cache().Get("allmodel")
		if modulelist != nil {
			modulelistarr := modulelist.([]db.Params)
			for _, vlaue := range modulelistarr {
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
		}

		data["modules"] = modules
	}

	data["name"] = userInfo["name"]
	expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	usermodel.(*coredata.User).ExtendExpireTime(self.GetDb(), id, expiretime) //延长时间
	self.AjaxReturn(libs.SuccessCode, nil, data)
}

//刷新token
func (self *UserController) RefreshToken() {
	userGroupModel := models.GetModel(names.USERGROUP)
	self.CheckFieldExitAndReturn(self.GetPost(), "id", "操作对象空")
	changedata := make(map[string]interface{})
	uid := self.GetPost()["id"].(string)
	userinfo := self.GetModel().GetInfoAndCache(self.GetDb(), uid, false)
	if userinfo == nil {
		self.AjaxReturnError(errors.New("角色不存在"))
	}
	curtime := time.Now().Unix()
	groupid := userinfo["user_group"].(string)
	groupinfo := userGroupModel.GetInfoAndCache(self.GetDb(), groupid, false)
	expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	usertoken := libs.GetToken(curtime, uid, userinfo["password"], groupid)
	changedata["user_token"] = usertoken
	changedata["token_get_time"] = curtime
	changedata["token_expire"] = curtime + int64(expiretime)
	changedata["last_login_time"] = curtime
	self.GetModel().ClearRowCache(uid)
	self.UpdateSqlByIdAndReturn(self, changedata, uid)

}

//商城获取用户信息
func (self *UserController) GetShopUserInfo() {
	usermodel := models.GetModel(names.USER)
	ordermodel := models.GetModel(names.SHOP_ORDER)
	userGroupModel := models.GetModel(names.USERGROUP)
	id := self.Ctx.Request.Header.Get("uid")
	userInfo := usermodel.GetInfoAndCache(self.GetDb(), id, false) //更新缓存

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
		data["order_waitpay"] = ordermodel.GetNumByField(self.GetDb(), map[string]interface{}{"user_id": self.GetUid(), "status": libs.OrderStatusWaitPay})
		data["order_pay"] = ordermodel.GetNumByField(self.GetDb(), map[string]interface{}{"user_id": self.GetUid(), "status": libs.OrderStatusWaitcheck})
		data["order_send"] = ordermodel.GetNumByField(self.GetDb(), map[string]interface{}{"user_id": self.GetUid(), "status": libs.OrderStatusSend})
		data["order_refund"] = ordermodel.GetNumByField(self.GetDb(), map[string]interface{}{"user_id": self.GetUid(), "status": libs.OrderStatusRefund})
	}
	update := self.Input().Get("update")
	if update != "" {
		groupinfo := userGroupModel.GetInfoAndCache(self.GetDb(), userInfo["user_group"].(string), false)
		expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
		if err != nil {
			self.AjaxReturnError(errors.WithStack(err))
		}
		usermodel.(*coredata.User).ExtendExpireTime(self.GetDb(), id, expiretime) //延长时间
	}
	self.AjaxReturn(libs.SuccessCode, nil, data)
}

func (self *UserController) GetWchatJsConf() {
	url := self.GetPost()["url"]
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
