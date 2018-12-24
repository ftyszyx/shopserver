package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/astaxie/beego/logs"
)

var curdb DBInfo

type Params map[string]interface{}
type ParamsList []interface{}

//一些常量
var (
	Debug            = true //配置是否打日志
	DefaultRowsLimit = 1000
	DefaultRelsDepth = 2
	DefaultTimeLoc   = time.Local
	ErrTxHasBegan    = errors.New("<Ormer.Begin> transaction already begin")
	ErrTxDone        = errors.New("<Ormer.Commit/Rollback> transaction not begin")
	ErrMultiRows     = errors.New("<QuerySeter> return multi rows")
	ErrNoRows        = errors.New("<QuerySeter> no row found")
	ErrStmtClosed    = errors.New("<QuerySeter> stmt already closed")
	ErrArgs          = errors.New("<Ormer> args error may be empty")
	ErrNotImplement  = errors.New("have not implement")
	formatTime       = "15:04:05"
	formatDate       = "2006-01-02"
	formatDateTime   = "2006-01-02 15:04:05"
	//获取行数列名
	SQLTotalName = "tbcount"
)

// SetDataBaseTZ Change the database default used timezone
func SetDataBaseTZ(name string, tz *time.Location) error {
	if al, ok := DBManger.get(name); ok {
		al.TZ = tz
	} else {
		return errors.Errorf("DataBase alias name `%s` not registered", name)
	}
	return nil
}

func addAliasWthDB(name, driverName string, db *sql.DB) (*DBInfo, error) {
	al := new(DBInfo)
	al.Name = name
	al.DriverName = driverName
	al.DB = db

	err := db.Ping()
	if err != nil {
		return nil, errors.Errorf("register db Ping `%s`, %s", name, err.Error())
	}

	if !DBManger.add(name, al) {
		return nil, errors.Errorf("DataBase alias name `%s` already registered, cannot reuse", name)
	}

	return al, nil
}

//AddAliasWithDB 添加一个数据库配置
func AddAliasWithDB(name, driverName string, db *sql.DB) error {
	_, err := addAliasWthDB(name, driverName, db)
	return errors.WithStack(err)
}

//RegisterDataBase 注册数据库
func RegisterDataBase(name, driverName, dataSource string) error {
	var (
		err  error
		db   *sql.DB
		info *DBInfo
	)
	db, err = sql.Open(driverName, dataSource)
	if err != nil {
		err = errors.Errorf("register db `%s`, %s", name, err.Error())

		goto end
	}

	info, err = addAliasWthDB(name, driverName, db)
	if err != nil {
		goto end
	}
	info.DataSource = dataSource
	info.DetectTZ()
end:
	if err != nil {
		if db != nil {
			db.Close()
		}

		logs.Error("%+v", err)
	}
	return errors.WithStack(err)
}

//GetDB 获取数据库
func GetDB(names ...string) (*DBInfo, error) {
	var name string
	if len(names) > 0 {
		name = names[0]
	} else {
		name = "default"
	}
	al, ok := DBManger.get(name)
	if ok {
		return al, nil
	}
	return nil, errors.Errorf("DataBase of alias name `%s` not found", name)
}

//新建一个操作
func NewOper() DBOperIO {
	o := new(DBOper)
	err := o.Using("default")
	if err != nil {
		panic(err)
	}
	return o
}

func LogQueies(info *DBInfo, operaton string, query string, t time.Time, err error, args ...interface{}) {
	if Debug == false {
		return
	}
	sub := time.Now().Sub(t) / 1e5
	elsp := float64(int(sub)) / 10.0
	flag := "  OK"
	if err != nil {
		flag = "FAIL"
	}
	con := fmt.Sprintf(" -[Queries/%s] - [%s / %11s / %7.1fms] - [%s]", info.Name, flag, operaton, elsp, query)
	cons := make([]string, 0, len(args))
	for _, arg := range args {
		cons = append(cons, fmt.Sprintf("%v", arg))
	}
	if len(cons) > 0 {
		con += fmt.Sprintf(" - `%s`", strings.Join(cons, "`, `"))
	}
	if err != nil {
		con += " - " + err.Error()
		logs.Error(con)
	} else {
		logs.Info(con)
	}
}
