package db

//prepare信息
import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

//数据库prepare的一项操作
// raw sql string prepared statement
type DBPreparer struct {
	DBSeter
	stmt   stmtQuerier
	closed bool
}

type DBPreparerIO interface {
	DBSeterIO
	Close() error
}

type stmtQuerier interface {
	Close() error
	Exec(args ...interface{}) (sql.Result, error)
	Query(args ...interface{}) (*sql.Rows, error)
	QueryRow(args ...interface{}) *sql.Row
}

func (self *DBPreparer) Exec() (sql.Result, error) {
	if self.closed {
		return nil, errors.WithStack(ErrStmtClosed)
	}
	return self.DBSeter.Exec()
}

func (self *DBPreparer) QueryExe() (*sql.Rows, error) {
	if self.closed {
		return nil, errors.WithStack(ErrStmtClosed)
	}
	return self.DBSeter.QueryExe()
}

// query data to []map[string]interface
func (self *DBPreparer) Values(container *[]Params, cols ...string) (int64, error) {
	return ReadValues(container, self, cols)
}

// query data to [][]interface
func (self *DBPreparer) ValuesList(container *[]ParamsList, cols ...string) (int64, error) {
	return ReadValues(container, self, cols)
}

// query data to []interface
func (self *DBPreparer) ValuesFlat(container *ParamsList, cols ...string) (int64, error) {
	return ReadValues(container, self, cols)
}

func (self *DBPreparer) Close() error {

	self.closed = true
	return self.stmt.Close()
}

func NewDBPreparer(rs *DBSeter) (DBPreparerIO, error) {
	o := new(DBPreparer)
	o.DBSeter = *rs

	query := rs.query
	begintime := time.Now()
	st, err := rs.oper.db.Prepare(query)
	LogQueies(rs.oper.info, "prepare", query, begintime, err)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	o.stmt = st
	return o, nil
}
