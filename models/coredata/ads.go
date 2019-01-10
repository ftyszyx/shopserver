package coredata

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/libs"
	"github.com/zyx/shop_server/libs/db"
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/names"
)

type Ads struct {
	models.Model
}

type AdsData struct {
	Name     string `empty:"名称不能为空"`
	Ads_pos  string `empty:"广告位不能为空"`
	Link     string
	Pic      string
	item_id  string
	Post_id  string
	Title    string
	Order_id int
}

var AdsUpdateEvent = libs.NewEvent()

func (self *Ads) InitSqlField(sql db.SqlType) db.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *Ads) InitJoinString(sql db.SqlType, allfield bool) db.SqlType {
	itemtablename := models.GetModel(names.SHOP_ITEM).TableName()
	posttablename := models.GetModel(names.POST).TableName()
	adspos := models.GetModel(names.ADSPOS).TableName()

	joinstring := ""
	if (allfield == true) || (sql.NeedJointable("adspos") == true) {

		joinstring += fmt.Sprintf("left join `%s` `adspos` ON `ads`.`ads_pos`=`adspos`.`id`", adspos)
	}
	if (allfield == true) || (sql.NeedJointable("post") == true) {

		joinstring += fmt.Sprintf("left join `%s` `post` ON `ads`.`post_id`=`post`.`id`", posttablename)
	}
	if (allfield == true) || (sql.NeedJointable("item") == true) {

		joinstring += fmt.Sprintf("left join `%s` `item` ON `ads`.`item_id`=`item`.`id`", itemtablename)
	}
	return sql.Alias("ads").Join(joinstring)
}

func (self *Ads) InitField(sql db.SqlType) db.SqlType {
	return sql.Field(map[string]string{
		"item.name":       "item_name",
		"ads.item_id":     "item_id",
		"ads.post_id":     "post_id",
		"post.title":      "post_title",
		"post.pic":        "post_pic",
		"post.build_time": "post_build_time",
		"post.summary":    "post_summary",
		"ads.id":          "id",
		"ads.name":        "name",
		"ads.title":       "title",
		"ads.ads_pos":     "ads_pos",
		"adspos.name":     "ads_pos_name",
		"adspos.title":    "ads_pos_title",
		"ads.link":        "link",
		"ads.pic":         "pic",
		"ads.order_id":    "order_id",
		"ads.is_del":      "is_del",
	})
}

func (self *Ads) ModelStruct() interface{} {
	return AdsData{}
}

func (self *Ads) Init() {
	self.resetCache()
}

func (self *Ads) resetCache() {
	dboper := db.NewOper()
	appname := beego.AppConfig.String("appname")
	self.Cache().ClearAll()
	if appname == "shop" {
		self.initShopAds(dboper)
		self.initAdsPosCatch(dboper, "adspos.name", "swipe", false)
	} else if appname == "ship" {
		//物流
		newstype := beego.AppConfig.String("newsposttype")
		var sqltext db.SqlType
		var dataList []db.Params
		sqltext = &db.SqlBuild{}
		postmodel := models.GetModel(names.POST)
		sqltext = sqltext.Name(postmodel.TableName())
		sqltext = postmodel.InitSqlField(sqltext)

		sqltext = sqltext.Where(map[string]interface{}{"post.type": newstype, "post.is_del": 0})
		sqltext = sqltext.Order(map[string]interface{}{"post.build_time": "desc"}).Limit([]int{0, 4})

		num, err := dboper.Raw(sqltext.Select()).Values(&dataList)
		if err == nil && num > 0 {
			self.Cache().Put("shipnews", dataList, 0)
			logs.Info("get ship news num:%d", num)
		} else {
			if num == 0 {
				logs.Error("get ship news empty")
			} else {
				logs.Error("get ship news err:%s", err.Error())
			}

		}

	} else if appname == "home" {
		self.initAdsPosCatch(dboper, "adspos.name", "swipecopany", false)
		self.initAdsPosCatch(dboper, "adspos.name", "cases", false)
		self.initAdsPosCatch(dboper, "adspos.name", "news", false)
		self.initAdsPosCatch(dboper, "adspos.name", "product1", false)
		self.initAdsPosCatch(dboper, "adspos.name", "product2", false)
		self.initAdsPosCatch(dboper, "ads.name", "about", true)
		self.initAdsPosCatch(dboper, "ads.name", "contact", true)
		self.initAdsPosCatch(dboper, "ads.name", "joinus", true)
	}

}

//商城首页的配置
func (self *Ads) initShopAds(dboper db.DBOperIO) {
	configmodel := models.GetModel(names.CONFIG)
	adsposmodel := models.GetModel(names.ADSPOS)
	itemmodel := models.GetModel(names.SHOP_ITEM)
	postmodel := models.GetModel(names.POST)

	res := configmodel.GetInfoByField(dboper, "name", "home_ads_set")
	if res == nil {
		return
	}
	var adsdata map[string]interface{}
	adsvalueStr := `{"data":` + res[0]["value"].(string) + "}"

	logs.Info("value str:%s", adsvalueStr)
	err := json.Unmarshal([]byte(adsvalueStr), &adsdata)
	if err != nil {
		logs.Error("get ads err:%s", err.Error())
		return
	}
	//商城的配置
	adsarr := adsdata["data"].([]interface{})

	for _, adsposItem := range adsarr {
		adsposItem := adsposItem.(map[string]interface{})
		posid := adsposItem["posid"].(string)

		resone := adsposmodel.GetInfoById(dboper, posid)
		if resone == nil {
			continue
		}
		adsposItem["adspos"] = resone
		res, err = self.GetInfoByWhere(dboper, fmt.Sprintf("`ads_pos`=%s and `is_del`=0 ", posid))
		if err != nil {
			logs.Error("get ads err:%s", err.Error())
			return
		}
		if res != nil {
			adsposItem["ads"] = res
			for _, adsitem := range res {
				if adsitem["item_id"] != nil {
					iteminfo := itemmodel.GetInfoById(dboper, adsitem["item_id"].(string))
					if iteminfo != nil {
						adsitem["iteminfo"] = iteminfo
					}
				}
				if adsitem["post_id"] != nil {
					postinfo := postmodel.GetInfoById(dboper, adsitem["post_id"].(string))
					if postinfo != nil {
						adsitem["postinfo"] = postinfo
					}
				}
			}
		}

	}
	self.Cache().Put("homeads", adsarr, 0)
}
func (self *Ads) initAdsPosCatch(dboper db.DBOperIO, fieldname string, name string, onlyone bool) {
	// o := orm.NewOrm()
	var sqltext db.SqlType
	var dataList []db.Params
	sqltext = &db.SqlBuild{}
	sqltext = sqltext.Name(self.TableName())
	sqltext = self.InitSqlField(sqltext)

	sqltext = sqltext.Where(map[string]interface{}{fieldname: name})
	num, err := dboper.Raw(sqltext.Select()).Values(&dataList)
	if err == nil && num > 0 {
		if onlyone {
			self.Cache().Put(name, dataList[0], 0)
		} else {
			self.Cache().Put(name, dataList, 0)
		}

	}
}

func (self *Ads) ClearCache() {
	self.resetCache()
}
