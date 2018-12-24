package admin

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type ShopOrderController struct {
	BaseController
}

func (self *ShopOrderController) Edit() {
	self.EditCommonAndReturn(self)
}

type orderinfo struct {
	lock     sync.Mutex
	lasttime int64
	Num      int
}

type RrpResultData struct {
	Id        string
	Logistics []string
}

type RrpResult struct {
	Code    string
	Message string
	Data    []RrpResultData
}

var ORDER_PRE = "WO"

func (self *ShopOrderController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.method == "UpdateIdNum" {
		oldroleinfo := make(map[string]interface{})
		oldroleinfo["idnum"] = oldinfo["idnum"]
		oldroleinfo["idnumpic1"] = oldinfo["idnumpic1"]
		oldroleinfo["idnumpic2"] = oldinfo["idnumpic2"]
		self.AddLog(fmt.Sprintf("change info:%+v oldroleinfo:%+v", data, oldroleinfo))
	} else {
		self.AddLog(fmt.Sprintf("change info:%+v ", self.postdata))
	}

	return nil

}

//导出物流单号
func (self *ShopOrderController) ExportToErp() {
	idarr := self.postdata["ids"].([]interface{})
	sendflag := self.postdata["sendflag"].(float64)

	var dataarr []interface{}
	for _, id := range idarr {
		dataarr = append(dataarr, self.getErpExportData(id.(string), nil, func(status int) bool {
			if status == libs.OrderStatusWaitSend || status == libs.OrderStatusSend {
				return true
			}
			return false
		}, ""))
	}
	self.exportToErpCommon(dataarr, sendflag, nil)

}

func (self *ShopOrderController) exportToErpCommon(dataarr []interface{}, sendflag float64, changevalue map[string]interface{}) {
	urlstr := beego.AppConfig.String("erp.url") + "Sell/addOneOrder"
	token := beego.AppConfig.String("erp.shoptoken")
	req := httplib.Post(urlstr)
	senddata := make(map[string]interface{})
	senddata["data"] = dataarr
	senddata["token"] = token
	senddata["sendflag"] = sendflag
	senddata["shop_id"] = beego.AppConfig.String("erp.shopid")
	reqbuf, err := json.Marshal(senddata)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}

	req.Body(string(reqbuf))
	req.Header("Content-Type", "application/json")
	logs.Info("send request:")
	respdata, err := req.Bytes()
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	getData := new(RrpResult)
	//logs.Info("get data:%s", string(respdata))
	err = json.Unmarshal(respdata, getData)
	if err != nil {
		logs.Info("parse data err")
		self.AjaxReturnError(errors.WithStack(err))
	}
	logs.Info("get data:%v", getData)
	if getData.Code != "1" {
		self.AjaxReturnError(errors.New(getData.Message))
	}
	// db := orm.NewOrm()
	self.dboper.Begin()
	// logicmodel := models.GetModel(models.LOGICSTICS).(*models.Logistics)
	for _, dataitem := range getData.Data {
		//读每一行
		changedata := make(map[string]interface{})

		if changevalue != nil {
			changedata["pay_time"] = changevalue["pay_time"]
			changedata["pay_check_info"] = changevalue["pay_check_info"]
			changedata["pay_type"] = changevalue["pay_type"]
		}
		logisticstr, err := json.Marshal(dataitem.Logistics)
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}

		if sendflag == 1 {
			//新增物流
			changedata["shipment_num"] = string(logisticstr)
			changedata["status"] = libs.OrderStatusSend
		} else {
			changedata["status"] = libs.OrderStatusWaitSend

		}

		err = self.updateSqlById(self, changedata, dataitem.Id)
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}
		if sendflag == 1 {
			self.AddLog(fmt.Sprintf("订单：%s 增加物流单号 %+v", dataitem.Id, dataitem.Logistics))
		} else {
			self.AddLog(fmt.Sprintf("订单：%s 导出erp", dataitem.Id))
		}
	}
	self.dboper.Commit()

	self.AjaxReturnSuccessNull()
}

//管理员审核支付
func (self *ShopOrderController) CheckPayOk() {
	idarr := self.postdata["ids"].([]interface{})
	sendflag := self.postdata["sendflag"].(float64)
	self.CheckFieldExitAndReturn(self.postdata, "pay_check_info", "审核信息不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "pay_type", "支付类型不能为空")
	var dataarr []interface{}
	paytime := time.Now().Unix()
	pay_check_info := self.postdata["pay_check_info"].(string)

	for _, id := range idarr {
		dataarr = append(dataarr, self.getErpExportData(id.(string), paytime, func(status int) bool {
			if status == libs.OrderStatusWaitcheck {
				return true
			}
			return false
		}, pay_check_info))
	}
	self.exportToErpCommon(dataarr, sendflag, map[string]interface{}{"pay_check_info": pay_check_info, "pay_time": paytime, "pay_type": self.postdata["pay_type"]})

}

type ErpSelldataList []models.ErpSelldata
type EditShipNumDataArr struct {
	List map[string]ErpSelldataList
}

//修改物流单号
func (self *ShopOrderController) UpdateOrderShipNum() {
	getData := new(EditShipNumDataArr)
	err := json.Unmarshal(self.Ctx.Input.RequestBody, getData)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	// db := orm.NewOrm()
	self.dboper.Begin()
	for orderid, dataitem := range getData.List {
		//读每一行

		orderinfo := self.model.GetInfoById(self.dboper, orderid)
		if orderinfo != nil {
			changedata := make(map[string]interface{})
			changedata["status"] = libs.OrderStatusSend
			var shiplist []string
			for _, sellitem := range dataitem {
				shiplist = append(shiplist, sellitem.Logistics)
			}
			shipliststr, err := json.Marshal(shiplist)
			if err != nil {
				self.AjaxReturnError(errors.WithStack(err))
			}
			changedata["shipment_num"] = string(shipliststr)
			err = self.updateSqlById(self, changedata, orderid)
			if err != nil {
				self.dboper.Rollback()
				self.AjaxReturnError(errors.WithStack(err))
			}
		}

	}
	self.dboper.Commit()
	self.AjaxReturnSuccess("成功", nil)
}

//导出到erp
func (self *ShopOrderController) getErpExportData(orderid string, paytime interface{}, checkstatus func(status int) bool, checkinfo string) map[string]interface{} {
	usermodel := models.GetModel(models.USER)
	orderinfo := self.model.GetInfoById(self.dboper, orderid)
	if orderinfo == nil {
		self.AjaxReturnError(errors.New("订单不存在"))
	}
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if checkstatus(status) == false {
		self.AjaxReturnError(errors.New("订单状态不对"))
	}
	sendata := make(map[string]interface{})
	sendata["shop_order"] = orderinfo["id"]
	if paytime == nil {
		sendata["pay_time"] = orderinfo["pay_time"]
	} else {
		sendata["pay_time"] = paytime
	}
	if checkinfo == "" {
		sendata["pay_check_info"] = orderinfo["pay_check_info"]
	} else {
		sendata["pay_check_info"] = checkinfo
	}
	userinfo := usermodel.GetInfoAndCache(self.dboper, orderinfo["user_id"].(string), false)
	sendata["customer_account"] = userinfo["account"]
	sendata["pay_id"] = orderinfo["pay_id"]
	sendata["supply_source"] = orderinfo["supply_source"]

	sendata["customer_username"] = userinfo["name"]
	sendata["customer_userid"] = userinfo["id"]

	//跟单员
	trackuserid := userinfo["track_admin"].(string)
	if trackuserid != "0" && trackuserid != "" {
		trackuserinfo := usermodel.GetInfoAndCache(self.dboper, trackuserid, false)
		sendata["track_man"] = trackuserinfo["name"]
	} else {
		sendata["track_man"] = ""
	}

	sendata["customer_name"] = orderinfo["client_name"]
	sendata["customer_addr"] = orderinfo["client_address"]
	sendata["customer_province"] = orderinfo["client_provice"]
	sendata["customer_city"] = orderinfo["client_city"]
	sendata["customer_area"] = orderinfo["client_area"]
	sendata["user_id_number"] = orderinfo["idnum"]
	sendata["client_phone"] = orderinfo["client_phone"]
	sendata["user_info"] = orderinfo["client_info"]
	sendata["order_time"] = orderinfo["order_time"]
	sendata["idnumpic1"] = orderinfo["idnumpic1"]
	sendata["idnumpic2"] = orderinfo["idnumpic2"]

	sendata["sell_vip_type"] = orderinfo["order_vip_type"]
	sendata["send_user_name"] = orderinfo["send_user_name"]
	sendata["freight_price"] = orderinfo["freight_price"]
	sendata["service_price"] = orderinfo["service_price"]
	sendata["send_user_phone"] = orderinfo["send_user_phone"]
	sendata["sell_type"] = 0

	sendata["item_code"] = orderinfo["itemcode"]
	sendata["num"] = orderinfo["num"]
	sendata["total_price"] = orderinfo["total_price"]
	sendata["pay_type"] = orderinfo["pay_type"]

	return sendata
}

func updateOrderTime(timecount int64) {
	orderidinfo.lock.Lock()
	if orderidinfo.lasttime == timecount {
		orderidinfo.Num++
	} else {
		orderidinfo.Num = 1
	}
	orderidinfo.lasttime = timecount
	orderidinfo.lock.Unlock()
}

var orderidinfo orderinfo

func getOrderid() (string, error) {
	nowtime := time.Now()
	updateOrderTime(nowtime.Unix())
	orderid := fmt.Sprintf("%s%d%d", ORDER_PRE, nowtime.Unix(), orderidinfo.Num)
	return orderid, nil
}

func (self *ShopOrderController) OrdersUpload() {
	err := self.UploadeCSV(self)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	} else {
		self.AjaxReturnSuccessNull()
	}
}

func (self *ShopOrderController) AddOneRow(rownum int, rowinfo []string) string {
	itemmodel := models.GetModel(models.SHOP_ITEM).(*models.ShopItem)
	usermodel := models.GetModel(models.USER)
	paycodeModel := models.GetModel(models.PAYCODE).(*models.PayCode)
	adddata := make(map[string]interface{})
	var colindex = 0
	//商品信息
	//order_iteminfo := make(map[string]interface{})
	if len(rowinfo) < 16 {
		return "模板列数不对"
	}
	colindex, itemname := getcolstr(colindex, rowinfo)
	iteminfo := itemmodel.GetInfoByField(self.dboper, "name", itemname)
	if iteminfo == nil {
		return getImportErr(colindex, rownum, fmt.Sprintf("商品:%s不存在", itemname))
	}
	colindex, itemnumstr := getcolstr(colindex, rowinfo)
	itemnum, err := strconv.Atoi(itemnumstr)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	basenum, err := strconv.Atoi(iteminfo[0]["basenum"].(string))
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	minnum, err := strconv.Atoi(iteminfo[0]["min_num"].(string))
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}

	if itemnum%basenum != 0 {
		return getImportErr(colindex, rownum, fmt.Sprintf("商品数量需要%d的倍数", basenum))
	}
	if itemnum < minnum {
		return getImportErr(colindex, rownum, fmt.Sprintf("商品数量需要大于%d", minnum))
	}
	adddata["itemid"] = iteminfo[0]["id"]
	adddata["num"] = itemnum
	adddata["supply_source"] = iteminfo[0]["supply_source"]
	// adddata["itemname"] = iteminfo[0]["name"]

	isonsale := iteminfo[0]["is_onsale"].(string)
	if isonsale == "0" {
		return getImportErr(colindex, rownum, fmt.Sprintf("商品id:%s已下架", adddata["itemid"].(string)))
	}

	colindex, itemcode := getcolstr(colindex, rowinfo)
	adddata["itemcode"] = itemcode
	err, specinfo := initSpecInfo(iteminfo[0], itemcode)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	storenum, err := strconv.Atoi(specinfo["store_num"].(string))
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	if storenum < itemnum {
		return getImportErr(colindex, rownum, fmt.Sprintf("商品id:%s 库存不足 需要:%d 库存:%d", adddata["itemid"].(string), itemnum, storenum))
	}

	adddata["user_id"] = self.uid

	userInfo := usermodel.GetInfoAndCache(self.dboper, adddata["user_id"].(string), true) //更新缓存
	unit_pricestr, err := itemmodel.GetItemPrice(userInfo["user_group"].(string), iteminfo[0], itemcode)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	adddata["unitprice"] = unit_pricestr
	unit_price, err := strconv.ParseFloat(unit_pricestr, 64)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}

	colindex, adddata["client_name"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_phone"] = getcolstr(colindex, rowinfo)
	client_phone := strings.Trim(adddata["client_phone"].(string), "#")
	adddata["client_phone"] = client_phone

	colindex, adddata["client_address"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_provice"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_city"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_area"] = getcolstr(colindex, rowinfo)
	colindex, adddata["idnum"] = getcolstr(colindex, rowinfo)
	idnum := strings.Trim(adddata["idnum"].(string), "#")
	if idnum != "" && len(idnum) < 18 {
		return getImportErr(colindex, rownum, "身份证格式不对")
	}

	needidnum := iteminfo[0]["idnum_need"].(string)
	if needidnum == "1" {
		if libs.CheckIdNum(idnum) == false {
			return getImportErr(colindex, rownum, "身份证格式不对")
		}
	}

	adddata["idnum"] = idnum

	colindex, adddata["idnumpic1"] = getcolstr(colindex, rowinfo)
	colindex, adddata["idnumpic2"] = getcolstr(colindex, rowinfo)
	colindex, adddata["send_user_name"] = getcolstr(colindex, rowinfo)
	colindex, adddata["send_user_phone"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_info"] = getcolstr(colindex, rowinfo)

	colindex, viptypestr := getcolstr(colindex, rowinfo)
	viptype, err := strconv.Atoi(viptypestr)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	if viptype < libs.Order_type_min || viptype > libs.Order_type_max {
		return getImportErr(colindex, rownum, "vip类型不对")
	}
	adddata["order_vip_type"] = viptypestr

	freight_price := getFreight(itemnum, adddata)
	service_price := getServiceMoney(itemnum, viptype)
	totalprice := getOrderPrice(itemnum, unit_price) + float64(freight_price) + float64(service_price)
	adddata["total_price"] = totalprice
	adddata["freight_price"] = freight_price
	adddata["service_price"] = service_price

	adddata["order_time"] = time.Now().Unix()

	adddata["status"] = libs.OrderStatusWaitPay
	idstr, err := getOrderid()
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	adddata["id"] = idstr

	adddata["pay_id"], err = paycodeModel.AddPayCodeByOrder(self.dboper, idstr, totalprice, self.uid)
	if err != nil {
		return err.Error()
	}
	keys, values := db.SqlGetInsertInfo(adddata)
	logs.Info("values:%s", values)

	_, err = self.dboper.Raw(fmt.Sprintf("insert into %s (%s) values (%s)", self.model.TableName(), keys, values)).Exec()
	if err != nil {
		//return err.Error()
		return getImportErr(colindex, rownum, err.Error())
	}
	self.AddLog(fmt.Sprintf("adddata:%+v", adddata))
	return ""
}

func initSpecInfo(iteminfo db.Params, code string) (error, map[string]interface{}) {
	logs.Info("itemcode:%s", code)
	resinfo := make(map[string]interface{})
	specstr := iteminfo["spec"].(string)
	if code != "" && specstr != "" {
		if code == iteminfo["code"].(string) {
			//直接是商品的
			resinfo["specname"] = ""
			resinfo["itempic"] = iteminfo["icon"]
			resinfo["store_num"] = iteminfo["store_num"]
			return nil, resinfo
		}
		specdata := new(models.ItemSpecData)
		err := json.Unmarshal([]byte(specstr), specdata)
		if err != nil {
			logs.Info("get spec err:", err.Error())
			return errors.New("商品数据错误"), nil
		}
		for _, detailitem := range specdata.DetailList {
			if detailitem.Code == code {
				var namestr = ""
				for _, specitem := range detailitem.Namearr {
					tagdata := specdata.SpecList[specitem.Specid].List[specitem.Tagid]
					namestr += tagdata.Name + ";"
				}
				resinfo["specname"] = namestr
				resinfo["pic"] = detailitem.Pic
				resinfo["store_num"] = detailitem.Store_num
				return nil, resinfo
			}
		}
		logs.Info("no found")
		return errors.New("商品编码不存在"), nil
	} else {
		resinfo["pic"] = iteminfo["icon"]
		resinfo["specname"] = ""
		resinfo["store_num"] = iteminfo["store_num"]
		return nil, resinfo
	}
}

//用户确认已支付
func (self *ShopOrderController) SetPayOk() {
	paycodeModel := models.GetModel(models.PAYCODE).(*models.PayCode)
	itemmodel := models.GetModel(models.SHOP_ITEM)
	self.CheckFieldExitAndReturn(self.postdata, "pay_type", "支付类型不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "pay_id", "支付码不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "money", "总金额不能为空")

	orderidtemp := self.postdata["order_id"]
	payid := self.postdata["pay_id"].(string)
	money, err := strconv.ParseFloat(self.postdata["money"].(string), 64)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	paycodeinfo := paycodeModel.GetInfoById(self.dboper, payid)
	if paycodeinfo == nil {
		self.AjaxReturnError(errors.New("支付码不存在"))
	}
	var orderlist []string
	if orderidtemp == nil {
		orderliststr := paycodeinfo["order_list"].(string)
		err = json.Unmarshal([]byte(orderliststr), &orderlist)
		if err != nil {
			self.AjaxReturnError(errors.WithStack(err))
		}
	} else {
		orderlist = append(orderlist, orderidtemp.(string))
	}

	var totalmoney = self.checkOrderList(orderlist)
	if money != totalmoney {
		self.AjaxReturnError(errors.New("金额不匹配"))
	}

	// db := orm.NewOrm()
	self.dboper.Begin()
	changedata := make(map[string]interface{})
	changedata["status"] = libs.OrderStatusWaitcheck
	changedata["pay_type"] = self.postdata["pay_type"]
	changedata["pay_id"] = payid
	for _, orderitem := range orderlist {
		logs.Info("update order:%s ", orderitem)
		err := self.updateSqlById(self, changedata, orderitem)
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}
		//增加商品销量
		orderinfo := self.model.GetInfoById(self.dboper, orderitem)
		itemnum := orderinfo["num"].(string)
		itemid := orderinfo["itemid"].(string)
		_, err = self.dboper.Raw(fmt.Sprintf("update %s set `sell_num`=`sell_num`+%s where `id`=?", itemmodel.TableName(), itemnum), itemid).Exec()
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}
	}
	self.dboper.Commit()
	self.AddLog(fmt.Sprintf("用户上传:%+v ", self.postdata))
	self.AjaxReturnSuccess("成功", nil)
}

//管理员审核未支付
func (self *ShopOrderController) CheckPayNO() {
	//orderinfo := self.checkOrderStatus(libs.OrderStatusWaitcheck)
	changedata := make(map[string]interface{})
	changedata["status"] = libs.OrderStatusWaitPay

	orderlist := self.postdata["id"].([]interface{})
	if len(orderlist) == 0 {
		self.AjaxReturnError(errors.New("id空"))
	}
	self.dboper.Begin()
	for _, id := range orderlist {
		idstr := id.(string)
		orderinfo := self.model.GetInfoById(self.dboper, idstr)
		if orderinfo == nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.New("订单不存在" + idstr))
		}

		statusstr := orderinfo["status"].(string)
		status, _ := strconv.Atoi(statusstr)
		if status != libs.OrderStatusWaitcheck {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.New("订单状态不对"))
		}

		err := self.updateSqlById(self, changedata, id)
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}
	}
	self.dboper.Commit()
	self.AddLog(fmt.Sprintf("用户上传:%+v ", self.postdata))
	self.AjaxReturnSuccess("成功", nil)

	//self.updateSqlByIdAndReturn(self.dboper,self,changedata, orderinfo["id"])
}

//获取运费
func getFreight(num int, adddata map[string]interface{}) int {
	provice := adddata["client_provice"].(string)
	supply := adddata["supply_source"].(string)
	if supply == libs.Supply_source_zhiyou {
		if strings.Contains(provice, "青海") || strings.Contains(provice, "西藏") || strings.Contains(provice, "新疆") {
			return int(math.Ceil(float64(num)/3.0)) * 50
		}
	} else if supply == libs.Supply_source_baoshui {
		if strings.Contains(provice, "青海") || strings.Contains(provice, "内蒙古") || strings.Contains(provice, "甘肃") ||
			strings.Contains(provice, "宁夏") || strings.Contains(provice, "西藏") || strings.Contains(provice, "新疆") {
			return num * 10
		}
	}
	return 0
}

func getServiceMoney(num int, viptype int) int {
	vipmoney := 0
	if viptype == libs.Order_type_photo {
		vipmoney = int(math.Ceil(float64(num)/3.0)) * 5
	} else if viptype == libs.Order_type_video {
		vipmoney = int(math.Ceil(float64(num)/3.0)) * 10
	}
	return vipmoney
}

func getOrderPrice(num int, unitprice float64) float64 {
	//logs.Info(" num: %d price:%v type:%d ", num, unitprice, viptype)
	totalmoney := unitprice * float64(num)
	return totalmoney
}

//补录身份证信息
func (self *ShopOrderController) UpdateIdNum() {
	self.CheckFieldExitAndReturn(self.postdata, "id", "id空")
	//self.CheckFieldExitAndReturn(self.postdata, "name", "姓名不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "idnum", "身份证号不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "idnumpic1", "身份证正面图片空")
	self.CheckFieldExitAndReturn(self.postdata, "idnumpic2", "身份证反面图片空")
	oldinfo := self.model.GetInfoByField(self.dboper, "id", self.postdata["id"])
	if oldinfo == nil {
		self.AjaxReturnError(errors.New("无此订单"))
	}

	if libs.CheckIdNum(self.postdata["idnum"].(string)) == false {
		self.AjaxReturnError(errors.New("身份证号不对"))
	}

	//if oldinfo[0]["client_name"].(string) != self.postdata["name"].(string) {
	//	self.AjaxReturnError(errors.Errorf("身份证号姓名:[%s]与收件人姓名[%s]不匹配",oldinfo[0]["client_name"].(string),self.postdata["name"].(string)))
	//}

	changedata := make(map[string]interface{})
	changedata["idnum"] = self.postdata["idnum"]
	changedata["idnumpic1"] = self.postdata["idnumpic1"]
	changedata["idnumpic2"] = self.postdata["idnumpic2"]
	self.updateSqlByIdAndReturn(self, changedata, self.postdata["id"])
}

//增加订单
//iteminfo{itemid:122,num:2,price:"234",specname:"dfsdf",code:"123",name:"dsfdf",pic:"sdfff"}
func (self *ShopOrderController) Add() {
	itemmodel := models.GetModel(models.SHOP_ITEM).(*models.ShopItem)
	datacheck := self.model.GetModelStruct()
	self.CheckExit(datacheck, self.postdata, false)
	adddata := libs.ClearMapByStruct(self.postdata, datacheck)
	paycodeModel := models.GetModel(models.PAYCODE).(*models.PayCode)
	self.dboper.Begin()
	itemarr, ok := self.postdata["item_info"].([]interface{})
	if (ok == false) && (len(itemarr) == 0) {
		self.dboper.Rollback()
		self.AjaxReturnError(errors.New("商品信息为空"))
	}
	adddata["order_time"] = time.Now().Unix()
	adddata["user_id"] = self.uid
	adddata["status"] = libs.OrderStatusWaitPay
	usermodel := models.GetModel(models.USER)
	userInfo := usermodel.GetInfoAndCache(self.dboper, self.uid, true) //更新缓存
	senddata := make(map[string]interface{})
	var idlist []string

	for _, item := range itemarr {
		iteminfo, ok := item.(map[string]interface{})
		if ok == false {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.New("商品信息错误"))
		}
		dbiteminfo := itemmodel.GetInfoById(self.dboper, iteminfo["itemid"])
		if dbiteminfo == nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.Errorf("商品id:%s不存在", iteminfo["itemid"].(string)))
		}
		isonsale := dbiteminfo["is_onsale"].(string)
		if isonsale == "0" {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.Errorf("商品id:%s已下架", iteminfo["itemid"].(string)))
		}

		needidnum := dbiteminfo["idnum_need"].(string)

		if needidnum == "1" {
			idnumstr, haveidnum := adddata["idnum"].(string)
			if haveidnum == true {
				idnumstr = strings.TrimSpace(idnumstr)
				adddata["idnum"] = idnumstr
			}
			if haveidnum == false || libs.CheckIdNum(idnumstr) == false {
				self.dboper.Rollback()
				self.AjaxReturnError(errors.Errorf("商品:%s需要身份证信息，请填入正确的身份证", iteminfo["name"].(string)))
			}
		}

		//获取单价
		unit_pricestr, err := itemmodel.GetItemPrice(userInfo["user_group"].(string), dbiteminfo, iteminfo["code"].(string))
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
			return
		}

		idstr, err := getOrderid()
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
			return
		}
		adddata["id"] = idstr

		adddata["itemid"] = iteminfo["itemid"]
		adddata["specname"] = iteminfo["specname"]
		adddata["itemcode"] = iteminfo["code"]
		adddata["num"] = iteminfo["num"]
		adddata["itempic"] = iteminfo["pic"]
		adddata["unitprice"] = unit_pricestr
		adddata["supply_source"] = dbiteminfo["supply_source"]
		trueprice, err := strconv.ParseFloat(unit_pricestr, 64)
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
			return
		}
		num := iteminfo["num"].(float64)

		err, specinfo := initSpecInfo(dbiteminfo, iteminfo["code"].(string))
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}
		storenum, err := strconv.Atoi(specinfo["store_num"].(string))
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}
		if storenum < int(num) {
			//self.dboper.Rollback()
			self.AjaxReturnError(errors.Errorf("商品id:%s 库存不足 需要:%d 库存:%d", iteminfo["itemid"].(string), int(num), storenum))

		}

		basenum, err := strconv.Atoi(dbiteminfo["basenum"].(string))
		minnum, err := strconv.Atoi(dbiteminfo["min_num"].(string))
		if int(num)%basenum != 0 {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.Errorf("商品id:%s 数量不是%d的倍数", iteminfo["specname"].(string), basenum))
		}
		if int(num) < minnum {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.Errorf("商品id:%s 数量小于最小数量%d", iteminfo["specname"].(string), minnum))
		}
		viptype := int(adddata["order_vip_type"].(float64))
		freight_price := getFreight(int(num), adddata)
		service_price := getServiceMoney(int(num), viptype)
		adddata["freight_price"] = freight_price
		adddata["service_price"] = service_price
		totalprice := getOrderPrice(int(num), trueprice) + float64(freight_price) + float64(service_price)
		adddata["total_price"] = totalprice
		//logs.Info("iteminfo:%s", adddata["item_info"].(string))
		adddata["pay_id"], err = paycodeModel.AddPayCodeByOrder(self.dboper, idstr, totalprice, self.uid)
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
			return
		}

		idlist = append(idlist, idstr)
		keys, values := db.SqlGetInsertInfo(adddata)
		logs.Info("values:%s", values)

		_, err = self.dboper.Raw(fmt.Sprintf("insert into %s (%s) values (%s)", self.model.TableName(), keys, values)).Exec()
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
			return
		}

		self.AddLog(fmt.Sprintf("adddata:%+v", adddata))
	}
	err := self.dboper.Commit()
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	senddata["ids"] = idlist
	self.AfterSql(senddata, nil)
	self.AjaxReturnSuccess("", senddata)
	return
}

//用户关闭订单
func (self *ShopOrderController) ClientClose() {
	orderinfo := self.checkOrderId()
	userid := orderinfo["user_id"].(string)
	if userid != self.uid {
		self.AjaxReturnError(errors.New("只能关闭自己的订单"))
	}
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status != libs.OrderStatusWaitPay {
		self.AjaxReturnError(errors.New("已付款，不能直接关闭订单"))
	}

	self.closeOrder(self.postdata["id"].(string), "", libs.OrderCloseByClient)
}

//删除订单
func (self *ShopOrderController) ClientDelOrder() {
	orderinfo := self.checkOrderId()
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	userid := orderinfo["user_id"].(string)
	if userid != self.uid {
		self.AjaxReturnError(errors.New("只能删除自己的订单"))
	}
	if status != libs.OrderStatusClose && status != libs.OrderStatusOver {
		self.AjaxReturnError(errors.New("只能删除已结束订单"))
	}
	id := self.postdata["id"].(string)
	updateinfo := make(map[string]interface{})
	updateinfo["status"] = libs.OrderStatusDelete
	self.updateSqlByIdAndReturn(self, updateinfo, id)
}

//确认收货
func (self *ShopOrderController) ClientConfirmOrder() {
	self.ConfirmOrder(true)
}

func (self *ShopOrderController) AdminConfirmOrder() {
	self.ConfirmOrder(false)
}

func (self *ShopOrderController) ConfirmOrder(onlymyself bool) {
	orderinfo := self.checkOrderId()
	if onlymyself {
		userid := orderinfo["user_id"].(string)
		if userid != self.uid {
			self.AjaxReturnError(errors.New("只能确认自己的订单"))
		}
	}

	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status != libs.OrderStatusSend {
		self.AjaxReturnError(errors.New("只能确认已发货订单"))
	}
	id := self.postdata["id"].(string)
	updateinfo := make(map[string]interface{})
	updateinfo["status"] = libs.OrderStatusOver
	self.updateSqlByIdAndReturn(self, updateinfo, id)
}

//取消恳款
func (self *ShopOrderController) AdminCancelRefund() {
	self.CancelRefund(false)
}

func (self *ShopOrderController) ClientCancelRefund() {
	self.CancelRefund(true)
}

func (self *ShopOrderController) CancelRefund(onlymyself bool) {
	orderinfo := self.checkOrderId()
	if onlymyself {
		userid := orderinfo["user_id"].(string)
		if userid != self.uid {
			self.AjaxReturnError(errors.New("只能取消自己的订单"))
		}
	}

	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status != libs.OrderStatusRefund {
		self.AjaxReturnError(errors.New("状态不对"))
	}
	id := self.postdata["id"].(string)
	updateinfo := make(map[string]interface{})
	shipnum, okstr := orderinfo["shipment_num"].(string)
	if okstr == false || shipnum == "" {
		updateinfo["status"] = libs.OrderStatusWaitSend
	} else {
		updateinfo["status"] = libs.OrderStatusSend
	}

	self.updateSqlByIdAndReturn(self, updateinfo, id)
}

//申请退款
func (self *ShopOrderController) ClientRefundOrder() {
	orderinfo := self.checkOrderId()
	userid := orderinfo["user_id"].(string)
	if userid != self.uid {
		self.AjaxReturnError(errors.New("只能确认自己的订单"))
	}
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if (status != libs.OrderStatusWaitSend) && (status != libs.OrderStatusSend) {
		self.AjaxReturnError(errors.New("只有已付款订单才能退款"))
	}

	id := self.postdata["id"].(string)
	updateinfo := make(map[string]interface{})
	updateinfo["status"] = libs.OrderStatusRefund
	updateinfo["refund_info"] = self.postdata["refund_info"]
	self.updateSqlByIdAndReturn(self, updateinfo, id)
}

//关闭订单
func (self *ShopOrderController) Adminclose() {
	orderinfo := self.checkOrderId()
	self.CheckFieldExitAndReturn(self.postdata, "close_info", "关闭原因不能为空")
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status > libs.OrderStatusWaitPay {
		self.AjaxReturnError(errors.New("玩家已经付款，不能直接关闭订单"))
	}

	self.closeOrder(self.postdata["id"].(string), self.postdata["close_info"].(string), libs.OrderCloseByAdmin)
}

//确认退款
func (self *ShopOrderController) AdminRefundSure() {
	self.checkOrderStatus(libs.OrderStatusRefund)
	self.CheckFieldExitAndReturn(self.postdata, "close_info", "退款原因不能为空")

	self.closeOrder(self.postdata["id"].(string), self.postdata["close_info"].(string), libs.OrderCloseRefund)
}

func (self *ShopOrderController) checkOrderStatus(needstatus int) db.Params {
	orderinfo := self.checkOrderId()
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status != needstatus {
		self.AjaxReturnError(errors.New("订单状态不对"))
	}
	return orderinfo
}

func (self *ShopOrderController) checkOrderId() db.Params {
	self.CheckFieldExitAndReturn(self.postdata, "id", "id空")
	id := self.postdata["id"]
	orderinfo := self.model.GetInfoById(self.dboper, id.(string))
	if orderinfo == nil {
		self.AjaxReturnError(errors.New("订单不存在"))
	}
	return orderinfo
}

//关闭订单
func (self *ShopOrderController) closeOrder(id string, closeinfo string, closetype int) {
	updateinfo := make(map[string]interface{})
	updateinfo["status"] = libs.OrderStatusClose
	updateinfo["close_info"] = closeinfo
	updateinfo["close_time"] = time.Now().Unix()
	updateinfo["close_type"] = closetype
	self.updateSqlByIdAndReturn(self, updateinfo, id)
}

func (self *ShopOrderController) GetMyOrder() {
	var data = models.AllReqData{And: true}
	err := json.Unmarshal(self.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.Info(err.Error())
		self.AjaxReturn(libs.ErrorCode, err.Error(), nil)
		return
	}
	data.Search["order.user_id"] = self.uid
	self.AllExc(data)
}

func (self *ShopOrderController) checkOrderList(orderlist []string) float64 {
	var totalmoney = 0.0
	for _, orderitem := range orderlist {
		orderinfo := self.model.GetInfoById(self.dboper, orderitem)
		if orderinfo == nil {
			self.AjaxReturnError(errors.New("订单" + orderitem + "不存在"))
		}
		statusstr := orderinfo["status"].(string)
		status, _ := strconv.Atoi(statusstr)
		if status != libs.OrderStatusWaitPay {
			self.AjaxReturnError(errors.New("订单" + orderitem + "状态不对"))
		}
		tempvalue, err := strconv.ParseFloat(orderinfo["total_price"].(string), 64)
		if err != nil {
			self.AjaxReturnError(errors.New("订单" + orderitem + "价格不对"))
		}
		totalmoney += tempvalue
	}
	return totalmoney
}

type OrderListType []string

type OrderTotalData struct {
	Orderlist OrderListType
}

func (p OrderListType) Len() int { return len(p) }
func (p OrderListType) Less(i, j int) bool {
	inta, err := strconv.Atoi(strings.TrimPrefix(p[i], ORDER_PRE))
	if err != nil {
		panic(err.Error())
	}
	intb, err := strconv.Atoi(strings.TrimPrefix(p[j], ORDER_PRE))
	if err != nil {
		panic(err.Error())
	}
	return inta < intb
}
func (p OrderListType) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

//获取汇总支付的支付码
func (self *ShopOrderController) GetTotalPayId() {
	paycodeModel := models.GetModel(models.PAYCODE).(*models.PayCode)

	getData := new(OrderTotalData)
	err := json.Unmarshal(self.Ctx.Input.RequestBody, getData)
	if err != nil {
		self.AjaxReturnError(errors.New("数据错误"))
	}

	if len(getData.Orderlist) == 0 {
		self.AjaxReturnError(errors.New("订单空"))
	}
	sort.Sort(OrderListType(getData.Orderlist))

	var totalmoney = self.checkOrderList(getData.Orderlist)
	payid, err := paycodeModel.AddPayCodeByOrderList(self.dboper, getData.Orderlist, totalmoney, self.uid)
	if err != nil {
		self.AjaxReturnError(errors.New("生成失败:" + err.Error()))
	}
	senddata := make(map[string]interface{})
	senddata["payid"] = payid
	self.AddLog(fmt.Sprintf("订单信息:%+v  生成的支付号:%s", getData.Orderlist, payid))
	self.AjaxReturn(libs.SuccessCode, "", senddata)
}

func (self *ShopOrderController) ExportCsv() {
	err, adddata := self.ExportCsvCommon()
	if err != nil {
		logs.Info("export err:%s", err.Error())
		self.AjaxReturnError(errors.WithStack(err))
	}
	self.AjaxReturnSuccess("", adddata)
}

func (self *ShopOrderController) ExportMyCsv() {
	search := self.postdata["search"].(map[string]interface{})
	search["order.user_id"] = self.uid
	err, adddata := self.ExportCsvCommonSearch(search)
	if err != nil {
		logs.Info("export err:%s", err.Error())
		self.AjaxReturnError(errors.WithStack(err))
	}
	self.AjaxReturnSuccess("", adddata)
}
