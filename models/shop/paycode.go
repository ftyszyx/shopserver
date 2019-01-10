package shop

//支付码
import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type PayCode struct {
	models.Model
}

type PayCodeData struct {
}

func (self *PayCode) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))

}

func (self *PayCode) GetModelStruct() interface{} {
	return PayCodeData{}
}

func (self *PayCode) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {

	usertable := models.GetModel(names.USER).TableName()
	fieldstr := ""
	if (allfield == true) || (sql.NeedJointable("user") == true) {
		fieldstr += fmt.Sprintf("left join `%s` `user` ON `user`.`id`=`paycode`.`user_id`", usertable)
	}

	return sql.Alias("paycode").Join(fieldstr)
}

func (self *PayCode) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"paycode.id":         "id",
		"paycode.order_list": "order_list",
		"paycode.money":      "money",
		"paycode.user_id":    "user_id",
		"paycode.build_time": "build_time",
		"user.name":          "user_name",
		"user.account":       "user_account",
	})
}

type payCodeinfoData struct {
	lock     sync.Mutex
	lasttime int64
	Num      int
}

var paycodeinfo payCodeinfoData

func updatePaycodeTime(timecount int64) {
	paycodeinfo.lock.Lock()
	if paycodeinfo.lasttime == timecount {
		paycodeinfo.Num++
	} else {
		paycodeinfo.Num = 1
	}
	paycodeinfo.lasttime = timecount
	paycodeinfo.lock.Unlock()
}

var alphaNum = []byte(`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`)

func getPayId() string {
	nowtime := time.Now()
	updatePaycodeTime(nowtime.Unix())
	var ordernum int64
	ordernum = nowtime.Unix() * int64(paycodeinfo.Num)
	var getstr []byte
	for ordernum >= 36 {
		remain := ordernum % 36
		ordernum = ordernum / 36
		getstr = append(getstr, alphaNum[remain])
	}
	getstr = append(getstr, alphaNum[ordernum])
	return string(getstr)
}

//增加支付订单号
func (self *PayCode) AddPayCodeByOrderList(dboper db.DBOperIO, orderlist []string, money float64, uid string) (string, error) {
	orderliststr, err := json.Marshal(orderlist)
	if err != nil {
		return "", err
	}
	liststr := string(orderliststr)

	var payid = ""
	oldinfo := self.GetInfoByField(dboper, "order_list", liststr)
	if oldinfo == nil {
		payid = getPayId()
		err = self.AddPayCode(dboper, string(orderliststr), money, uid, payid)
		if err != nil {
			return "", err
		}
	} else {
		if len(oldinfo) > 1 {
			return "", errors.New("有重复单号")
		} else {
			payid = oldinfo[0]["id"].(string)
			oldmoney := oldinfo[0]["money"].(string)
			olduid := oldinfo[0]["user_id"].(string)

			oldmoneyvalue, err := strconv.ParseFloat(oldmoney, 64)
			if err != nil {
				return "", err
			}

			if olduid != uid || oldmoneyvalue != money {
				models.AddLog(dboper, uid, fmt.Sprintf("pyacode 冲突:old:%+v uid:%s money:%f", oldinfo[0], uid, money), "error", "AddPayCodeByOrderList")
				changedata := make(map[string]interface{})
				changedata["user_id"] = uid
				changedata["money"] = money
				changestr := db.SqlGetKeyValue(changedata, "=")
				_, err := dboper.Raw(fmt.Sprintf("update %s set %s where `id`=?", self.TableName(), changestr), payid).Exec()
				if err != nil {
					return "", err
				}
			}
			return payid, nil
		}
	}
	return payid, nil
}

func (self *PayCode) AddPayCodeByOrder(dboper db.DBOperIO, order string, money float64, uid string) (string, error) {
	var orderlist = []string{order}
	return self.AddPayCodeByOrderList(dboper, orderlist, money, uid)
}

func (self *PayCode) AddPayCode(dboper db.DBOperIO, order string, money float64, uid string, payid string) error {
	curtime := time.Now().Unix()
	_, err := dboper.Raw(fmt.Sprintf("insert into %s (`id`,`order_list`,`money`,`user_id`,`build_time`) Values (?,?,?,?,?)", self.TableName()), payid, order, money, uid, curtime).Exec()
	if err != nil {
		return err
	}
	return nil
}
