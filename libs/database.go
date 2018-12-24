package libs

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/zyx/shop_server/libs/db"
)

var SqlLineEnd = "\r\n"

//获取所有的表
func GetTableList(oper db.DBOperIO) ([]db.Params, error) {

	var dataList []db.Params
	// db := orm.NewOrm()
	_, err := oper.Raw("show table status").Values(&dataList)
	if err == nil {
		return dataList, nil
	}
	return nil, err
}

//获取表的字符串
func GetTableString(tablename string, oper db.DBOperIO) (string, error) {
	var buffetstr bytes.Buffer
	buffetstr.WriteString(fmt.Sprintf("drop table if exists `%s`;%s", tablename, SqlLineEnd))
	var dataList []db.Params
	// db := orm.NewOrm()
	num, err := oper.Raw(fmt.Sprintf("show create table %s", tablename)).Values(&dataList)
	if err == nil {
		if num == 0 {
			return "", errors.New(fmt.Sprintf("%s is empty", tablename))
		}
		tablestr := dataList[0]["Create Table"].(string)
		tablestr = strings.Replace(tablestr, "\n", "\r\n", -1)
		return tablestr + ";" + SqlLineEnd, nil
	}
	return "", err
}

func GetInsertSql(tablename string, start int, size int, oper db.DBOperIO) (string, error) {
	// var sqlstr string
	var dataList []db.ParamsList
	var outstr bytes.Buffer
	// db := orm.NewOrm()
	num, err := oper.Raw(fmt.Sprintf("select * from %s limit %d,%d", tablename, start, size)).ValuesList(&dataList)
	if err == nil {
		if num > 0 {
			var rowstemp bytes.Buffer
			for _, row := range dataList {
				rowstemp.Reset()
				rowstemp.WriteString("(")
				for _, value := range row {
					if value == nil {
						rowstemp.WriteString("null,")
					} else {
						rowstemp.WriteString("'")
						valuestr := value.(string)
						replacer := strings.NewReplacer(`\`, `\\`, `'`, `\'`, `"`, `\"`, "\n", "\\n", "\r", "\\r")
						valuestr = replacer.Replace(valuestr)
						rowstemp.WriteString(valuestr)
						// logs.Info("write:%s", valuestr)
						rowstemp.WriteString("',")
					}
				}
				rowstemp.Truncate((rowstemp.Len() - 1))
				rowstemp.WriteString(")")
				// rowstemp.WriteString(SqlLineEnd)
				outstr.WriteString(fmt.Sprintf("INSERT INTO `%s` VALUES %s; %s", tablename, rowstemp.String(), SqlLineEnd))
			}
			// taillen := len(SqlLineEnd) + 1
			// rowstemp.Truncate((rowstemp.Len() - taillen)) //去掉最后一个回车和逗号
			// logs.Info("row:%s", rowstemp.String())
			// outstr.WriteString(fmt.Sprintf("INSERT INTO `%s` VALUES %s %s; %s", tablename, SqlLineEnd, rowstemp.String(), SqlLineEnd))
			return outstr.String(), nil
		}
		return "", nil
	}
	return "", err
}
