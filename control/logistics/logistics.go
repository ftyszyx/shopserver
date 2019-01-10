package logistics

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
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/baiduai"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/logistics"
	"github.com/zyx/shop_server/models/names"
	"github.com/zyx/shop_server/models/shop"
)

type LogisticsController struct {
	base.BaseController
}

func (self *LogisticsController) BeforeSql(data map[string]interface{}) error {
	if self.GetMethod() == "Add" {
		client_phone := data["client_phone"]
		client_name := data["client_name"]
		if client_phone == nil || client_name == nil {
			return errors.New("信息错误")
		}
		if client_phone.(string) != "" && client_name.(string) != "" && (data["idnumpic1"] == nil || data["idnumpic2"] == nil || data["user_id_number"] == nil) {
			//填入用户之前的数据
			olddatalist, err := self.GetModel().GetInfoByWhere(self.GetDb(), fmt.Sprintf("`client_phone`='%s' and `client_name`=%s ", db.SqlGetString(client_phone), db.SqlGetString(client_name)))
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
	if self.GetMethod() == "ClientChangeInfo" {
		self.AddLog(fmt.Sprintf("oldinfo:%+v change:%#v", oldinfo, data))
	} else {
		self.AddLog(fmt.Sprintf("change:%+v ", self.GetPost()))

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
	List []logistics.ErpShipdata
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
	logicmodel := self.GetModel().(*logistics.Logistics)
	if len(rowinfo) < 11 {
		return "模板列数不对"
	}
	var adddata logistics.ErpShipdata

	colindex, adddata.Logistics = base.Getcolstr(colindex, rowinfo)
	if adddata.Logistics == "" {
		return base.GetImportErr(colindex, rownum, "物流号不能为空")
	}
	if strings.HasPrefix(adddata.Logistics, "AB") == false {
		return base.GetImportErr(colindex, rownum, "只能导入AB**AU格式单号")
	}
	if strings.HasSuffix(adddata.Logistics, "AU") == false {
		return base.GetImportErr(colindex, rownum, "只能导入AB**AU格式单号")
	}
	colindex, inter_ship_companyname := base.Getcolstr(colindex, rowinfo)
	if inter_ship_companyname != "" {
		inter_ship_companycode, ok := logistics.LogisticNameMap[inter_ship_companyname]
		if ok == false {
			return base.GetImportErr(colindex, rownum, "快递公司不存在:"+inter_ship_companyname)
		}
		adddata.Internal_ship_company_code = inter_ship_companycode
	}
	colindex, adddata.Internal_ship_num = base.Getcolstr(colindex, rowinfo)
	colindex, adddata.Customer_name = base.Getcolstr(colindex, rowinfo)
	colindex, adddata.Client_phone = base.Getcolstr(colindex, rowinfo)
	adddata.Client_phone = strings.Trim(adddata.Client_phone, "#")
	colindex, adddata.Client_address = base.Getcolstr(colindex, rowinfo)
	colindex, adddata.User_id_number = base.Getcolstr(colindex, rowinfo)
	adddata.User_id_number = strings.Trim(adddata.User_id_number, "#")
	colindex, adddata.Idnumpic1 = base.Getcolstr(colindex, rowinfo)
	colindex, adddata.Idnumpic2 = base.Getcolstr(colindex, rowinfo)
	var taskname = ""
	colindex, taskname = base.Getcolstr(colindex, rowinfo)
	var err error
	var taskidstr = ""
	if taskname != "" {
		logisticssTaskModel := models.GetModel(names.LOGISTICS_TASK)
		info := logisticssTaskModel.GetInfoByField(self.GetDb(), "name", taskname)
		if info == nil {
			return base.GetImportErr(colindex, rownum, "物流进度不存在:"+taskname)
		}
		taskidstr = info[0]["id"].(string)
	}
	colindex, adddata.Logistics_task_starttime = base.Getcolstr(colindex, rowinfo)
	if adddata.Logistics_task_starttime != "" {

		err, starttime := libs.ParseTime(adddata.Logistics_task_starttime)
		if err != nil {
			return base.GetImportErr(colindex, rownum, err.Error())
		}
		adddata.Logistics_task_starttime = strconv.FormatInt(starttime.Unix(), 10)
		logs.Info("get time:%s", adddata.Logistics_task_starttime)
	}

	err = logicmodel.AddList(self.GetDb(), []logistics.ErpShipdata{adddata}, taskidstr)
	if err != nil {
		return base.GetImportErr(colindex, rownum, err.Error())
	}
	self.AddLog(fmt.Sprintf("adddata:%+v", adddata))
	return ""
}

//导物流接口
func (self *LogisticsController) AddLogicAPI() {
	//logs.Info("AddLogicAPI:%+v", self.GetPost())
	logicmodel := self.GetModel().(*logistics.Logistics)
	dataarr := new(EditLogicAPI)
	err := json.Unmarshal(self.Ctx.Input.RequestBody, dataarr)
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}

	err = logicmodel.AddList(self.GetDb(), dataarr.List, "")
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
	idarr := self.GetPost()["ids"].([]interface{})
	senddata := make(map[string]interface{})
	logicmodel := self.GetModel().(*logistics.Logistics)
	oklist, errlist, err := logicmodel.UpdateTask(self.GetDb(), idarr)
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
	res, err := self.GetModel().GetInfoByWhere(self.GetDb(), fmt.Sprintf("`logistics_task_starttime`<%d and  `state` <> %d and `id` like ", curtime, libs.ShipOverseaOverValue)+"'AB%AU'")
	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	if res != nil {
		logicmodel := self.GetModel().(*logistics.Logistics)
		//oklist, errlist, err := logicmodel.UpdateTaskByDataList(self.GetDb(), res)
		oklist, _, err := logicmodel.UpdateTaskByDataList(self.GetDb(), res)
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
	id := self.GetPost()["id"].(string)
	if strings.HasPrefix(id, "AB") == false {
		self.AjaxReturnErrorMsg("不存在此物流信息")
	}
	if strings.HasSuffix(id, "AU") == false {
		self.AjaxReturnErrorMsg("不存在此物流信息")
	}
	model := self.GetModel().(*logistics.Logistics)
	res, status, err, idinfo := model.GetLogicsInfo(self.GetDb(), id)
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
	idlist := self.GetPost()["id"].([]interface{})
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
		info := self.GetModel().GetInfoById(self.GetDb(), id)
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
	getData := new(shop.ErpResult)
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
	self.GetDb().Begin()
	changedata := make(map[string]interface{})
	changedata["sync_erp_flag"] = 1
	for _, id := range idlist {
		err := self.UpdateSqlById(self, changedata, id)
		if err != nil {
			self.GetDb().Rollback()
			self.AjaxReturnError(errors.WithStack(err))
		}
	}
	self.GetDb().Commit()
	self.AjaxReturnSuccessNull()
}

//客户修改物流信息
func (self *LogisticsController) ClientChangeInfo() {
	self.CheckFieldExitAndReturn(self.GetPost(), "client_name", "收件人姓名不能为空")
	self.CheckFieldExitAndReturn(self.GetPost(), "idnum", "收件人身份证号不能为空")
	//self.CheckFieldExitAndReturn(self.GetPost(), "idnumpic1", "收件人姓名身份证号正面图片不能为空")
	//self.CheckFieldExitAndReturn(self.GetPost(), "idnumpic2", "收件人姓名身份证号反面图片不能为空")
	self.CheckFieldExitAndReturn(self.GetPost(), "client_phone", "收件人电话不能为空")

	idnum := self.GetPost()["idnum"].(string)

	if libs.CheckIdNum(idnum) == false {
		self.AjaxReturnError(errors.New("身份证格式不对"))
	}
	client_name := self.GetPost()["client_name"].(string)
	client_phone := self.GetPost()["client_phone"].(string)

	olddatalist, err := self.GetModel().GetInfoByWhere(self.GetDb(), fmt.Sprintf("`client_phone`=%s and `client_name`=%s ", db.SqlGetString(client_phone), db.SqlGetString(client_name)))
	if err != nil {
		logs.Error(" err %+v", err)
		self.AjaxReturnError(errors.WithStack(err))
	}
	if olddatalist == nil {
		self.AjaxReturnError(errors.New("不存在该电话和收件人的物流信息"))
	}

	idnumpic1, idok1 := self.GetPost()["idnumpic1"].(string)
	//检查身份证正面
	var have_idfront = false
	if idok1 == true && idnumpic1 != "" {
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
		have_idfront = true
	}
	var have_idback = false
	//检查身份证反面
	idnumpic2, idok2 := self.GetPost()["idnumpic2"].(string)
	if idok2 == true && idnumpic2 != "" {
		_, err = baiduai.GlobalaiIDCarddata.GetIdResByUrl(idnumpic2, "back")
		if err != nil {
			logs.Info("Globalaidata err %+v", err)
			self.AjaxReturnError(errors.WithStack(err))
		}
		have_idback = true

	}

	var sqlprestr = ""
	if have_idfront && have_idback {
		sqlprestr = fmt.Sprintf("update %s set `idnum`=? ,`idnumpic1`=?,`idnumpic2`=?,`sync_erp_flag`=? where `id`=?", self.GetModel().TableName())
	} else if have_idback && have_idfront == false {
		sqlprestr = fmt.Sprintf("update %s set `idnum`=? ,`idnumpic2`=?,`sync_erp_flag`=? where `id`=?", self.GetModel().TableName())
	} else if have_idback == false && have_idfront {
		sqlprestr = fmt.Sprintf("update %s set `idnum`=? ,`idnumpic1`=?,`sync_erp_flag`=? where `id`=?", self.GetModel().TableName())
	} else if have_idback == false && have_idfront == false {
		sqlprestr = fmt.Sprintf("update %s set `idnum`=? ,`sync_erp_flag`=? where `id`=?", self.GetModel().TableName())
	}

	dbpre, err := self.GetDb().Raw(sqlprestr).Prepare()
	if err != nil {
		logs.Info("err %+v", err)
		self.AjaxReturnError(errors.WithStack(err))
	}

	defer dbpre.Close()
	for _, logisticsItem := range olddatalist {
		logisticsid := logisticsItem["id"].(string)
		var err error
		if have_idfront && have_idback {
			_, err = dbpre.SetArgs(idnum, idnumpic1, idnumpic2, 0, logisticsid).Exec()
		} else if have_idback && have_idfront == false {
			_, err = dbpre.SetArgs(idnum, idnumpic2, 0, logisticsid).Exec()
		} else if have_idback == false && have_idfront {
			_, err = dbpre.SetArgs(idnum, idnumpic1, 0, logisticsid).Exec()
		} else if have_idback == false && have_idfront == false {
			_, err = dbpre.SetArgs(idnum, 0, logisticsid).Exec()
		}
		//_, err := dbpre.SetArgs(idnum, idnumpic1, idnumpic2, 0, logisticsid).Exec()
		if err != nil {
			logs.Info("err %+v", err)
			self.AjaxReturnError(errors.WithStack(err))
		}
	}
	self.SetUid("1")
	self.AddLog(fmt.Sprintf("update postdata:%+v", self.GetPost()))
	self.AjaxReturnSuccessNull()
}
