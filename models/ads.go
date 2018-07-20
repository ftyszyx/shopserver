package models

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/zyx/shop_server/libs"
)

type Ads struct {
	Model
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

func (self *Ads) InitSqlField(sql libs.SqlType) libs.SqlType {
	return self.InitField(self.InitJoinString(sql, true))
}

func (self *Ads) InitJoinString(sql libs.SqlType, allfield bool) libs.SqlType {
	itemtablename := GetModel(SHOP_ITEM).TableName()
	posttablename := GetModel(POST).TableName()
	adspos := GetModel(ADSPOS).TableName()

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

func (self *Ads) InitField(sql libs.SqlType) libs.SqlType {
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

func (self *Ads) GetModelStruct() interface{} {
	return AdsData{}
}

var AdsHomeCache cache.Cache //主页广告缓存

func (self *Ads) Init() {
	AdsHomeCache, _ = cache.NewCache("memory", `{"interval":0}`) //不过期
	//AdsUpdateEvent.AddLister(libs.HandlerFunc(self.initHomeCache))

	self.initHomeCache()
}

func (self *Ads) initHomeCache() {
	configmodel := GetModel(CONFIG)
	adsposmodel := GetModel(ADSPOS)
	itemmodel := GetModel(SHOP_ITEM)
	postmodel := GetModel(POST)

	res := configmodel.GetInfoByField("name", "home_ads_set")
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
	adsarr := adsdata["data"].([]interface{})

	for _, adsposItem := range adsarr {
		adsposItem := adsposItem.(map[string]interface{})
		posid := adsposItem["posid"].(string)

		resone := adsposmodel.GetInfoById(posid)
		if resone == nil {
			continue
		}
		adsposItem["adspos"] = resone
		res = self.GetInfoByWhere(fmt.Sprintf("`ads_pos`=%s and `is_del`=0 ", posid))
		if res != nil {
			adsposItem["ads"] = res
			for _, adsitem := range res {
				if adsitem["item_id"] != nil {
					iteminfo := itemmodel.GetInfoById(adsitem["item_id"].(string))
					if iteminfo != nil {
						adsitem["iteminfo"] = iteminfo
					}
				}
				if adsitem["post_id"] != nil {
					postinfo := postmodel.GetInfoById(adsitem["post_id"].(string))
					if postinfo != nil {
						adsitem["postinfo"] = postinfo
					}
				}
			}
		}

	}
	AdsHomeCache.Put("homeads", adsarr, 0)
	self.initAdsPosCatch("adspos.name", "swipe", false)
	self.initAdsPosCatch("adspos.name", "swipecopany", false)
	self.initAdsPosCatch("adspos.name", "cases", false)
	self.initAdsPosCatch("adspos.name", "news", false)
	self.initAdsPosCatch("adspos.name", "product1", false)
	self.initAdsPosCatch("adspos.name", "product2", false)
	self.initAdsPosCatch("ads.name", "about", true)
	self.initAdsPosCatch("ads.name", "contact", true)
	self.initAdsPosCatch("ads.name", "joinus", true)
}

func (self *Ads) initAdsPosCatch(fieldname string, name string, onlyone bool) {
	o := orm.NewOrm()
	var sqltext libs.SqlType
	var dataList []orm.Params
	sqltext = &libs.SqlBuild{}
	sqltext = sqltext.Name(self.tablename)
	sqltext = self.InitSqlField(sqltext)

	sqltext = sqltext.Where(map[string]interface{}{fieldname: name})
	num, err := o.Raw(sqltext.Select()).Values(&dataList)
	if err == nil && num > 0 {
		if onlyone {
			AdsHomeCache.Put(name, dataList[0], 0)
		} else {
			AdsHomeCache.Put(name, dataList, 0)
		}

	}
}

func (self *Ads) updateHomeCache() {
	AdsHomeCache.ClearAll()
	self.initHomeCache()
}

func (self *Ads) ClearCache() {
	self.cache = make(map[string]orm.Params)
	self.updateHomeCache()
}
