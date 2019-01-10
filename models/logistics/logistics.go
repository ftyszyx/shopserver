package logistics

// 物流
import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type Logistics struct {
	models.Model
}

type LogisticsData struct {
	Id                         string `empty:"物流号不能为空" `
	Info_oversea               string
	Internal_ship_company_code string
	internal_ship_num          string
	logistics_task_starttime   int
	state                      int
	logistics_task             string
	Idnumpic1                  string
	Idnumpic2                  string
	Client_name                string `empty:"收件人姓名不能空" `
	Client_phone               string `empty:"收件人电话不能空" `
	Client_address             string `empty:"收件人地址不能空" `
	idnum                      string
}

type ErpShipdata struct {
	Logistics                  string
	Idnumpic1                  string
	Idnumpic2                  string
	Customer_name              string
	User_id_number             string
	Internal_ship_company_code string
	Client_phone               string
	Client_address             string
	Internal_ship_num          string
	Logistics_task_starttime   string
	Ok                         bool
}

var LogisticCodeMap map[string]string
var LogisticNameMap map[string]string

func (self *Logistics) Init() {
	LogisticCodeMap = make(map[string]string)
	LogisticNameMap = make(map[string]string)
	for _, item := range libs.LogisticsCodeArr {
		LogisticCodeMap[item.Code] = item.Title
		LogisticNameMap[item.Title] = item.Code
	}
}

func (self *Logistics) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))

}

func (self *Logistics) GetModelStruct() interface{} {
	return LogisticsData{}
}

func (self *Logistics) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {

	fieldstr := ""

	return sql.Alias("logistics").Join(fieldstr)
}

func (self *Logistics) ExportNameProcess(name string, celldata interface{}, row db.Params) (string, error) {
	if celldata == nil {
		logs.Info("field %s is nil", name)
		return "", nil
	}
	value, ok := celldata.(string)
	if ok == false {
		return "", errors.New("upload file err:" + name + " not exit")
	}

	if name == "build_time" || name == "logistics_task_starttime" {
		return libs.FormatTableTime(value), nil
	} else if name == "client_phone" || name == "idnum" {
		return "\t" + value, nil
	} else if name == "state" {
		valuenum, err := strconv.Atoi(value)
		if err != nil {
			return "", errors.WithStack(err)
		}
		return libs.ShipOverseaStatusArr[valuenum], nil
	}

	return value, nil
}

func (self *Logistics) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"logistics.id":                         "id",
		"logistics.internal_ship_company_code": "internal_ship_company_code",
		"logistics.internal_ship_num":          "internal_ship_num",
		"logistics.build_time":                 "build_time",
		"logistics.logistics_task":             "logistics_task",
		"logistics.logistics_task_starttime":   "logistics_task_starttime",
		"logistics.state":                      "state",
		"logistics.idnum":                      "idnum",
		"logistics.client_name":                "client_name",
		"logistics.idnumpic1":                  "idnumpic1",
		"logistics.idnumpic2":                  "idnumpic2",
		"logistics.client_phone":               "client_phone",
		"logistics.client_address":             "client_address",
		// "logistics.is_del":                     "is_del",
		"logistics.sync_erp_flag": "sync_erp_flag",
	})
}

func (self *Logistics) BuildTask(dboper db.DBOperIO, taskid interface{}, curtime int64) string {
	logisticssTaskModel := models.GetModel(names.LOGISTICS_TASK)
	if taskid == nil {
		return ""
	}
	info := logisticssTaskModel.GetInfoAndCache(dboper, taskid.(string), false)
	if info == nil {
		return ""
	}
	taskliststr := info["tasklist"].(string)
	if taskliststr == "" {
		return ""
	}

	var taskList []LogisticsTaskItem
	err := json.Unmarshal([]byte(taskliststr), &taskList)
	if err != nil {

		return ""
	}
	if len(taskList) == 0 {
		return ""
	}

	var endtime = curtime
	for index, taskitem := range taskList {
		endtime += int64(taskitem.Use_time * 3600)
		//logs.Info("build time:%d", endtime)
		taskList[index].Time = endtime
	}

	tasknewstr, err := json.Marshal(taskList)
	if err != nil {
		logs.Error("BuildTask err" + err.Error())
		return ""
	}
	return string(tasknewstr)

}

//增加一条
func (self *Logistics) AddList(dboper db.DBOperIO, logicsarr []ErpShipdata, taskid string) error {

	// logs.Info("AddList :%#v", logicsarr)
	curtime := time.Now().Unix()
	configmodel := models.GetModel(names.CONFIG)
	configcache := configmodel.Cache()
	default_logistics_task := configcache.Get("logistics_default_task")
	if taskid != "" {
		default_logistics_task = taskid
	}

	// logs.Info("logistics_default_task:%v", default_logistics_task)

	for _, logicsinfo := range logicsarr {
		if strings.HasPrefix(logicsinfo.Logistics, "AB") == false {
			logicsinfo.Ok = false
			continue
		}
		if strings.HasSuffix(logicsinfo.Logistics, "AU") == false {
			logicsinfo.Ok = false
			continue
		}
		logisticssinfo := self.GetInfoById(dboper, logicsinfo.Logistics)
		var changedata = make(map[string]interface{})

		if logicsinfo.Idnumpic1 != "" {
			changedata["idnumpic1"] = logicsinfo.Idnumpic1
		}
		if logicsinfo.Idnumpic2 != "" {
			changedata["idnumpic2"] = logicsinfo.Idnumpic2
		}
		if logicsinfo.Customer_name != "" {
			changedata["client_name"] = logicsinfo.Customer_name
		}

		if logicsinfo.Logistics_task_starttime != "" {
			changedata["logistics_task_starttime"] = logicsinfo.Logistics_task_starttime
		}

		if logicsinfo.User_id_number != "" {
			changedata["idnum"] = logicsinfo.User_id_number
		}

		if logicsinfo.Client_address != "" {
			changedata["client_address"] = logicsinfo.Client_address
		}

		if logicsinfo.Client_phone != "" {
			changedata["client_phone"] = logicsinfo.Client_phone
		}

		if logicsinfo.Internal_ship_company_code != "" {
			changedata["internal_ship_company_code"] = logicsinfo.Internal_ship_company_code
		}

		if logicsinfo.Internal_ship_num != "" {
			changedata["internal_ship_num"] = logicsinfo.Internal_ship_num
		}

		if logisticssinfo != nil {
			//原来就有，修改信息
			_, err := dboper.Raw(fmt.Sprintf("update %s set %s where `%s`=?", self.TableName(), db.SqlGetKeyValue(changedata, "="), "id"), logicsinfo.Logistics).Exec()
			if err != nil {
				return errors.WithStack(err)
			}
		} else {
			starttime := curtime
			if logicsinfo.Logistics_task_starttime != "" {
				// logicsinfo.Logistics_task_starttime = curtime
				var err error
				starttime, err = strconv.ParseInt(logicsinfo.Logistics_task_starttime, 10, 64)
				if err != nil {
					return errors.WithStack(err)
				}
			}

			taskstr := self.BuildTask(dboper, default_logistics_task, starttime)

			//新的
			changedata["id"] = logicsinfo.Logistics
			changedata["build_time"] = curtime
			changedata["logistics_task"] = taskstr
			changedata["logistics_task_starttime"] = starttime
			if logicsinfo.Client_phone != "" && logicsinfo.Customer_name != "" && (logicsinfo.User_id_number == "" || logicsinfo.Idnumpic1 == "" || logicsinfo.Idnumpic2 == "") {
				//填入用户之前的数据
				olddatalist, err := self.GetInfoByWhere(dboper, fmt.Sprintf("`client_phone`=%s and `client_name`=%s ", db.SqlGetString(logicsinfo.Client_phone), db.SqlGetString(logicsinfo.Customer_name)))
				if err != nil {
					return errors.WithStack(err)
				}
				if olddatalist != nil {
					if logicsinfo.User_id_number == "" {
						changedata["idnum"] = olddatalist[0]["idnum"]
					}
					if logicsinfo.Idnumpic2 == "" {
						changedata["idnumpic1"] = olddatalist[0]["idnumpic1"]
					}
					if logicsinfo.Idnumpic1 == "" {
						changedata["idnumpic2"] = olddatalist[0]["idnumpic2"]
					}

				}
			}
			keys, values := db.SqlGetInsertInfo(changedata)
			_, err := dboper.Raw(fmt.Sprintf("insert into %s (%s) Values (%s)", self.TableName(), keys, values)).Exec()
			if err != nil {
				return errors.WithStack(err)
			}
			logicsinfo.Ok = true
		}
	}
	return nil
}

type LogisticsTaskItem struct {
	Use_time int    `json:"use_time"`
	Time     int64  `json:"time"`
	Info     string `json:"info"`
	Check    int    `json:"check"`
}

func (self *Logistics) UpdateTask(dboper db.DBOperIO, idarr []interface{}) ([]string, []string, error) {
	dboper.Begin()

	var errlist []string
	var successList []string
	for _, id := range idarr {
		idstr := id.(string)
		logicinfo := self.GetInfoById(dboper, idstr)
		if logicinfo == nil {
			errlist = append(errlist, "id:"+idstr+"不存在")
			continue
		}
		err, errstr := self.UpdateTaskByData(dboper, logicinfo)
		if err != nil {
			return nil, nil, errors.Wrap(err, "update ")
		}
		if errstr == "" {
			successList = append(successList, logicinfo["id"].(string))
		} else {
			errlist = append(errlist, errstr)
		}

	}
	dboper.Commit()
	return successList, errlist, nil
}

func (self *Logistics) UpdateTaskByDataList(dboper db.DBOperIO, datalist []db.Params) ([]string, []string, error) {
	dboper.Begin()
	var errlist []string
	var successList []string
	for _, logicinfo := range datalist {
		err, errstr := self.UpdateTaskByData(dboper, logicinfo)
		if err != nil {
			return nil, nil, err
		}
		if errstr == "" {
			successList = append(successList, logicinfo["id"].(string))
		} else {
			errlist = append(errlist, logicinfo["id"].(string)+errstr)
		}

	}
	dboper.Commit()
	return successList, errlist, nil
}

//更新一个
func (self *Logistics) UpdateTaskByData(dboper db.DBOperIO, logicinfo db.Params) (error, string) {

	curtime := time.Now().Unix()
	idstr := logicinfo["id"].(string)
	logistic_task_starttime, err := strconv.ParseInt(logicinfo["logistics_task_starttime"].(string), 10, 64)
	if err != nil {
		return nil, idstr + "|task_starttime:" + err.Error()
	}
	oldstate, err := strconv.Atoi(logicinfo["state"].(string))
	if err != nil {
		return nil, idstr + "|state:" + err.Error()
	}
	if logistic_task_starttime > curtime {
		return nil, idstr + "|未到开始时间"
	}
	if logicinfo["logistics_task"] == nil {
		return nil, idstr + "|没有进度"
	}
	taskliststr := logicinfo["logistics_task"].(string)
	if taskliststr == "" {
		return nil, idstr + "|没有进度"
	}

	var taskList []LogisticsTaskItem
	err = json.Unmarshal([]byte(taskliststr), &taskList)
	if err != nil {
		return nil, idstr + "|数据错误"
	}
	if len(taskList) == 0 {
		return nil, idstr + "|进度空"
	}
	var endtime = logistic_task_starttime
	var updateok = false

	for index, taskitem := range taskList {
		endtime += int64(taskitem.Use_time * 3600)
		//logs.Info("build time:%d", endtime)
		if taskList[index].Time != endtime {
			taskList[index].Time = endtime
			updateok = true
		}

	}

	var newstate = libs.ShipNotBeginValue

	var allover = true
	for index, taskitem := range taskList {
		if taskitem.Time < curtime {
			if taskList[index].Check == 0 {
				updateok = true
				taskList[index].Check = 1
			}

		} else {

			if taskList[index].Check == 1 {
				updateok = true
				taskList[index].Check = 0
			}

			allover = false
		}
	}
	if allover {
		newstate = libs.ShipOverseaOverValue
	} else {
		newstate = libs.ShiOverseaValue
	}

	if updateok == false && oldstate == newstate {
		return nil, "无更新"
	}

	tasknewstr, err := json.Marshal(taskList)
	if err != nil {

		return nil, idstr + "|格式化出错" + err.Error()
	}
	var changedata = make(map[string]interface{})
	changedata["logistics_task"] = string(tasknewstr)
	changedata["state"] = newstate

	_, err = dboper.Raw(fmt.Sprintf("update %s set %s where `id`=?", self.TableName(), db.SqlGetKeyValue(changedata, "=")), idstr).Exec()
	if err != nil {
		return err, idstr + "数据库出错" + err.Error()
	}

	return nil, ""
}

//查物流

//https://www.kuaidi100.com/query?type=yunda&postid=3903300539521&id=11&valicode=&temp=0.31565693514011617
//{"message":"ok","nu":"3903300539521","ischeck":"1","condition":"F00","com":"yunda","status":"200","state":"3",
//"data":[{"time":"2016-08-20 18:06:39","ftime":"2016-08-20 18:06:39","context":"[河南南阳公司卧龙分部]快件已被 已签收 签收","location":null},
//{"time":"2016-08-20 18:02:01","ftime":"2016-08-20 18:02:01","context":"[河南南阳公司卧龙分部]进行派件扫描；派送业务员：李师傅；联系电话：15993190833","location":null},{"time":"2016-08-20 18:01:20","ftime":"2016-08-20 18:01:20","context":"[河南南阳公司卧龙分部]到达目的地网点，快件将很快进行派送","location":null},{"time":"2016-08-18 17:58:24","ftime":"2016-08-18 17:58:24","context":"[河南南阳公司卧龙分部]进行派件扫描；派送业务员：徐奎西；联系电话：15938888811","location":null},{"time":"2016-08-18 15:09:20","ftime":"2016-08-18 15:09:20","context":"[河南南阳公司卧龙分部]进行派件扫描；派送业务员：徐奎西；联系电话：15938888811","location":null},{"time":"2016-08-18 12:55:26","ftime":"2016-08-18 12:55:26","context":"[河南南阳公司卧龙分部]到达目的地网点，快件将很快进行派送","location":null},{"time":"2016-08-18 10:32:23","ftime":"2016-08-18 10:32:23","context":"[河南南阳公司]进行快件扫描，将发往：河南南阳公司卧龙分部","location":null},{"time":"2016-08-18 03:30:05","ftime":"2016-08-18 03:30:05","context":"[河南漯河分拨中心]从站点发出，本次转运目的地：河南南阳公司","location":null},{"time":"2016-08-18 03:25:17","ftime":"2016-08-18 03:25:17","context":"[河南漯河分拨中心]在分拨中心进行卸车扫描","location":null},{"time":"2016-08-17 21:42:13","ftime":"2016-08-17 21:42:13","context":"[河南郑州分拨中心]进行装车扫描，即将发往：河南漯河分拨中心","location":null},{"time":"2016-08-17 21:39:58","ftime":"2016-08-17 21:39:58","context":"[河南郑州分拨中心]在分拨中心进行卸车扫描","location":null},{"time":"2016-08-17 05:03:25","ftime":"2016-08-17 05:03:25","context":"[浙江杭州分拨中心]进行装车扫描，即将发往：河南郑州分拨中心","location":null},{"time":"2016-08-17 04:57:02","ftime":"2016-08-17 04:57:02","context":"[浙江杭州分拨中心]在分拨中心进行称重扫描","location":null},{"time":"2016-08-16 23:00:38","ftime":"2016-08-16 23:00:38","context":"[浙江宁波分拨中心]进行装车扫描，即将发往：浙江杭州分拨中心","location":null},{"time":"2016-08-16 22:16:25","ftime":"2016-08-16 22:16:25","context":"[浙江宁波分拨中心]在分拨中心进行称重扫描","location":null},{"time":"2016-08-16 18:59:44","ftime":"2016-08-16 18:59:44","context":"[浙江宁波鄞州区潘火公司]进行揽件扫描","location":null},{"time":"2016-08-16 15:36:10","ftime":"2016-08-16 15:36:10","context":"[浙江宁波鄞州区潘火公司]进行下级地点扫描，将发往：河南漯河分拨中心","location":null},{"time":"2016-08-16 13:53:49","ftime":"2016-08-16 13:53:49","context":"[浙江宁波鄞州区潘火公司]进行揽件扫描","location":null}]}

type LogicstaticSendData struct {
	Time    string `json:"time"`
	Context string `json:"context"`
}

func (self *Logistics) GetLogicsInfo(dboper db.DBOperIO, id interface{}) (sendmsg []LogicstaticSendData, status int, err error, idinfo map[string]interface{}) {
	info := self.GetInfoById(dboper, id)
	status = libs.ShipNotExit
	if info == nil {
		err = errors.New("单号不存在")
		return
	}
	status = libs.ShipOnWay
	var overseainfo = ""
	idinfo = make(map[string]interface{})
	if info["logistics_task"] != nil {
		overseainfo = info["logistics_task"].(string)
	}

	idinfo["idnumpic1"] = info["idnumpic1"]
	idinfo["idnumpic2"] = info["idnumpic2"]
	idinfo["idnum"] = info["idnum"]
	idinfo["state"] = info["state"]
	idinfo["id"] = id

	//加上海外信息
	var overseaOk = true
	lasttime := time.Now().Unix()

	if overseainfo != "" {
		var taskList []LogisticsTaskItem
		err = json.Unmarshal([]byte(overseainfo), &taskList)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		for _, taskitem := range taskList {
			if taskitem.Check == 1 {
				//buildtime := int64(taskitem.Use_time*3600) + logistic_task_starttime
				lasttime = taskitem.Time + 3600
				buildtimestr := libs.FormatTableTime(strconv.FormatInt(taskitem.Time, 10))
				var itemdata = LogicstaticSendData{Time: buildtimestr, Context: taskitem.Info}
				sendmsg = append(sendmsg, itemdata)

			} else {
				overseaOk = false
			}
		}
	}
	if overseaOk == true {
		company := info["internal_ship_company_code"].(string)
		shipnum := info["internal_ship_num"].(string)
		if company != "" && shipnum != "" {

			idinfo["company"] = company
			idinfo["shipnum"] = shipnum

			if overseainfo != "" {
				companyname, ok := LogisticCodeMap[company]
				if ok == false {
					err = errors.Errorf("找不到:%s", company)
					return
				}
				timestr := libs.FormatTableTime(strconv.FormatInt(lasttime, 10))
				var itemdata = LogicstaticSendData{Time: timestr, Context: fmt.Sprintf("包裹清关成功!现已转发国内快递:%s 单号:%s", companyname, shipnum)}
				sendmsg = append(sendmsg, itemdata)
			}

			/*
				customerid := beego.AppConfig.String("kuai100.customerid")
				querykey := beego.AppConfig.String("kuai100.key")
				interdata, err2 := libs.GetKuai100Info(customerid, querykey, company, shipnum)
				if err2 != nil {
					err = errors.WithStack(err2)
					return
				}
				datalen := len(interdata.Data)
				for i := datalen - 1; i >= 0; i-- {
					datainfo := interdata.Data[i]
					if i == datalen-1 {
						if overseainfo != "" {
							companyname, ok := LogisticCodeMap[company]
							if ok == false {
								err = errors.Errorf("找不到:%s", company)
								return
							}
							var itemdata = LogicstaticSendData{Time: datainfo.Time, Context: fmt.Sprintf("包裹清关成功!现已转发国内快递:%s 单号:%s", companyname, shipnum)}
							sendmsg = append(sendmsg, itemdata)
						}
					}
					var tempdata = LogicstaticSendData{Time: datainfo.Time, Context: datainfo.Context}
					sendmsg = append(sendmsg, tempdata)
				}
				status, err = strconv.Atoi(interdata.State)
				if err != nil {
					logs.Info("parse state err")
					err = errors.WithStack(err)
					return
				}
			*/

		}
	}
	return
}
