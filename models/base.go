package models

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
)

type Model struct {
	tablename string
	cache     map[string]db.Params
}

type ModelInterface interface {
	Init()
	TableName() string //表名
	InitSqlField(db.SqlType) db.SqlType
	InitJoinString(db.SqlType, bool) db.SqlType
	InitField(db.SqlType) db.SqlType

	Cache() map[string]db.Params
	ClearCache()
	ClearRowCache(string)
	GetModelStruct() interface{}

	GetFieldName(string) string
	ExportNameProcess(string, interface{}, db.Params) (string, error)
	GetInfoAndCache(db.DBOperIO, string, bool) db.Params
	CheckExit(db.DBOperIO, string, interface{}) bool
	GetInfoById(db.DBOperIO, interface{}) db.Params
	AllExcCommon(db.DBOperIO, ModelInterface, AllReqData, int) (error, int, []db.Params)
	GetInfoByField(db.DBOperIO, string, interface{}) []db.Params
	GetNumByField(db.DBOperIO, map[string]interface{}) int
	GetInfoByWhere(db.DBOperIO, string) ([]db.Params, error)
}

//获取所有时请求数据
type AllReqData struct {
	Page   int
	Rownum int
	Order  map[string]interface{}
	And    bool
	Search map[string]interface{}
}

func (self *Model) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
	return sql
}
func (self *Model) InitField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *Model) ExportNameProcess(name string, value interface{}, row db.Params) (string, error) {
	if value == nil {
		logs.Info("field %s is nil", name)
		return "", nil
	}
	datastr, ok := value.(string)
	if ok == false {
		return "", errors.New("upload file err:" + name + " not exit")
	}
	return datastr, nil

}

func (self *Model) GetFieldName(name string) string {
	return name
}

func (self *Model) ClearCache() {
	self.cache = make(map[string]db.Params)
}

func (self *Model) Cache() map[string]db.Params {
	return self.cache
}

func (self *Model) TableName() string {
	return self.tablename
}

func (self *Model) InitSqlField(sql db.SqlType) db.SqlType {
	return sql
}

func (self *Model) Init() {
	logs.Info("init:%s", self.tablename)
}

func (self *Model) GetModelStruct() interface{} {
	return nil
}

//检查是否存在某个数据
func (self *Model) CheckExitMap(oper db.DBOperIO, fieldinfo map[string]interface{}) bool {
	// db := orm.NewOrm()
	var dataList []db.Params
	var sqltext db.SqlType
	sqltext = &db.SqlBuild{}
	sqltext = sqltext.Name(self.TableName())
	num, err := oper.Raw(sqltext.Where(fieldinfo).Find()).Values(&dataList)
	if err == nil && num > 0 {
		return true
	}
	return false
}

//检查是否存在
func (self *Model) CheckExit(oper db.DBOperIO, field string, value interface{}) bool {
	data := make(map[string]interface{})
	data[field] = value
	return self.CheckExitMap(oper, data)
}

//获取表里面的一项，默认从内存取，如果内存没有，就从数据库取，并缓存。
func (self *Model) GetInfoAndCache(oper db.DBOperIO, uid string, forceUpdate bool) db.Params {
	if forceUpdate == false {
		//读旧的
		info, ok := self.cache[uid]
		if ok {
			// logs.Info("old info")
			return info
		}
	}
	// logs.Info("find info")
	// o := orm.NewOrm()
	var dataList []db.Params
	num, err := oper.Raw(fmt.Sprintf(`select * from %s where id=?`, self.TableName()), uid).Values(&dataList)
	if err == nil && num > 0 {
		self.cache[uid] = dataList[0] //添加
		// logs.Info("add info")
		return dataList[0]
	}
	return nil
}

func (self *Model) GetInfoById(oper db.DBOperIO, id interface{}) db.Params {
	res := self.GetInfoByField(oper, "id", id)
	if res != nil {
		return res[0]
	}
	return nil
}

func (self *Model) GetInfoByField(oper db.DBOperIO, field string, value interface{}) []db.Params {
	// o := orm.NewOrm()
	var dataList []db.Params
	num, err := oper.Raw(fmt.Sprintf("select * from %s where `%s`=?", self.TableName(), field), value).Values(&dataList)
	if err == nil && num > 0 {
		return dataList
	}
	if err != nil {
		logs.Error("err:%s", err.Error())
	}

	return nil
}

//获取数量
func (self *Model) GetNumByField(oper db.DBOperIO, search map[string]interface{}) int {
	// o := orm.NewOrm()
	totalnum := 0
	var dataList []db.Params
	var sqltext db.SqlType
	sqltext = &db.SqlBuild{}
	sqltext = sqltext.Name(self.TableName())
	num, err := oper.Raw(sqltext.Where(search).Count()).Values(&dataList)
	if err == nil && num > 0 {
		totalnum, err = strconv.Atoi(dataList[0][db.SQLTotalName].(string))
		if err == nil {
			return totalnum
		}
	}
	if err != nil {
		logs.Error("err:%s", err.Error())
	}

	return 0
}

func (self *Model) GetInfoByWhere(oper db.DBOperIO, where string) ([]db.Params, error) {
	// o := orm.NewOrm()
	var dataList []db.Params
	num, err := oper.Raw(fmt.Sprintf("select * from %s where %s", self.TableName(), where)).Values(&dataList)
	if err == nil && num > 0 {
		return dataList, nil
	}
	if err != nil {

		return nil, errors.WithStack(err)
	}

	return nil, nil
}

//清除缓存
func (self *Model) ClearRowCache(id string) {

	delete(self.cache, id)
}

//添加日志
func AddLog(oper db.DBOperIO, uid string, info string, control string, method string) {
	AddLogLink(oper, uid, info, control, method, "")
}

func AddLogLink(oper db.DBOperIO, uid string, info string, control string, method string, link string) {
	// db := orm.NewOrm()
	logmodel := GetModel(LOG)
	curtime := time.Now().Unix()
	_, err := oper.Raw(fmt.Sprintf("insert into %s (`userid`,`time`,`info`,`controller`,`method`,`link`) Values (?,?,?,?,?,?)", logmodel.TableName()), uid, curtime, info, control, method, link).Exec()
	if err != nil {
		// return errors.Wrap(err, "addlog failed")
		logs.Error("write log error:%+v uid:%s control:%s method:%s", err, uid, control, method)
	}
}

func (self *Model) AllExcCommon(oper db.DBOperIO, model ModelInterface, data AllReqData, gettype int) (error, int, []db.Params) {

	var totalnum = 0
	var dataList []db.Params
	var sqltext db.SqlType
	sqltext = &db.SqlBuild{}
	sqltext = sqltext.Name(model.TableName())

	if data.And {
		sqltext = sqltext.Where(data.Search)
	} else {
		sqltext = sqltext.WhereOr(data.Search)
	}
	sqltext = model.InitJoinString(sqltext, false)
	num, err := oper.Raw(sqltext.Count()).Values(&dataList)
	if err == nil && num > 0 {
		totalnum, err = strconv.Atoi(dataList[0][db.SQLTotalName].(string))
		if err != nil {
			logs.Error("err:%+v statck:\n %s", err, string(debug.Stack()))
			return errors.WithStack(err), 0, nil
		}
		if gettype == libs.GetAll_type_num {
			return nil, totalnum, nil
		}
		sqltext = sqltext.Order(data.Order)
		sqltext = model.InitJoinString(sqltext, false)
		if data.Page == 0 {
			//不用分页
			sqltext = model.InitJoinString(model.InitField(sqltext), true)

			num, err = oper.Raw(sqltext.Select()).Values(&dataList)
			if err == nil {
				return nil, totalnum, dataList
			} else {
				//logs.Error("err:%s statck:\n %s", err.Error(), string(debug.Stack()))
				return errors.WithStack(err), 0, nil
			}
		} else {
			//用分页
			var start = (data.Page - 1) * data.Rownum
			if totalnum > 1000 {
				//总数很多
				tablealias := sqltext.GetAlias()
				selfidname := "id"
				if tablealias != "" {
					selfidname = tablealias + ".id"
				}
				subsql := sqltext.Limit([]int{start, data.Rownum}).Field(map[string]string{selfidname: "id"}).Select()
				var newsqltext db.SqlType
				newsqltext = &db.SqlBuild{}
				newsqltext = newsqltext.Name(model.TableName()).Order(data.Order)
				// newsqltext = newsqltext.Name(model.TableName())
				newsqltext = model.InitJoinString(model.InitField(newsqltext), true)
				oldjoinstr := newsqltext.GetJoinStr()
				newsqltext.Join(oldjoinstr + fmt.Sprintf(" INNER join (%s) a ON `a`.`id`=%s ", subsql, db.SqlGetKey(selfidname)))
				num, err = oper.Raw(newsqltext.Select()).Values(&dataList)
				if err == nil {
					return nil, totalnum, dataList
				} else {
					//logs.Error("err:%s statck:\n %s", err.Error(), string(debug.Stack()))
					return errors.WithStack(err), 0, nil
				}

			} else {
				//按正常方式
				sqltext = model.InitJoinString(model.InitField(sqltext), true)
				num, err = oper.Raw(sqltext.Limit([]int{start, data.Rownum}).Select()).Values(&dataList)
				if err == nil {
					return nil, totalnum, dataList
				} else {
					//logs.Error("err:%s statck:\n %s", err.Error(), string(debug.Stack()))

					return errors.WithStack(err), 0, nil
				}
			}

		}
	} else {
		return errors.WithStack(err), 0, nil

	}
}
