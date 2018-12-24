package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/admin"
)

func initShop() {

	//商城
	beego.Router("/shop/getinfo", &admin.ShopController{}, "get:GetInfo")

	//广告
	beego.Router("/ads/all", &admin.AdsController{}, "post:All")
	beego.Router("/ads/edit", &admin.AdsController{}, "post:Edit")
	beego.Router("/ads/add", &admin.AdsController{}, "post:Add")
	beego.Router("/ads/del", &admin.AdsController{}, "post:Del")

	//广告位
	beego.Router("/adspos/all", &admin.AdsPosController{}, "post:All")
	beego.Router("/adspos/edit", &admin.AdsPosController{}, "post:Edit")
	beego.Router("/adspos/add", &admin.AdsPosController{}, "post:Add")
	beego.Router("/adspos/del", &admin.AdsPosController{}, "post:Del")

	//品牌
	beego.Router("/ShopBrand/all", &admin.ShopBrandController{}, "post:All")
	beego.Router("/ShopBrand/edit", &admin.ShopBrandController{}, "post:Edit")
	beego.Router("/ShopBrand/add", &admin.ShopBrandController{}, "post:Add")
	beego.Router("/ShopBrand/del", &admin.ShopBrandController{}, "post:Del")

	//商品
	beego.Router("/ShopItem/all", &admin.ShopItemController{}, "post:All")
	beego.Router("/ShopItem/edit", &admin.ShopItemController{}, "post:Edit")
	beego.Router("/ShopItem/add", &admin.ShopItemController{}, "post:Add")
	beego.Router("/ShopItem/del", &admin.ShopItemController{}, "post:Del")
	beego.Router("/ShopItem/ExportCsv", &admin.ShopItemController{}, "post:ExportCsv")

	//商品类型
	beego.Router("/ShopItemType/all", &admin.ShopItemTypeController{}, "post:All")
	beego.Router("/ShopItemType/edit", &admin.ShopItemTypeController{}, "post:Edit")
	beego.Router("/ShopItemType/add", &admin.ShopItemTypeController{}, "post:Add")
	beego.Router("/ShopItemType/del", &admin.ShopItemTypeController{}, "post:Del")

	//公告
	beego.Router("/ShopNotice/all", &admin.ShopNoticeController{}, "post:All")
	beego.Router("/ShopNotice/edit", &admin.ShopNoticeController{}, "post:Edit")
	beego.Router("/ShopNotice/add", &admin.ShopNoticeController{}, "post:Add")
	beego.Router("/ShopNotice/del", &admin.ShopNoticeController{}, "post:Del")

	//订单
	beego.Router("/ShopOrder/edit", &admin.ShopOrderController{}, "post:Edit")
	beego.Router("/ShopOrder/all", &admin.ShopOrderController{}, "post:All")
	beego.Router("/ShopOrder/add", &admin.ShopOrderController{}, "post:Add")
	beego.Router("/ShopOrder/orderuploade", &admin.ShopOrderController{}, "post:OrdersUpload")

	beego.Router("/ShopOrder/GetMyOrder", &admin.ShopOrderController{}, "post:GetMyOrder")
	beego.Router("/ShopOrder/ExportMyCsv", &admin.ShopOrderController{}, "post:ExportMyCsv")
	beego.Router("/ShopOrder/ClientClose", &admin.ShopOrderController{}, "post:ClientClose")
	beego.Router("/ShopOrder/ClientDelOrder", &admin.ShopOrderController{}, "post:ClientDelOrder")
	beego.Router("/ShopOrder/ClientRefundOrder", &admin.ShopOrderController{}, "post:ClientRefundOrder")
	beego.Router("/ShopOrder/UpdateIdNum", &admin.ShopOrderController{}, "post:UpdateIdNum")
	beego.Router("/ShopOrder/ExportCsv", &admin.ShopOrderController{}, "post:ExportCsv")

	beego.Router("/ShopOrder/SetPayOk", &admin.ShopOrderController{}, "post:SetPayOk")
	beego.Router("/ShopOrder/CheckPayOk", &admin.ShopOrderController{}, "post:CheckPayOk")
	beego.Router("/ShopOrder/CheckPayNO", &admin.ShopOrderController{}, "post:CheckPayNO")
	beego.Router("/ShopOrder/ExportToErp", &admin.ShopOrderController{}, "post:ExportToErp")

	beego.Router("/ShopOrder/UpdateOrderShipNum", &admin.ShopOrderController{}, "post:UpdateOrderShipNum")
	beego.Router("/ShopOrder/Adminclose", &admin.ShopOrderController{}, "post:Adminclose")
	beego.Router("/ShopOrder/AdminRefundSure", &admin.ShopOrderController{}, "post:AdminRefundSure")
	beego.Router("/ShopOrder/GetTotalPayId", &admin.ShopOrderController{}, "post:GetTotalPayId")
	beego.Router("/ShopOrder/ClientConfirmOrder", &admin.ShopOrderController{}, "post:ClientConfirmOrder")
	beego.Router("/ShopOrder/AdminConfirmOrder", &admin.ShopOrderController{}, "post:AdminConfirmOrder")
	beego.Router("/ShopOrder/AdminCancelRefund", &admin.ShopOrderController{}, "post:AdminCancelRefund")
	beego.Router("/ShopOrder/ClientCancelRefund", &admin.ShopOrderController{}, "post:ClientCancelRefund")

	//支付码
	beego.Router("/paycode/all", &admin.PayCodeController{}, "post:All")

	//标签
	beego.Router("/ShopTag/all", &admin.ShopTagController{}, "post:All")
	beego.Router("/ShopTag/edit", &admin.ShopTagController{}, "post:Edit")
	beego.Router("/ShopTag/add", &admin.ShopTagController{}, "post:Add")
	beego.Router("/ShopTag/del", &admin.ShopTagController{}, "post:Del")
}
