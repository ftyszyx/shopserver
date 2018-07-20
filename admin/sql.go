package admin

import (
	"github.com/astaxie/beego/orm"
)

type SqlIO interface {
	BeforeSql(data map[string]interface{})
	AfterSql(data map[string]interface{}, oldinfo orm.Params)
}
