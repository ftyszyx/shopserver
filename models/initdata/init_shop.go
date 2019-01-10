package initdata

import (
	"github.com/zyx/shop_server/models"
	"github.com/zyx/shop_server/models/coredata"
	"github.com/zyx/shop_server/models/names"
	"github.com/zyx/shop_server/models/shop"
)

//模块名
func InitShopModel() {
	models.InitAllModel()
	allmodels := models.GetAllModel()
	allmodels[names.USER] = &coredata.User{models.NewModel("aq_user", true)}
	allmodels[names.USERGROUP] = &coredata.UserGroup{models.NewModel("aq_usergroup", true)}
	allmodels[names.MODULE] = &coredata.Module{models.NewModel("aq_module", true)}
	allmodels[names.LOG] = &coredata.Log{models.NewModel("aq_log", false)}
	allmodels[names.POST] = &coredata.Post{models.NewModel("aq_post", false)}
	allmodels[names.POSTTYPE] = &coredata.PostType{models.NewModel("aq_post_type", false)}
	allmodels[names.ALBUM] = &coredata.Album{models.NewModel("aq_album", false)}
	allmodels[names.PHOTO] = &coredata.Photo{models.NewModel("aq_photo", false)}
	allmodels[names.CONFIG] = &coredata.Config{models.NewModel("aq_config", true)}
	allmodels[names.ADS] = &coredata.Ads{models.NewModel("aq_ads", true)}
	allmodels[names.ADSPOS] = &coredata.AdsPos{models.NewModel("aq_ads_pos", false)}
	allmodels[names.EXPORT_TASK] = &coredata.ExportTask{models.NewModel("aq_export_task", false)}
	allmodels[names.EXPORT] = &coredata.Export{models.NewModel("aq_export", false)}
	allmodels[names.DATABASE] = &coredata.DataBase{models.NewModel("aq_database", false)}

	allmodels[names.SHOP_BRAND] = &shop.ShopBrand{models.NewModel("aq_brand", false)}
	allmodels[names.SHOP_ITEM] = &shop.ShopItem{models.NewModel("aq_item", false)}
	allmodels[names.SHOP_ITEMTYPE] = &shop.ShopItemType{models.NewModel("aq_item_type", false)}
	allmodels[names.SHOP_NOTICE] = &shop.ShopNotice{models.NewModel("aq_notice", true)}
	allmodels[names.SHOP_ORDER] = &shop.ShopOrder{models.NewModel("aq_order", false)}
	allmodels[names.SHOP_TAG] = &shop.ShopTag{models.NewModel("aq_tag", false)}
	allmodels[names.PAYCODE] = &shop.PayCode{models.NewModel("aq_paycode", false)}

	for _, value := range allmodels {
		value.Init()
	}
}
