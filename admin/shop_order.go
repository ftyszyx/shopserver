package admin

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/models"
)

type ShopOrderController struct {
	BaseController
}

func (self *ShopOrderController) Edit() {
	self.EditCommon(self)
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

func (self *ShopOrderController) AfterSql(data map[string]interface{}, oldinfo orm.Params) {
	if self.method == "UpdateOrderShipNum" {
		self.AddLog(fmt.Sprintf("源订单信息：id:%s 物流单号:%s ||修改后： %+v", oldinfo["id"].(string), oldinfo["shipment_num"].(string), data))
	} else {
		self.AddLog(fmt.Sprintf("%+v", data))
	}

}

func (self *ShopOrderController) ExportToErp() {
	idarr := self.postdata["ids"].([]interface{})
	sendflag := self.postdata["sendflag"].(float64)
	// if ok == false {
	// 	self.AjaxReturnError("id错误")
	// }
	var dataarr []interface{}
	for _, id := range idarr {
		dataarr = append(dataarr, self.getErpExportData(id.(string)))
	}
	urlstr := beego.AppConfig.String("erp.url")
	token := beego.AppConfig.String("erp.shoptoken")
	req := httplib.Post(urlstr)
	senddata := make(map[string]interface{})
	senddata["data"] = dataarr
	senddata["token"] = token
	senddata["shop_id"] = beego.AppConfig.String("erp.shopid")
	reqbuf, err := json.Marshal(senddata)
	if err != nil {
		self.AjaxReturnError(err.Error())
	}

	req.Body(string(reqbuf))
	req.Header("Content-Type", "application/json")

	respdata, err := req.Bytes()
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	getData := new(RrpResult)
	logs.Info("get data:%s", string(respdata))
	err = json.Unmarshal(respdata, getData)
	if err != nil {
		logs.Info("parse data err")
		self.AjaxReturnError(err.Error())
	}
	logs.Info("get data:%v", getData)
	if getData.Code != "1" {
		self.AjaxReturnError(getData.Message)
	}
	db := orm.NewOrm()
	db.Begin()
	for _, dataitem := range getData.Data {
		//读每一行
		changedata := make(map[string]interface{})
		if sendflag == 1 {
			changedata["status"] = libs.OrderStatusSend
		}
		logisticstr, err := json.Marshal(dataitem.Logistics)
		if err != nil {
			self.AjaxReturnError(err.Error())
		}
		changedata["shipment_num"] = string(logisticstr)
		self.updateSqlById(self, changedata, dataitem.Id)
	}
	db.Commit()

	self.AjaxReturnSuccess("成功", nil)
}

type EditShipNumDataArr struct {
	List map[string]string
}

//修改物流单号
func (self *ShopOrderController) UpdateOrderShipNum() {
	getData := new(EditShipNumDataArr)
	err := json.Unmarshal(self.Ctx.Input.RequestBody, getData)
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	db := orm.NewOrm()
	db.Begin()
	for orderid, dataitem := range getData.List {
		//读每一行
		changedata := make(map[string]interface{})
		changedata["status"] = libs.OrderStatusSend
		// logisticstr, err := json.Marshal(dataitem)
		// if err != nil {
		// 	self.AjaxReturnError(err.Error())
		// }
		changedata["shipment_num"] = dataitem
		self.updateSqlById(self, changedata, orderid)
	}
	db.Commit()
	self.AjaxReturnSuccess("成功", nil)
}

//导出到erp
func (self *ShopOrderController) getErpExportData(orderid string) map[string]interface{} {
	usermodel := models.GetModel(models.USER)
	orderinfo := self.model.GetInfoById(orderid)
	if orderinfo == nil {
		self.AjaxReturnError("订单不存在")
	}
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status != libs.OrderStatusWaitSend {
		self.AjaxReturnError("订单状态不对")
	}
	sendata := make(map[string]interface{})
	sendata["shop_order"] = orderinfo["id"]

	sendata["pay_time"] = orderinfo["pay_time"]
	userinfo := usermodel.GetInfoAndCache(orderinfo["user_id"].(string), false)
	sendata["customer_account"] = userinfo["account"]
	sendata["customer_name"] = orderinfo["client_name"]
	sendata["customer_addr"] = orderinfo["client_address"]
	sendata["customer_province"] = orderinfo["client_provice"]
	sendata["customer_city"] = orderinfo["client_city"]
	sendata["customer_area"] = orderinfo["client_area"]
	sendata["user_id_number"] = orderinfo["idnum"]
	sendata["client_phone"] = orderinfo["client_phone"]
	sendata["user_info"] = orderinfo["client_info"]
	sendata["sell_vip_type"] = orderinfo["order_vip_type"]
	sendata["sell_type"] = 0

	iteminfostr := orderinfo["item_info"].(string)
	var iteminfo map[string]interface{}
	err := json.Unmarshal([]byte(iteminfostr), &iteminfo)
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	sendata["item_code"] = iteminfo["code"]
	sendata["num"] = iteminfo["num"]
	// totalprice := strconv.Atoi(iteminfo["total_price"].(string))
	sendata["total_price"] = orderinfo["total_price"]
	sendata["unit_price"] = iteminfo["price"]
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
var alphaNum = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`)

func getOrderid() (string, string, error) {
	nowtime := time.Now()
	updateOrderTime(nowtime.Unix())
	orderid := fmt.Sprintf("%d%d", nowtime.Unix(), orderidinfo.Num)
	var ordernum int64
	ordernum = nowtime.Unix() * int64(orderidinfo.Num)
	buyid := getPayId(ordernum)
	return orderid, buyid, nil
}

func getPayId(num int64) string {
	var getstr []byte
	for num >= 36 {
		remain := num % 36
		num = num / 36
		getstr = append(getstr, alphaNum[remain])
	}
	getstr = append(getstr, alphaNum[num])
	return string(getstr)
}

//批量下单
func (self *ShopOrderController) OrdersUpload() {

	logs.Info("OrderUpload")
	err, fileinfo := self.upload()
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	logs.Info("OrderUpload:%+v", fileinfo)
	filetype := fileinfo["filetype"].(string)
	if filetype != "csv" {
		self.AjaxReturnError("表格格式错误，只支持CSV")
	}
	fielpath := fileinfo["filePath"].(string)

	fileio, err := os.Open(fielpath)

	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	defer fileio.Close()

	reader := csv.NewReader(fileio)

	_, err = reader.Read()
	if err != nil {
		self.AjaxReturnError(err.Error())
	}
	rownum := 0
	db := orm.NewOrm()
	db.Begin()
	for {
		//读每一行
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			self.AjaxReturnError(err.Error())
		}
		rownum++
		if rownum > 9000 {
			self.AjaxReturnError("超过最大行数9000")
		}
		errstr := self.AddOneRow(rownum, record, db)
		if errstr != "" {
			self.AjaxReturnError(errstr)
		}
	}
	db.Commit()
	self.AjaxReturnSuccess("成功", nil)
}

func getcolstr(col int, rowinfo []string) (int, string) {
	utf8byte, err := libs.GbkToUtf8([]byte(rowinfo[col]))
	if err != nil {
		logs.Info(err.Error())
	}
	return col + 1, strings.TrimSpace(string(utf8byte))
}

func getImportErr(col int, row int, msg string) string {
	return fmt.Sprintf("第%d行 第%d列 错误:%s", row+1, col, msg)
}

func initSpecInfo(order_iteminfo map[string]interface{}, iteminfo orm.Params, code string) bool {
	logs.Info("itemcode:%s", code)
	specstr := iteminfo["spec"].(string)

	if code != "" && specstr != "" {
		if code == iteminfo["code"].(string) {
			//直接是商品的
			order_iteminfo["specname"] = ""
			order_iteminfo["pic"] = iteminfo["icon"].(string)
			order_iteminfo["price"] = iteminfo["price"].(string)
			return true
		}
		var sepcdata map[string]interface{}
		err := json.Unmarshal([]byte(specstr), &sepcdata)
		if err != nil {
			logs.Info("get spec err:", err.Error())
			return false
		}
		detaillist := sepcdata["detailList"].([]interface{})
		speclist := sepcdata["specList"].([]interface{})
		logs.Info("detaillist:%+v", detaillist)
		logs.Info("speclist:%+v", speclist)
		for _, value := range detaillist {
			valuedata := value.(map[string]interface{})

			codenum := valuedata["code"].(string)
			// logs.Info("codenum:%s", codenum)
			if codenum == code {
				namearrdata := valuedata["namearr"].([]interface{})
				var namestr = ""
				for _, specitem := range namearrdata {
					specitemdata := specitem.(map[string]interface{})
					specid := int(specitemdata["specid"].(float64))
					tagid := int(specitemdata["tagid"].(float64))
					specdata := speclist[specid].(map[string]interface{})
					taglist := specdata["list"].([]interface{})
					tagdata := taglist[tagid].(map[string]interface{})
					namestr += tagdata["name"].(string) + ";"

				}
				order_iteminfo["specname"] = namestr
				order_iteminfo["pic"] = valuedata["pic"]
				order_iteminfo["price"] = valuedata["price"]
				return true
			}
		}
		logs.Info("no found")
		return false
	} else {
		order_iteminfo["pic"] = iteminfo["icon"]
		order_iteminfo["specname"] = ""
		order_iteminfo["price"] = iteminfo["price"]
		return true
	}
}

//确认已支付
func (self *ShopOrderController) SetPayOk() {
	orderinfo := self.checkOrderStatus(libs.OrderStatusWaitPay)
	changedata := make(map[string]interface{})
	changedata["status"] = libs.OrderStatusWaitcheck
	self.updateSqlById(self, changedata, orderinfo["id"])
}

//审核支付
func (self *ShopOrderController) CheckPayOk() {
	orderinfo := self.checkOrderStatus(libs.OrderStatusWaitcheck)
	changedata := make(map[string]interface{})
	changedata["status"] = libs.OrderStatusWaitSend
	changedata["pay_time"] = time.Now().Unix()
	self.updateSqlById(self, changedata, orderinfo["id"])
}

//审核未支付
func (self *ShopOrderController) CheckPayNO() {
	orderinfo := self.checkOrderStatus(libs.OrderStatusWaitcheck)
	changedata := make(map[string]interface{})
	changedata["status"] = libs.OrderStatusWaitPay
	self.updateSqlById(self, changedata, orderinfo["id"])
}

func (self *ShopOrderController) AddOneRow(rownum int, rowinfo []string, db orm.Ormer) string {
	itemmodel := models.GetModel(models.SHOP_ITEM)
	usermodel := models.GetModel(models.USER)
	adddata := make(map[string]interface{})
	var colindex = 0
	//商品信息
	order_iteminfo := make(map[string]interface{})
	colindex, itemname := getcolstr(colindex, rowinfo)
	iteminfo := itemmodel.GetInfoByField("name", itemname)
	if iteminfo == nil {
		return getImportErr(colindex, rownum, fmt.Sprintf("商品:%s不存在", itemname))
	}
	colindex, itemnumstr := getcolstr(colindex, rowinfo)
	itemnum, err := strconv.Atoi(itemnumstr)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	order_iteminfo["itemid"] = iteminfo[0]["id"]
	order_iteminfo["num"] = itemnum
	order_iteminfo["name"] = iteminfo[0]["name"]

	colindex, itemcode := getcolstr(colindex, rowinfo)
	order_iteminfo["code"] = itemcode
	if initSpecInfo(order_iteminfo, iteminfo[0], itemcode) == false {
		return getImportErr(colindex, rownum, "商品编码错误")
	}
	iteminfostr, err := json.Marshal(order_iteminfo)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	unit_price, err := strconv.ParseFloat(order_iteminfo["price"].(string), 64)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	adddata["total_price"] = float64(itemnum) * unit_price
	//买家账号
	colindex, username := getcolstr(colindex, rowinfo)
	username = strings.TrimSpace(username)
	if username == "" {
		adddata["user_id"] = self.uid
	} else {
		userinfo := usermodel.GetInfoByField("account", username)
		if userinfo == nil {
			return getImportErr(colindex, rownum, "账号名错误")
		} else {
			adddata["user_id"] = userinfo[0]["id"]
		}
	}

	adddata["item_info"] = string(iteminfostr)
	colindex, adddata["client_name"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_phone"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_address"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_provice"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_city"] = getcolstr(colindex, rowinfo)
	colindex, adddata["client_area"] = getcolstr(colindex, rowinfo)
	colindex, adddata["idnum"] = getcolstr(colindex, rowinfo)
	adddata["idnum"] = strings.Trim(adddata["idnum"].(string), "#")
	colindex, adddata["idnumpic1"] = getcolstr(colindex, rowinfo)
	colindex, adddata["idnumpic2"] = getcolstr(colindex, rowinfo)
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

	adddata["order_time"] = time.Now().Unix()

	adddata["status"] = libs.OrderStatusWaitPay
	idstr, payid, err := getOrderid()
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	adddata["id"] = idstr
	adddata["pay_id"] = payid
	keys, values := libs.SqlGetInsertInfo(adddata)
	logs.Info("values:%s", values)

	_, err = db.Raw(fmt.Sprintf("insert into %s (%s) values (%s)", self.model.TableName(), keys, values)).Exec()
	if err != nil {
		return err.Error()
	}
	return ""
}

//补录身份证信息
func (self *ShopOrderController) UpdateIdNum() {
	self.CheckFieldExit(self.postdata, "id", "id空")
	self.CheckFieldExit(self.postdata, "idnum", "身份证号不能为空")
	self.CheckFieldExit(self.postdata, "idnumpic1", "身份证正面图片空")
	self.CheckFieldExit(self.postdata, "idnumpic2", "身份证反面图片空")
	changedata := make(map[string]interface{})
	changedata["idnum"] = self.postdata["idnum"]
	changedata["idnumpic1"] = self.postdata["idnumpic1"]
	changedata["idnumpic2"] = self.postdata["idnumpic2"]
	self.updateSqlById(self, changedata, self.postdata["id"])
}

//增加订单
//iteminfo{itemid:122,num:2,price:"234",specname:"dfsdf",code:"123",name:"dsfdf",pic:"sdfff"}
func (self *ShopOrderController) Add() {
	itemmodel := models.GetModel(models.SHOP_ITEM)
	datacheck := self.model.GetModelStruct()
	self.CheckExit(datacheck, self.postdata, false)
	adddata := libs.ClearMapByStruct(self.postdata, datacheck)
	o := orm.NewOrm()
	o.Begin()
	itemarr, ok := self.postdata["item_info"].([]interface{})
	if (ok == false) && (len(itemarr) == 0) {
		o.Rollback()
		self.AjaxReturnError("商品信息为空")
	}

	adddata["order_time"] = time.Now().Unix()
	adddata["user_id"] = self.uid
	adddata["status"] = libs.OrderStatusWaitPay
	senddata := make(map[string]interface{})
	var idlist []string

	for _, item := range itemarr {
		iteminfo, ok := item.(map[string]interface{})
		if ok == false {
			// o.Rollback()
			self.AjaxReturnError("商品信息错误")
		}
		if itemmodel.CheckExit("id", iteminfo["itemid"]) == false {
			// o.Rollback()
			self.AjaxReturnError(fmt.Sprintf("商品id:%s不存在", iteminfo["itemid"].(string)))
		}
		idstr, payid, err := getOrderid()
		if err != nil {
			self.AjaxReturnError(err.Error())
			return
		}
		adddata["id"] = idstr
		adddata["pay_id"] = payid
		iteminfostr, err := json.Marshal(iteminfo)
		if err != nil {
			self.AjaxReturnError(err.Error())
			return
		}
		adddata["item_info"] = string(iteminfostr)
		price := iteminfo["price"].(float64)
		num := iteminfo["num"].(float64)
		adddata["total_price"] = price * num
		//logs.Info("iteminfo:%s", adddata["item_info"].(string))
		idlist = append(idlist, idstr)
		keys, values := libs.SqlGetInsertInfo(adddata)
		logs.Info("values:%s", values)

		_, err = o.Raw(fmt.Sprintf("insert into %s (%s) values (%s)", self.model.TableName(), keys, values)).Exec()
		if err != nil {
			self.AjaxReturnError(err.Error())
			return
		}
	}
	o.Commit()
	senddata["ids"] = idlist
	self.AfterSql(senddata, nil)
	self.AjaxReturn(libs.SuccessCode, "", senddata)
	return
}

//用户关闭订单
func (self *ShopOrderController) ClientClose() {
	orderinfo := self.checkOrderId()
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status != libs.OrderStatusWaitPay {
		self.AjaxReturnError("已付款，不能直接关闭订单")
	}

	self.closeOrder(self.postdata["id"].(string), "", libs.OrderCloseByClient)
}

func (self *ShopOrderController) ClientDelOrder() {
	orderinfo := self.checkOrderId()
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status != libs.OrderStatusClose {
		self.AjaxReturnError("只能删除已关闭订单")
	}
	id := self.postdata["id"].(string)
	updateinfo := make(map[string]interface{})
	updateinfo["status"] = libs.OrderStatusDelete
	self.updateSqlById(self, updateinfo, id)
}

func (self *ShopOrderController) ClientRefundOrder() {
	orderinfo := self.checkOrderId()
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if (status != libs.OrderStatusWaitSend) && (status != libs.OrderStatusSend) {
		self.AjaxReturnError("只有已付款订单才能退款")
	}

	id := self.postdata["id"].(string)
	updateinfo := make(map[string]interface{})
	updateinfo["status"] = libs.OrderStatusRefund
	updateinfo["refund_info"] = self.postdata["refund_info"]
	self.updateSqlById(self, updateinfo, id)
}

//关闭订单
func (self *ShopOrderController) Adminclose() {
	orderinfo := self.checkOrderId()
	self.CheckFieldExit(self.postdata, "close_info", "关闭原因不能为空")
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status > libs.OrderStatusWaitPay {
		self.AjaxReturnError("玩家已经付款，不能直接关闭订单")
	}

	self.closeOrder(self.postdata["id"].(string), self.postdata["close_info"].(string), libs.OrderCloseByAdmin)
}

//确认退款
func (self *ShopOrderController) AdminRefundSure() {
	self.checkOrderStatus(libs.OrderStatusRefund)
	self.CheckFieldExit(self.postdata, "close_info", "关闭原因不能为空")

	self.closeOrder(self.postdata["id"].(string), self.postdata["close_info"].(string), libs.OrderCloseRefund)
}

func (self *ShopOrderController) checkOrderStatus(needstatus int) orm.Params {
	orderinfo := self.checkOrderId()
	statusstr := orderinfo["status"].(string)
	status, _ := strconv.Atoi(statusstr)
	if status != needstatus {
		self.AjaxReturnError("订单状态不对")
	}
	return orderinfo
}

func (self *ShopOrderController) checkOrderId() orm.Params {
	self.CheckFieldExit(self.postdata, "id", "id空")
	id := self.postdata["id"]
	orderinfo := self.model.GetInfoById(id.(string))
	if orderinfo == nil {
		self.AjaxReturnError("订单不存在")
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
	self.updateSqlById(self, updateinfo, id)
}

func (self *ShopOrderController) GetMyOrder() {
	var data = AllReqData{And: true}
	err := json.Unmarshal(self.Ctx.Input.RequestBody, &data)
	if err != nil {
		logs.Info(err.Error())
		self.AjaxReturn(libs.ErrorCode, err.Error(), nil)
		return
	}
	data.Search["user_id"] = self.uid
	self.AllExc(data)
}

func (self *ShopOrderController) ExportCsv() {
	self.ExportCsvCommon()
}
