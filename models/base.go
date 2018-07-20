package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

type ModelInterface interface {
	Init()
	TableName() string //表名
	InitSqlField(libs.SqlType) libs.SqlType
	InitJoinString(libs.SqlType, bool) libs.SqlType
	InitField(libs.SqlType) libs.SqlType
	GetInfoAndCache(string, bool) orm.Params
	Cache() map[string]orm.Params
	ClearCache()
	ClearRowCache(string)
	CheckExit(string, interface{}) bool
	GetInfoById(interface{}) orm.Params
	GetModelStruct() interface{}
	GetInfoByField(string, interface{}) []orm.Params
	GetNumByField(map[string]interface{}) int
	GetInfoByWhere(string) []orm.Params
	GetFieldName(string) string
	ExportNameProcess(string, string) string
}

type Model struct {
	tablename string
	cache     map[string]orm.Params
}

func (self *Model) InitJoinString(sql libs.SqlType, allfield bool) libs.SqlType {
	return sql
}
func (self *Model) InitField(sql libs.SqlType) libs.SqlType {
	return sql
}

func (self *Model) ExportNameProcess(name string, value string) string {
	return value
}

func (self *Model) GetFieldName(name string) string {
	return name
}

func (self *Model) ClearCache() {
	self.cache = make(map[string]orm.Params)
}

func (self *Model) Cache() map[string]orm.Params {
	return self.cache
}

func (self *Model) TableName() string {
	return self.tablename
}

func (self *Model) InitSqlField(sql libs.SqlType) libs.SqlType {
	return sql
}

func (self *Model) Init() {
	logs.Info("init:%s", self.tablename)
}

func (self *Model) GetModelStruct() interface{} {
	return nil
}

//检查是否存在某个数据
func (self *Model) CheckExitMap(fieldinfo map[string]interface{}) bool {
	db := orm.NewOrm()
	var dataList []orm.Params
	var sqltext libs.SqlType
	sqltext = &libs.SqlBuild{}
	sqltext = sqltext.Name(self.TableName())
	num, err := db.Raw(sqltext.Where(fieldinfo).Find()).Values(&dataList)
	if err == nil && num > 0 {
		return true
	}
	return false
}

//检查是否存在
func (self *Model) CheckExit(field string, value interface{}) bool {
	data := make(map[string]interface{})
	data[field] = value
	return self.CheckExitMap(data)
}

//获取表里面的一项，默认从内存取，如果内存没有，就从数据库取，并缓存。
func (self *Model) GetInfoAndCache(uid string, forceUpdate bool) orm.Params {
	if forceUpdate == false {
		//读旧的
		info, ok := self.cache[uid]
		if ok {
			// logs.Info("old info")
			return info
		}
	}
	// logs.Info("find info")
	o := orm.NewOrm()
	var dataList []orm.Params
	num, err := o.Raw(fmt.Sprintf(`select * from %s where id=?`, self.TableName()), uid).Values(&dataList)
	if err == nil && num > 0 {
		self.cache[uid] = dataList[0] //添加
		// logs.Info("add info")
		return dataList[0]
	}
	return nil
}

func (self *Model) GetInfoById(id interface{}) orm.Params {
	res := self.GetInfoByField("id", id)
	if res != nil {
		return res[0]
	}
	return nil
}

func (self *Model) GetInfoByField(field string, value interface{}) []orm.Params {
	o := orm.NewOrm()
	var dataList []orm.Params
	num, err := o.Raw(fmt.Sprintf("select * from %s where `%s`=?", self.TableName(), field), value).Values(&dataList)
	if err == nil && num > 0 {
		return dataList
	}
	if err != nil {
		logs.Error("err:%s", err.Error())
	}

	return nil
}

//获取数量
func (self *Model) GetNumByField(search map[string]interface{}) int {
	o := orm.NewOrm()
	totalnum := 0
	var dataList []orm.Params
	var sqltext libs.SqlType
	sqltext = &libs.SqlBuild{}
	sqltext = sqltext.Name(self.TableName())
	num, err := o.Raw(sqltext.Where(search).Count()).Values(&dataList)
	if err == nil && num > 0 {
		totalnum, err = strconv.Atoi(dataList[0][libs.SQL_COUNT_NAME].(string))
		if err == nil {
			return totalnum
		}
	}
	if err != nil {
		logs.Error("err:%s", err.Error())
	}

	return 0
}

func (self *Model) GetInfoByWhere(where string) []orm.Params {
	o := orm.NewOrm()
	var dataList []orm.Params
	num, err := o.Raw(fmt.Sprintf("select * from %s where %s", self.TableName(), where)).Values(&dataList)
	if err == nil && num > 0 {
		return dataList
	}
	if err != nil {
		logs.Error("err:%s", err.Error())
	}

	return nil
}

//清除缓存
func (self *Model) ClearRowCache(id string) {

	delete(self.cache, id)
}

//添加日志
func AddLog(uid string, info string, control string, method string) {
	AddLogLink(uid, info, control, method, "")
}

func AddLogLink(uid string, info string, control string, method string, link string) {
	db := orm.NewOrm()
	logmodel := GetModel(LOG)
	curtime := time.Now().Unix()
	_, err := db.Raw(fmt.Sprintf("insert into %s (`userid`,`time`,`info`,`controller`,`method`,`link`) Values (?,?,?,?,?,?)", logmodel.TableName()), uid, curtime, info, control, method, link).Exec()
	if err != nil {
		logs.Error("write log error:%s uid:%s control:%s method:%s", err.Error(), uid, control, method)
	}
}
