package db

//数据库一个操作
import (
	"database/sql"
	"time"

	"github.com/astaxie/beego/logs"

	"github.com/pkg/errors"
)

//DBOper 一次sql操作
type DBOper struct {
	info *DBInfo
	db   dbQuerier
	isTx bool
}
type DBOperIO interface {
	Raw(query string, args ...interface{}) *DBSeter
	Using(name string) error
	Begin() error
	Commit() error
	RollbackIfNotNull() error
	Rollback() error
}

//底层数据库的接口
type dbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// transaction beginner
type txer interface {
	Begin() (*sql.Tx, error)
}

// transaction ending
type txEnder interface {
	Commit() error
	Rollback() error
}

func (self *DBOper) Raw(query string, args ...interface{}) *DBSeter {

	return newDBSeter(self, query, args)
}

func (self *DBOper) Using(name string) error {
	if self.isTx {
		panic(errors.Errorf("<Ormer.Using> transaction has been start, cannot change db"))
	}
	if info, ok := DBManger.get(name); ok {
		self.info = info
		self.db = info.DB
	} else {
		return errors.Errorf("<Ormer.Using> unknown db alias name `%s`", name)
	}
	return nil
}

func (self *DBOper) Begin() error {
	if self.isTx {
		logs.Error(ErrTxHasBegan)
		return ErrTxHasBegan
	}
	var tx *sql.Tx
	begintime := time.Now()
	tx, err := self.db.(txer).Begin()
	LogQueies(self.info, "db.begin", "START TRANSACTION", begintime, err)
	if err != nil {
		logs.Error("err:%+v", err)
		return errors.WithStack(err)
	}
	self.isTx = true
	self.db = tx
	return nil
}

// commit transaction
func (self *DBOper) Commit() error {
	if !self.isTx {
		logs.Error(ErrTxDone)
		return ErrTxDone
	}
	begintime := time.Now()
	err := self.db.(txEnder).Commit()
	LogQueies(self.info, "db.commit", "COMMIT", begintime, err)
	if err == nil {
		self.isTx = false
		self.Using(self.info.Name) //返回用原来的
	} else if err == sql.ErrTxDone {
		return ErrTxDone
	}
	return errors.WithStack(err)
}

func (self *DBOper) RollbackIfNotNull() error {
	if !self.isTx {
		return nil
	}
	logs.Error("close not end begin")
	return self.Rollback()
}

func (self *DBOper) Rollback() error {
	if !self.isTx {
		logs.Error(ErrTxDone)
		return errors.WithStack(ErrTxDone)
	}
	begintime := time.Now()
	err := self.db.(txEnder).Rollback()
	LogQueies(self.info, "db.ROLLBACK", "ROLLBACK", begintime, err)
	if err == nil {
		self.isTx = false
		self.Using(self.info.Name)
		return nil
	} else if err == sql.ErrTxDone {
		logs.Error("err:%+v", err)
		return errors.WithStack(ErrTxDone)
	}
	logs.Error("err:%+v", err)
	return errors.WithStack(err)
}
