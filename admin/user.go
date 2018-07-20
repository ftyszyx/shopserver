package admin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
)

type UserController struct {
	BaseController
}

//检查数据正确性
func (self *UserController) checkData(data map[string]interface{}) {
	groupModel := models.GetModel(models.USERGROUP)
	logs.Info("checkData")
	if value, ok := data["user_group"]; ok {

		groupid := value.(string)
		groupinfo := groupModel.GetInfoById(groupid)
		if groupinfo == nil {
			self.AjaxReturnError("用户组不对")
		}
		grouptype := groupinfo["group_type"].(string)
		if grouptype == strconv.Itoa(libs.UserSystem) {
			self.AjaxReturnError("系统用户组不可设")
		}
	}
}

func (self *UserController) BeforeSql(data map[string]interface{}) {
	if self.method == "Add" {
		self.checkData(self.postdata)
		defaultPass := beego.AppConfig.String("user.defaultPssword")
		logs.Info("default pass:%s", defaultPass)
		data["password"] = libs.GetStrMD5(defaultPass)
		data["reg_time"] = time.Now().Unix()
	} else if self.method == "Edit" {
		self.checkData(self.postdata)
	} else if self.method == "ChangeValid" {
		self.CheckFieldExit(self.postdata, "is_del", "数据空")
		data["is_del"] = self.postdata["is_del"]
	} else if self.method == "ChangePassword" {
		self.CheckFieldExit(self.postdata, "password", "密码为空")
		data["password"] = libs.GetStrMD5(self.postdata["password"].(string))
	}
}
func (self *UserController) AfterSql(data map[string]interface{}, oldinfo orm.Params) {
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
		oldinfo := self.model.GetInfoById(self.uid)
		username := oldinfo["name"]
		valid := self.postdata["is_valid"]
		self.AddLog(fmt.Sprintf("修改角色有效:%s %+v", username, valid))
		self.model.ClearRowCache(id)
	} else if self.method == "ChangePassword" {
		oldinfo := self.model.GetInfoById(self.uid)
		username := oldinfo["name"]
		self.AddLog(fmt.Sprintf("修改角色密码:%s", username))

		self.model.ClearRowCache(self.uid)
	} else if self.method == "UpdateCart" {
		oldinfo := self.model.GetInfoById(self.uid)
		username := oldinfo["name"]
		self.AddLog(fmt.Sprintf("修改角色购物车:%s", username))

		self.model.ClearRowCache(self.uid)
	} else if self.method == "RefreshToken" {
		self.AddLog(fmt.Sprintf("%+v", data))
		senddata := make(map[string]interface{})
		senddata["user_token"] = data["user_token"]
		senddata["token_expire"] = data["token_expire"]
		self.AjaxReturnSuccess("", senddata)

	} else {
		self.AddLog(fmt.Sprintf("%+v", data))
	}
}

func (self *UserController) ChangeValid() {
	self.EditCommon(self)
}

func (self *UserController) ChangePassword() {
	self.postdata["id"] = self.uid
	self.EditCommon(self)
}

func (self *UserController) Add() {
	self.AddCommon(self)
}

func (self *UserController) Edit() {
	self.EditCommon(self)
}

func (self *UserController) Del() {
	self.DelCommon(self)
}

func (self *UserController) UpdateName() {
	self.CheckFieldExit(self.postdata, "name", "姓名不能为空")
	changedata := make(map[string]interface{})
	changedata["name"] = self.postdata["name"]
	self.updateSqlById(self, changedata, self.uid)
}

func (self *UserController) UpdateHead() {
	self.CheckFieldExit(self.postdata, "head", "头像不能为空")
	changedata := make(map[string]interface{})
	changedata["head"] = self.postdata["head"]
	self.updateSqlById(self, changedata, self.uid)
}

func (self *UserController) UpdatePhone() {
	self.CheckFieldExit(self.postdata, "phone", "手机号不能为空")
	self.CheckFieldExit(self.postdata, "code", "验证码不能为空")
	phone := self.postdata["phone"].(string)
	code := self.postdata["code"].(string)

	codestr, ok := models.PhoneCodeCache.Get(phone).(string)
	if ok == false || codestr == "" || codestr != code {
		self.AjaxReturnError("验证码不对")
	}
	changedata := make(map[string]interface{})
	changedata["phone"] = phone
	self.updateSqlById(self, changedata, self.uid)

}

//添加购物车
func (self *UserController) UpdateCart() {
	if self.uid == "" {
		self.AjaxReturn(libs.AuthFail, "请先登录", nil)
	}
	changedata := make(map[string]interface{})
	changedata["shop_cart"] = self.postdata["shop_cart"]
	self.updateSqlById(self, changedata, self.uid)
}

func (self *UserController) UpdateAddress() {
	if self.uid == "" {
		self.AjaxReturn(libs.AuthFail, "请先登录", nil)
	}
	changedata := make(map[string]interface{})
	changedata["address"] = self.postdata["address"]
	self.updateSqlById(self, changedata, self.uid)
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
	userInfo := usermodel.GetInfoAndCache(id, false) //更新缓存
	if userInfo == nil {
		self.AjaxReturn(libs.ErrorCode, "无效", nil)
	}
	var data = make(map[string]interface{})
	groupinfo := groupmodel.GetInfoAndCache(userInfo["user_group"].(string), false)
	data["groupinfo"] = groupinfo
	data["head"] = userInfo["head"]

	if groupinfo["group_type"].(string) == strconv.Itoa(libs.UserAdmin) {
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
		self.AjaxReturnError(err.Error())
	}
	usermodel.(*models.User).ExtendExpireTime(id, expiretime) //延长时间
	self.AjaxReturn(libs.SuccessCode, nil, data)
}

//刷新token
func (self *UserController) RefreshToken() {
	userGroupModel := models.GetModel(models.USERGROUP)
	self.CheckFieldExit(self.postdata, "id", "操作对象空")
	changedata := make(map[string]interface{})
	uid := self.postdata["id"].(string)
	userinfo := self.model.GetInfoAndCache(uid, false)
	if userinfo == nil {
		self.AjaxReturnError("角色不存在")
	}
	curtime := time.Now().Unix()
	groupid := userinfo["user_group"].(string)
	groupinfo := userGroupModel.GetInfoAndCache(groupid, false)
	expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	usertoken := libs.GetToken(curtime, uid, userinfo["password"], groupid)
	changedata["user_token"] = usertoken
	changedata["token_get_time"] = curtime
	changedata["token_expire"] = curtime + int64(expiretime)
	changedata["last_login_time"] = curtime
	self.model.ClearRowCache(uid)
	self.updateSqlById(self, changedata, uid)

}

//商城获取用户信息
func (self *UserController) GetShopUserInfo() {
	usermodel := models.GetModel(models.USER)
	ordermodel := models.GetModel(models.SHOP_ORDER)
	userGroupModel := models.GetModel(models.USERGROUP)
	id := self.Ctx.Request.Header.Get("uid")
	userInfo := usermodel.GetInfoAndCache(id, false) //更新缓存

	var data = make(map[string]interface{})
	data["name"] = userInfo["name"]
	data["account"] = userInfo["account"]
	data["mail"] = userInfo["mail"]
	data["phone"] = userInfo["phone"]
	data["shop_cart"] = userInfo["shop_cart"]
	data["head"] = userInfo["head"]
	data["address"] = userInfo["address"]
	data["wchat_openid"] = userInfo["wchat_openid"]
	order := self.Input().Get("order")
	if order != "" {
		//要获取订单信息
		data["order_waitpay"] = ordermodel.GetNumByField(map[string]interface{}{"user_id": self.uid, "status": libs.OrderStatusWaitPay})
		data["order_pay"] = ordermodel.GetNumByField(map[string]interface{}{"user_id": self.uid, "status": libs.OrderStatusWaitSend})
		data["order_send"] = ordermodel.GetNumByField(map[string]interface{}{"user_id": self.uid, "status": libs.OrderStatusSend})
		data["order_refund"] = ordermodel.GetNumByField(map[string]interface{}{"user_id": self.uid, "status": libs.OrderStatusRefund})
	}
	update := self.Input().Get("update")
	if update != "" {

		groupinfo := userGroupModel.GetInfoAndCache(userInfo["user_group"].(string), false)
		expiretime, err := strconv.Atoi(groupinfo["expire_time"].(string))
		if err != nil {
			self.AjaxReturnError(err.Error())
		}
		usermodel.(*models.User).ExtendExpireTime(id, expiretime) //延长时间
	}
	self.AjaxReturn(libs.SuccessCode, nil, data)
}
