package db

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/astaxie/beego/logs"
)

//DriverType 数据库类型
type DriverType int

const (
	_          DriverType = iota // int enum type
	DRMySQL                      // mysql
	DRSqlite                     // sqlite
	DROracle                     // oracle
	DRPostgres                   // pgsql
	DRTiDB                       // TiDB
)

type DBInfo struct {
	Name         string
	Driver       DriverType
	DriverName   string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	DB           *sql.DB
	TZ           *time.Location //地区
	Engine       string         //引擎
}

type DBInfoIO interface {
	SetMaxIdleConns(maxIdleConns int)
	SetMaxOpenConns(maxOpenConns int)
	DetectTZ()
}

// SetMaxIdleConns Change the max idle conns for *sql.DB, use specify database alias name
func (self *DBInfo) SetMaxIdleConns(maxIdleConns int) {

	self.MaxIdleConns = maxIdleConns
	self.DB.SetMaxIdleConns(maxIdleConns)
}

// SetMaxOpenConns Change the max open conns for *sql.DB, use specify database alias name
func (self *DBInfo) SetMaxOpenConns(maxOpenConns int) {

	self.MaxOpenConns = maxOpenConns
	// for tip go 1.2
	if fun := reflect.ValueOf(self.DB).MethodByName("SetMaxOpenConns"); fun.IsValid() {
		fun.Call([]reflect.Value{reflect.ValueOf(maxOpenConns)})
	}
}

func (self *DBInfo) DetectTZ() {
	// orm timezone system match database
	// default use Local
	self.TZ = DefaultTimeLoc
	switch self.Driver {
	case DRMySQL:
		row := self.DB.QueryRow("SELECT TIMEDIFF(NOW(), UTC_TIMESTAMP)")
		var tz string
		row.Scan(&tz)
		if len(tz) >= 8 {
			if tz[0] != '-' {
				tz = "+" + tz
			}
			t, err := time.Parse("-07:00:00", tz)
			if err == nil {
				if t.Location().String() != "" {
					self.TZ = t.Location()
				}
			} else {
				logs.Error("Detect DB timezone: %s %+v\n", tz, err)
			}
		}

		// get default engine from current database
		row = self.DB.QueryRow("SELECT ENGINE, TRANSACTIONS FROM information_schema.engines WHERE SUPPORT = 'DEFAULT'")
		var engine string
		var tx bool
		row.Scan(&engine, &tx)

		if engine != "" {
			self.Engine = engine
		} else {
			self.Engine = "INNODB"
		}
	}
}
