package db

//保存的参数信息
import (
	"database/sql"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

type DBSeterIO interface {
	SetArgs(...interface{}) DBSeterIO

	//execute sql and get result
	Exec() (sql.Result, error)
	QueryExe() (*sql.Rows, error)

	Prepare() (DBPreparerIO, error)

	// query data to []map[string]interface
	// see QuerySeter's Values
	Values(container *[]Params, cols ...string) (int64, error)
	// query data to [][]interface
	// see QuerySeter's ValuesList
	ValuesList(container *[]ParamsList, cols ...string) (int64, error)
	// query data to []interface
	// see QuerySeter's ValuesFlat
	ValuesFlat(container *ParamsList, cols ...string) (int64, error)
}

//数据库的一项操作
type DBSeter struct {
	query string
	args  []interface{}
	oper  *DBOper
}

func newDBSeter(oper *DBOper, query string, args []interface{}) *DBSeter {
	o := new(DBSeter)
	o.query = query
	o.args = args
	o.oper = oper
	return o
}

// set args for every query
func (self *DBSeter) SetArgs(args ...interface{}) DBSeterIO {
	self.args = args
	return self
}

// execute raw sql and return sql.Result
func (self *DBSeter) Exec() (sql.Result, error) {
	a := time.Now()
	res, err := self.oper.db.Exec(self.query, self.args...)
	LogQueies(self.oper.info, "Exec", self.query, a, err, self.args...)
	return res, errors.WithStack(err)
}

func (self *DBSeter) QueryExe() (*sql.Rows, error) {
	begintime := time.Now()
	res, err := self.oper.db.Query(self.query, self.args...)
	LogQueies(self.oper.info, "Query", self.query, begintime, err, self.args)
	return res, errors.WithStack(err)
}

func ReadValues(container interface{}, dbexe DBSeterIO, needCols []string) (int64, error) {
	var (
		maps  []Params
		lists []ParamsList
		list  ParamsList
	)

	typ := 0
	switch container.(type) {
	case *[]Params:
		typ = 1
	case *[]ParamsList:
		typ = 2
	case *ParamsList:
		typ = 3
	default:
		panic(errors.Errorf("<RawSeter> unsupport read values type `%T`", container))
	}

	//args := getFlatParams(nil, self.args, self.oper.info.TZ)

	var rs *sql.Rows
	rs, err := dbexe.QueryExe()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer rs.Close()

	var (
		refs   []interface{}
		cnt    int64
		cols   []string
		indexs []int
	)

	for rs.Next() {
		if cnt == 0 {
			columns, err := rs.Columns()
			if err != nil {
				return 0, errors.WithStack(err)
			}
			if len(needCols) > 0 {
				indexs = make([]int, 0, len(needCols))
			} else {
				indexs = make([]int, 0, len(columns))
			}

			cols = columns
			refs = make([]interface{}, len(cols))
			for i := range refs {
				var ref sql.NullString
				refs[i] = &ref

				if len(needCols) > 0 {
					for _, c := range needCols {
						if c == cols[i] {
							indexs = append(indexs, i)
						}
					}
				} else {
					indexs = append(indexs, i)
				}
			}
		}

		if err := rs.Scan(refs...); err != nil {
			return 0, errors.WithStack(err)
		}

		switch typ {
		case 1:
			params := make(Params, len(cols))
			for _, i := range indexs {
				ref := refs[i]
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					params[cols[i]] = value.String
				} else {
					params[cols[i]] = nil
				}
			}
			maps = append(maps, params)
		case 2:
			params := make(ParamsList, 0, len(cols))
			for _, i := range indexs {
				ref := refs[i]
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					params = append(params, value.String)
				} else {
					params = append(params, nil)
				}
			}
			lists = append(lists, params)
		case 3:
			for _, i := range indexs {
				ref := refs[i]
				value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
				if value.Valid {
					list = append(list, value.String)
				} else {
					list = append(list, nil)
				}
			}
		}

		cnt++
	}

	switch v := container.(type) {
	case *[]Params:
		*v = maps
	case *[]ParamsList:
		*v = lists
	case *ParamsList:
		*v = list
	}

	return cnt, nil
}

// query data to []map[string]interface
func (self *DBSeter) Values(container *[]Params, cols ...string) (int64, error) {
	return ReadValues(container, self, cols)
}

// query data to [][]interface
func (self *DBSeter) ValuesList(container *[]ParamsList, cols ...string) (int64, error) {
	return ReadValues(container, self, cols)
}

// query data to []interface
func (self *DBSeter) ValuesFlat(container *ParamsList, cols ...string) (int64, error) {
	return ReadValues(container, self, cols)
}

// return prepared raw statement for used in times.
func (self *DBSeter) Prepare() (DBPreparerIO, error) {
	return NewDBPreparer(self)
}
