package base

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
)

func (self *BaseController) AddOneRow(rownum int, rowinfo []string) string {
	return ""
}

//表格上传
func (self *BaseController) UploadeCSV(sqlcall SqlIO) error {
	logs.Info("UploadeCSV")
	err, fileinfo := self.uploadCommon(false)
	if err != nil {
		return err
	}
	logs.Info("UploadeCSVafter:%+v", fileinfo)
	filetype := fileinfo["filetype"].(string)
	if filetype != "csv" {

		return errors.New("表格格式错误，只支持CSV")
	}
	fielpath := fileinfo["filePath"].(string)

	fileio, err := os.Open(fielpath)

	if err != nil {
		self.AjaxReturnError(errors.WithStack(err))
	}
	defer func() {
		fileio.Close()
		os.Remove(fielpath)
	}()

	reader := csv.NewReader(fileio)

	_, err = reader.Read()
	if err != nil {
		return err
	}
	rownum := 0
	self.dboper.Begin()
	for {
		//读每一行
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			self.dboper.Rollback()
			return err
		}
		rownum++
		if rownum > 9000 {
			self.dboper.Rollback()

			return errors.New("超过最大行数9000")
		}
		errstr := sqlcall.AddOneRow(rownum, record)
		if errstr != "" {
			self.dboper.Rollback()
			return errors.New(errstr)
		}
	}
	self.dboper.Commit()
	return nil
}

func Getcolstr(col int, rowinfo []string) (int, string) {

	utf8byte, err := libs.GbkToUtf8([]byte(rowinfo[col]))
	if err != nil {
		logs.Info(err.Error())
	}
	return col + 1, strings.TrimSpace(string(utf8byte))
}

func GetImportErr(col int, row int, msg string) string {
	return fmt.Sprintf("第%d行 第%d列 错误:%s", row+1, col, msg)
}
