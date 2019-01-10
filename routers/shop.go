package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/control/corecontrol"
	"github.com/zyx/shop_server/control/shop"
)

func initShop() {

	//商城
	beego.Router("/shop/getinfo", &shop.ShopController{}, "get:GetInfo")

	//广告
	beego.Router("/ads/all", &corecontrol.AdsController{}, "post:All")
	beego.Router("/ads/edit", &corecontrol.AdsController{}, "post:Edit")
	beego.Router("/ads/add", &corecontrol.AdsController{}, "post:Add")
	beego.Router("/ads/del", &corecontrol.AdsController{}, "post:Del")

	//广告位
	beego.Router("/adspos/all", &corecontrol.AdsPosController{}, "post:All")
	beego.Router("/adspos/edit", &corecontrol.AdsPosController{}, "post:Edit")
	beego.Router("/adspos/add", &corecontrol.AdsPosController{}, "post:Add")
	beego.Router("/adspos/del", &corecontrol.AdsPosController{}, "post:Del")

	//品牌
	beego.Router("/ShopBrand/all", &shop.ShopBrandController{}, "post:All")
	beego.Router("/ShopBrand/edit", &shop.ShopBrandController{}, "post:Edit")
	beego.Router("/ShopBrand/add", &shop.ShopBrandController{}, "post:Add")
	beego.Router("/ShopBrand/del", &shop.ShopBrandController{}, "post:Del")

	//商品
	beego.Router("/ShopItem/all", &shop.ShopItemController{}, "post:All")
	beego.Router("/ShopItem/edit", &shop.ShopItemController{}, "post:Edit")
	beego.Router("/ShopItem/add", &shop.ShopItemController{}, "post:Add")
	beego.Router("/ShopItem/del", &shop.ShopItemController{}, "post:Del")
	beego.Router("/ShopItem/ExportCsv", &shop.ShopItemController{}, "post:ExportCsv")

	//商品类型
	beego.Router("/ShopItemType/all", &shop.ShopItemTypeController{}, "post:All")
	beego.Router("/ShopItemType/edit", &shop.ShopItemTypeController{}, "post:Edit")
	beego.Router("/ShopItemType/add", &shop.ShopItemTypeController{}, "post:Add")
	beego.Router("/ShopItemType/del", &shop.ShopItemTypeController{}, "post:Del")

	//公告
	beego.Router("/ShopNotice/all", &shop.ShopNoticeController{}, "post:All")
	beego.Router("/ShopNotice/edit", &shop.ShopNoticeController{}, "post:Edit")
	beego.Router("/ShopNotice/add", &shop.ShopNoticeController{}, "post:Add")
	beego.Router("/ShopNotice/del", &shop.ShopNoticeController{}, "post:Del")

	//订单
	beego.Router("/ShopOrder/edit", &shop.ShopOrderController{}, "post:Edit")
	beego.Router("/ShopOrder/all", &shop.ShopOrderController{}, "post:All")
	beego.Router("/ShopOrder/add", &shop.ShopOrderController{}, "post:Add")
	beego.Router("/ShopOrder/orderuploade", &shop.ShopOrderController{}, "post:OrdersUpload")

	beego.Router("/ShopOrder/GetMyOrder", &shop.ShopOrderController{}, "post:GetMyOrder")
	beego.Router("/ShopOrder/ExportMyCsv", &shop.ShopOrderController{}, "post:ExportMyCsv")
	beego.Router("/ShopOrder/ClientClose", &shop.ShopOrderController{}, "post:ClientClose")
	beego.Router("/ShopOrder/ClientDelOrder", &shop.ShopOrderController{}, "post:ClientDelOrder")
	beego.Router("/ShopOrder/ClientRefundOrder", &shop.ShopOrderController{}, "post:ClientRefundOrder")
	beego.Router("/ShopOrder/UpdateIdNum", &shop.ShopOrderController{}, "post:UpdateIdNum")
	beego.Router("/ShopOrder/ExportCsv", &shop.ShopOrderController{}, "post:ExportCsv")

	beego.Router("/ShopOrder/SetPayOk", &shop.ShopOrderController{}, "post:SetPayOk")
	beego.Router("/ShopOrder/CheckPayOk", &shop.ShopOrderController{}, "post:CheckPayOk")
	beego.Router("/ShopOrder/CheckPayNO", &shop.ShopOrderController{}, "post:CheckPayNO")
	beego.Router("/ShopOrder/ExportToErp", &shop.ShopOrderController{}, "post:ExportToErp")

	beego.Router("/ShopOrder/UpdateOrderShipNum", &shop.ShopOrderController{}, "post:UpdateOrderShipNum")
	beego.Router("/ShopOrder/Adminclose", &shop.ShopOrderController{}, "post:Adminclose")
	beego.Router("/ShopOrder/AdminRefundSure", &shop.ShopOrderController{}, "post:AdminRefundSure")
	beego.Router("/ShopOrder/GetTotalPayId", &shop.ShopOrderController{}, "post:GetTotalPayId")
	beego.Router("/ShopOrder/ClientConfirmOrder", &shop.ShopOrderController{}, "post:ClientConfirmOrder")
	beego.Router("/ShopOrder/AdminConfirmOrder", &shop.ShopOrderController{}, "post:AdminConfirmOrder")
	beego.Router("/ShopOrder/AdminCancelRefund", &shop.ShopOrderController{}, "post:AdminCancelRefund")
	beego.Router("/ShopOrder/ClientCancelRefund", &shop.ShopOrderController{}, "post:ClientCancelRefund")

	//支付码
	beego.Router("/paycode/all", &shop.PayCodeController{}, "post:All")

	//标签
	beego.Router("/ShopTag/all", &shop.ShopTagController{}, "post:All")
	beego.Router("/ShopTag/edit", &shop.ShopTagController{}, "post:Edit")
	beego.Router("/ShopTag/add", &shop.ShopTagController{}, "post:Add")
	beego.Router("/ShopTag/del", &shop.ShopTagController{}, "post:Del")
}
