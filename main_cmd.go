package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/admin"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
)

//数据库备份
func saveDatabaseTask() {
	dboper := db.NewOper()
	path, err := admin.SaveDatabase(dboper)
	if err != nil {
		logs.Error("back up system err:%", err.Error())
	}
	adddata := make(map[string]interface{})
	adddata["name"] = fmt.Sprintf("系统备份 %s", time.Now().Format("2006-01-02"))
	adddata["user_id"] = "1"
	adddata["build_time"] = time.Now().Unix()
	adddata["path"] = path
	// o := orm.NewOrm()
	keys, values := db.SqlGetInsertInfo(adddata)
	_, err = db.NewOper().Raw(fmt.Sprintf("insert into aq_database (%s) values (%s)", keys, values)).Exec()
	if err != nil {
		logs.Error("back up system err:%", err.Error())
	}
}

//更新所有物流进度
func updateAllShip() {
	dboper := db.NewOper()
	curtime := time.Now().Unix()
	logisticsModel := models.GetModel(models.LOGISTICS).(*models.Logistics)

	res, err := logisticsModel.GetInfoByWhere(dboper, fmt.Sprintf("`logistics_task_starttime`<%d and  `state` <> %d and `id` like ", curtime, libs.ShipOverseaOverValue)+"'AB%AU'")
	if err != nil {
		logs.Error("updateAllShip system err:%+v", err)
	}
	if res != nil {
		oklist, errlist, err := logisticsModel.UpdateTaskByDataList(dboper, res)
		if err != nil {
			logs.Error("updateAllShip system err:%+v", err)
			return
		}
		// logs.Info("oklist:%+v ", oklist)
		// logs.Info("errlist:%+v", errlist)
		models.AddLog(dboper, "1", fmt.Sprintf("更新全部物流进度,成功:%d errlist:%d", len(oklist), len(errlist)), "logistics", "UpdateAllTask")
		// logs.Info("oklist:%d errlist:%d", len(oklist), len(errlist))
	}
}

//批量给订单生成支付码（不用了)
func initpaycodeData() {
	paycodeModel := models.GetModel(models.LOGISTICS).(*models.PayCode)
	orderModel := models.GetModel(models.SHOP_ORDER).(*models.ShopOrder)

	dboper := db.NewOper()
	var dataList []db.Params
	num, err := dboper.Raw(fmt.Sprintf("select * from %s ", orderModel.TableName())).Values(&dataList)
	if err != nil {
		logs.Error("err:%s", err.Error())
		return
	}
	if err == nil && num > 0 {
		for _, orderitem := range dataList {
			payid := orderitem["pay_id"].(string)
			var orderlist = []string{orderitem["id"].(string)}
			orderliststr, err := json.Marshal(orderlist)
			if err != nil {
				panic(err.Error())
			}
			if payid != "" {
				totalprice, err := strconv.ParseFloat(orderitem["total_price"].(string), 64)
				if err != nil {
					panic(err.Error())
				}
				err = paycodeModel.AddPayCode(dboper, string(orderliststr), totalprice, orderitem["user_id"].(string), payid)
				if err != nil {
					panic(err.Error())
				}
			}
		}
	}
}

//清除系统垃圾（数据库备份，导出的表）
func CleanSystem() {
	dboper := db.NewOper()
	manger := libs.GetManger()
	host := beego.AppConfig.String("qiniu.host")
	bucket := beego.AppConfig.String("qiniu.bucket")

	DatabaseModel := models.GetModel(models.DATABASE).(*models.DataBase)
	ExportTaskModel := models.GetModel(models.EXPORT_TASK).(*models.ExportTask)

	//清除数据库备份
	var dataList []db.Params
	num, err := dboper.Raw(fmt.Sprintf("select * from %s order by `build_time` DESC", DatabaseModel.TableName())).Values(&dataList)
	if err != nil {
		logs.Error("get table system err:%+v", err)
		return
	}

	if num > 1 {
		for index, dbitem := range dataList {
			if index > 0 {
				//留下最后一条
				key := strings.TrimPrefix(dbitem["path"].(string), host)
				logs.Info("del key:%s", key)
				err := manger.Delete(bucket, key)
				if err != nil {
					logs.Error("delete qiniu err:%+v", err)
					//return
				}
				_, err = dboper.Raw(fmt.Sprintf("delete from %s where `id` = ?", DatabaseModel.TableName()), dbitem["id"]).Exec()
				if err != nil {
					logs.Error("delete err:%+v", err)
					return
				}
			}
		}
	}

	num, err = dboper.Raw(fmt.Sprintf("select * from %s  ", ExportTaskModel.TableName())).Values(&dataList)
	if err != nil {
		logs.Error("get tasklist system err:%+v", err)
		return
	}
	if num > 0 {
		for _, taskitem := range dataList {
			//留下最后一条
			key := strings.TrimPrefix(taskitem["path"].(string), host)
			logs.Info("del key:%s", key)
			err := manger.Delete(bucket, key)
			if err != nil {
				logs.Error("delete qiniu err:%+v", err)
				//return
			}
			_, err = dboper.Raw(fmt.Sprintf("delete from %s where `id` = ?", ExportTaskModel.TableName()), taskitem["id"]).Exec()
			if err != nil {
				logs.Error("delete err:%+v", err)
				return
			}
		}
	}

}
