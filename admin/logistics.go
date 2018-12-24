package admin

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/baiduai"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

type LogisticsController struct {
	BaseController
}

func (self *LogisticsController) BeforeSql(data map[string]interface{}) error {
	if self.method == "Add" {
		client_phone := data["client_phone"]
		client_name := data["client_name"]
		if client_phone == nil || client_name == nil {
			return errors.New("信息错误")
		}
		if client_phone.(string) != "" && client_name.(string) != "" && (data["idnumpic1"] == nil || data["idnumpic2"] == nil || data["user_id_number"] == nil) {
			//填入用户之前的数据
			olddatalist, err := self.model.GetInfoByWhere(self.dboper, fmt.Sprintf("`client_phone`='%s' and `client_name`=%s ", db.SqlGetString(client_phone), db.SqlGetString(client_name)))
			if err != nil {
				return err
			}
			if olddatalist != nil {
				data["idnum"] = olddatalist[0]["idnum"]
				data["idnumpic1"] = olddatalist[0]["idnumpic1"]
				data["idnumpic2"] = olddatalist[0]["idnumpic2"]
			}
		}
		curtime := time.Now().Unix()
		data["build_time"] = curtime

		data["logistics_task_starttime"] = curtime

	}
	return nil
}

func (self *LogisticsController) AfterSql(data map[string]interface{}, oldinfo db.Params) error {
	if self.method == "ClientChangeInfo" {
		self.AddLog(fmt.Sprintf("oldinfo:%+v change:%#v", oldinfo, data))
	} else {
		self.AddLog(fmt.Sprintf("change:%+v ", self.postdata))

	}

	return nil
}

func (self *LogisticsController) Add() {
	self.AddCommonAndReturn(self)
}

func (self *LogisticsController) Edit() {
	self.EditCommonAndReturn(self)
}

func (self *LogisticsController) Del() {
	self.DelCommonAndReturn(self)
}

func (self *LogisticsController) ExportCsv() {
	self.ExportCsvCommonAndReturn()
}

type EditLogicAPI struct {
	List []models.ErpShipdata
}

//批量下单
func (self *LogisticsController) UploadeLogistics() {
	err := self.UploadeCSV(self)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	} else {
		self.AjaxReturnSuccessNull()
	}
}

func (self *LogisticsController) AddOneRow(rownum int, rowinfo []string) string {

	var colindex = 0
	logicmodel := self.model.(*models.Logistics)
	if len(rowinfo) < 11 {
		return "模板列数不对"
	}
	var adddata models.ErpShipdata

	colindex, adddata.Logistics = getcolstr(colindex, rowinfo)
	if adddata.Logistics == "" {
		return getImportErr(colindex, rownum, "物流号不能为空")
	}
	if strings.HasPrefix(adddata.Logistics, "AB") == false {
		return getImportErr(colindex, rownum, "只能导入AB**AU格式单号")
	}
	if strings.HasSuffix(adddata.Logistics, "AU") == false {
		return getImportErr(colindex, rownum, "只能导入AB**AU格式单号")
	}
	colindex, inter_ship_companyname := getcolstr(colindex, rowinfo)
	if inter_ship_companyname != "" {
		inter_ship_companycode, ok := models.LogisticNameMap[inter_ship_companyname]
		if ok == false {
			return getImportErr(colindex, rownum, "快递公司不存在:"+inter_ship_companyname)
		}
		adddata.Internal_ship_company_code = inter_ship_companycode
	}
	colindex, adddata.Internal_ship_num = getcolstr(colindex, rowinfo)
	colindex, adddata.Customer_name = getcolstr(colindex, rowinfo)
	colindex, adddata.Client_phone = getcolstr(colindex, rowinfo)
	adddata.Client_phone = strings.Trim(adddata.Client_phone, "#")
	colindex, adddata.Client_address = getcolstr(colindex, rowinfo)
	colindex, adddata.User_id_number = getcolstr(colindex, rowinfo)
	adddata.User_id_number = strings.Trim(adddata.User_id_number, "#")
	colindex, adddata.Idnumpic1 = getcolstr(colindex, rowinfo)
	colindex, adddata.Idnumpic2 = getcolstr(colindex, rowinfo)
	var taskname = ""
	colindex, taskname = getcolstr(colindex, rowinfo)
	var err error
	var taskidstr = ""
	if taskname != "" {
		logisticssTaskModel := models.GetModel(models.LOGISTICS_TASK)
		info := logisticssTaskModel.GetInfoByField(self.dboper, "name", taskname)
		if info == nil {
			return getImportErr(colindex, rownum, "物流进度不存在:"+taskname)
		}
		taskidstr = info[0]["id"].(string)
	}
	colindex, adddata.Logistics_task_starttime = getcolstr(colindex, rowinfo)
	if adddata.Logistics_task_starttime != "" {

		err, starttime := libs.ParseTime(adddata.Logistics_task_starttime)
		if err != nil {
			return getImportErr(colindex, rownum, err.Error())
		}
		adddata.Logistics_task_starttime = strconv.FormatInt(starttime.Unix(), 10)
		logs.Info("get time:%s", adddata.Logistics_task_starttime)
	}

	err = logicmodel.AddList(self.dboper, []models.ErpShipdata{adddata}, taskidstr)
	if err != nil {
		return getImportErr(colindex, rownum, err.Error())
	}
	self.AddLog(fmt.Sprintf("adddata:%+v", adddata))
	return ""
}

//导物流接口
func (self *LogisticsController) AddLogicAPI() {
	//logs.Info("AddLogicAPI:%+v", self.postdata)
	logicmodel := self.model.(*models.Logistics)
	dataarr := new(EditLogicAPI)
	err := json.Unmarshal(self.Ctx.Input.RequestBody, dataarr)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}

	err = logicmodel.AddList(self.dboper, dataarr.List, "")
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	idarrlist := make([]string, len(dataarr.List))
	for _, logicsinfo := range dataarr.List {
		idarrlist = append(idarrlist, logicsinfo.Logistics)
	}
	self.AddLog(fmt.Sprintf("list:%+v ", idarrlist))
	self.AjaxReturnSuccessNull()

}

//更新海外物流进度
func (self *LogisticsController) UpdateTask() {
	idarr := self.postdata["ids"].([]interface{})
	senddata := make(map[string]interface{})
	logicmodel := self.model.(*models.Logistics)
	oklist, errlist, err := logicmodel.UpdateTask(self.dboper, idarr)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	senddata["oklist"] = oklist
	senddata["errlist"] = errlist
	self.AddLog(fmt.Sprintf("list:%+v ", senddata))
	self.AjaxReturnSuccess("", senddata)
}

//更新所有
func (self *LogisticsController) UpdateAllTask() {
	curtime := time.Now().Unix()
	senddata := make(map[string]interface{})
	res, err := self.model.GetInfoByWhere(self.dboper, fmt.Sprintf("`logistics_task_starttime`<%d and  `state` <> %d and `id` like ", curtime, libs.ShipOverseaOverValue)+"'AB%AU'")
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	if res != nil {
		logicmodel := self.model.(*models.Logistics)
		//oklist, errlist, err := logicmodel.UpdateTaskByDataList(self.dboper, res)
		oklist, _, err := logicmodel.UpdateTaskByDataList(self.dboper, res)
		if err != nil {
			self.AjaxReturnError(errors.WithStack(err))
		}
		senddata["oklist"] = oklist
		// senddata["errlist"] = errlist
	}
	self.AddLog(fmt.Sprintf("list:%+v ", senddata))
	self.AjaxReturnSuccess("", senddata)
}

//查询信息
func (self *LogisticsController) GetLogicsInfo() {
	id := self.postdata["id"].(string)
	if strings.HasPrefix(id, "AB") == false {
		self.AjaxReturnErrorMsg("不存在此物流信息")
	}
	if strings.HasSuffix(id, "AU") == false {
		self.AjaxReturnErrorMsg("不存在此物流信息")
	}
	model := self.model.(*models.Logistics)
	res, status, err, idinfo := model.GetLogicsInfo(self.dboper, id)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	senddata := make(map[string]interface{})
	senddata["state"] = status
	senddata["data"] = res
	senddata["idinfo"] = idinfo
	logs.Info("data:%+v ", senddata)
	self.AjaxReturnSuccess("", senddata)
}

//同步到Erp
func (self *LogisticsController) SyncErpData() {
	idlist := self.postdata["id"].([]interface{})
	urlstr := beego.AppConfig.String("erp.url") + "Sell/syncLogistics"
	token := beego.AppConfig.String("erp.shoptoken")
	req := httplib.Post(urlstr)
	logs.Info("url:%s", urlstr)
	if len(idlist) == 0 {
		self.AjaxReturnError(errors.New("同步项空"))
	}
	senddata := make(map[string]interface{})
	senddata["token"] = token
	senddata["shop_id"] = beego.AppConfig.String("erp.shopid")
	var datalist []interface{}
	for _, id := range idlist {
		info := self.model.GetInfoById(self.dboper, id)
		if info == nil {
			self.AjaxReturnError(errors.Errorf("id错误%+v", id))
		}
		itemdata := make(map[string]interface{})
		itemdata["idnum"] = info["idnum"]
		itemdata["client_name"] = info["client_name"]
		itemdata["idnumpic1"] = info["idnumpic1"]
		itemdata["idnumpic2"] = info["idnumpic2"]
		itemdata["id"] = info["id"]
		datalist = append(datalist, info)
	}
	senddata["list"] = datalist
	reqbuf, err := json.Marshal(senddata)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}

	req.Body(string(reqbuf))
	req.Header("Content-Type", "application/json")

	respdata, err := req.Bytes()
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	getData := new(RrpResult)
	// logs.Info("get data:%s", string(respdata))
	err = json.Unmarshal(respdata, getData)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	logs.Info("get data:%v", getData)
	if getData.Code != "1" {
		self.AjaxReturnError(errors.New(getData.Message))
	}

	//修改状态
	self.dboper.Begin()
	changedata := make(map[string]interface{})
	changedata["sync_erp_flag"] = 1
	for _, id := range idlist {
		err := self.updateSqlById(self, changedata, id)
		if err != nil {
			self.dboper.Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}
	}
	self.dboper.Commit()
	self.AjaxReturnSuccessNull()
}

//客户修改物流信息
func (self *LogisticsController) ClientChangeInfo() {
	self.CheckFieldExitAndReturn(self.postdata, "client_name", "收件人姓名不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "idnum", "收件人身份证号不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "idnumpic1", "收件人姓名身份证号正面图片不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "idnumpic2", "收件人姓名身份证号反面图片不能为空")
	self.CheckFieldExitAndReturn(self.postdata, "client_phone", "收件人电话不能为空")

	idnumpic1 := self.postdata["idnumpic1"].(string)
	idnumpic2 := self.postdata["idnumpic2"].(string)
	idnum := self.postdata["idnum"].(string)
	client_name := self.postdata["client_name"].(string)
	client_phone := self.postdata["client_phone"].(string)

	olddatalist, err := self.model.GetInfoByWhere(self.dboper, fmt.Sprintf("`client_phone`=%s and `client_name`=%s ", db.SqlGetString(client_phone), db.SqlGetString(client_name)))
	if err != nil {
		logs.Error(" err %+v", err)
		self.AjaxReturnError(errors.WithStack(err))
	}
	if olddatalist == nil {
		self.AjaxReturnError(errors.New("不存在该电话和收件人的物流信息"))
	}

	//检查身份证正面
	cardfrontres, err := baiduai.GlobalaiIDCarddata.GetIdResByUrl(idnumpic1, "front")
	if err != nil {
		logs.Info("Globalaidata err %+v", err)
		self.AjaxReturnError(errors.WithStack(err))
	}

	getname := cardfrontres.Words_result["姓名"].Words
	getidnum := cardfrontres.Words_result["公民身份号码"].Words
	if getidnum != idnum {
		self.AjaxReturnError(errors.New("身份证号不匹配"))
	}

	if getname != client_name {
		self.AjaxReturnError(errors.New("身份证号姓名不匹配"))
	}

	//检查身份证反面
	_, err = baiduai.GlobalaiIDCarddata.GetIdResByUrl(idnumpic2, "back")
	if err != nil {
		logs.Info("Globalaidata err %+v", err)
		self.AjaxReturnError(errors.WithStack(err))
	}

	dbpre, err := self.dboper.Raw(fmt.Sprintf("update %s set `idnum`=? ,`idnumpic1`=?,`idnumpic2`=?,`sync_erp_flag`=? where `id`=?", self.model.TableName())).Prepare()
	if err != nil {
		logs.Info("err %+v", err)
		self.AjaxReturnError(errors.WithStack(err))
	}

	defer dbpre.Close()
	for _, logisticsItem := range olddatalist {
		logisticsid := logisticsItem["id"].(string)
		_, err := dbpre.SetArgs(idnum, idnumpic1, idnumpic2, 0, logisticsid).Exec()
		if err != nil {
			logs.Info("err %+v", err)
			self.AjaxReturnError(errors.WithStack(err))
		}
	}
	self.uid = "1"
	self.AddLog(fmt.Sprintf("update postdata:%+v", self.postdata))
	self.AjaxReturnSuccessNull()
}
