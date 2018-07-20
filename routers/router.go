package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/admin"
	"github.com/zyx/shop_server/wechat"
)

func init() {
	//登录
	beego.Router("/login/login", &admin.LoginController{}, "*:Login")
	beego.Router("/login/logout", &admin.LoginController{}, "*:LoginOut")
	beego.Router("/login/getcaptcha", &admin.LoginController{}, "get:GetCaptchaCode")
	beego.Router("/login/getphonecode", &admin.LoginController{}, "post:GetPhoneCode")
	beego.Router("/login/loginwithphone", &admin.LoginController{}, "post:LoginWithPhone")

	//商城
	beego.Router("/shop/getinfo", &admin.ShopController{}, "get:GetInfo")

	//用户管理
	beego.Router("/user/all", &admin.UserController{}, "post:All")
	beego.Router("/user/edit", &admin.UserController{}, "post:Edit")
	beego.Router("/user/add", &admin.UserController{}, "post:Add")
	beego.Router("/user/getUserInfo", &admin.UserController{}, "post:GetUserInfo")
	beego.Router("/user/getshopuserinfo", &admin.UserController{}, "get:GetShopUserInfo")
	beego.Router("/user/del", &admin.UserController{}, "post:Del")
	beego.Router("/user/changePassword", &admin.UserController{}, "post:ChangePassword")
	beego.Router("/user/ChangeValid", &admin.UserController{}, "post:ChangeValid")
	beego.Router("/user/UpdateCart", &admin.UserController{}, "post:UpdateCart")

	beego.Router("/user/UpdateName", &admin.UserController{}, "post:UpdateName")
	beego.Router("/user/UpdateHead", &admin.UserController{}, "post:UpdateHead")
	beego.Router("/user/UpdatePhone", &admin.UserController{}, "post:UpdatePhone")
	beego.Router("/user/UpdateAddress", &admin.UserController{}, "post:UpdateAddress")
	beego.Router("/user/RefreshToken", &admin.UserController{}, "post:RefreshToken")

	//用户组
	beego.Router("/user_group/all", &admin.UserGroupController{}, "post:All")
	beego.Router("/user_group/edit", &admin.UserGroupController{}, "post:Edit")
	beego.Router("/user_group/add", &admin.UserGroupController{}, "post:Add")
	beego.Router("/user_group/del", &admin.UserGroupController{}, "post:Del")

	//文章
	beego.Router("/post/all", &admin.PostController{}, "post:All")
	beego.Router("/post/edit", &admin.PostController{}, "post:Edit")
	beego.Router("/post/add", &admin.PostController{}, "post:Add")
	beego.Router("/post/del", &admin.PostController{}, "post:Del")
	beego.Router("/post/abandon", &admin.PostController{}, "post:Abandon")

	//文章类型
	beego.Router("/posttype/all", &admin.PostTypeController{}, "post:All")
	beego.Router("/posttype/edit", &admin.PostTypeController{}, "post:Edit")
	beego.Router("/posttype/add", &admin.PostTypeController{}, "post:Add")
	beego.Router("/posttype/del", &admin.PostTypeController{}, "post:Del")
	beego.Router("/posttype/abandon", &admin.PostTypeController{}, "post:Abandon")

	//图片
	beego.Router("/photo/all", &admin.PhotoController{}, "post:All")
	beego.Router("/photo/edit", &admin.PhotoController{}, "post:Edit")
	beego.Router("/photo/add", &admin.PhotoController{}, "post:Add")
	beego.Router("/photo/del", &admin.PhotoController{}, "post:Del")
	beego.Router("/photo/movemulti", &admin.PhotoController{}, "post:MoveMulti")
	beego.Router("/photo/getuploadtoken", &admin.PhotoController{}, "post:GetUploadtoken")

	//相册
	beego.Router("/album/all", &admin.AlbumController{}, "post:All")
	beego.Router("/album/edit", &admin.AlbumController{}, "post:Edit")
	beego.Router("/album/add", &admin.AlbumController{}, "post:Add")
	beego.Router("/album/del", &admin.AlbumController{}, "post:Del")
	beego.Router("/album/changecover", &admin.AlbumController{}, "post:ChangeCover")
	beego.Router("/album/changedefault", &admin.AlbumController{}, "post:ChangeDefault")

	//模块
	beego.Router("/module/all", &admin.ModuleController{}, "post:All")
	//日志
	beego.Router("/log/all", &admin.LogController{}, "post:All")
	//设置
	beego.Router("/config/all", &admin.ConfigController{}, "post:All")
	beego.Router("/config/edit", &admin.ConfigController{}, "post:Edit")

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

	//标签
	beego.Router("/ShopTag/all", &admin.ShopTagController{}, "post:All")
	beego.Router("/ShopTag/edit", &admin.ShopTagController{}, "post:Edit")
	beego.Router("/ShopTag/add", &admin.ShopTagController{}, "post:Add")
	beego.Router("/ShopTag/del", &admin.ShopTagController{}, "post:Del")

	//导表模板
	beego.Router("/Export/all", &admin.ExportController{}, "post:All")
	beego.Router("/Export/edit", &admin.ExportController{}, "post:Edit")
	beego.Router("/Export/add", &admin.ExportController{}, "post:Add")
	beego.Router("/Export/del", &admin.ExportController{}, "post:Del")

	//数据库备份
	beego.Router("/database/all", &admin.DatabaseController{}, "post:All")
	beego.Router("/database/edit", &admin.DatabaseController{}, "post:Edit")
	beego.Router("/database/add", &admin.DatabaseController{}, "post:Add")
	beego.Router("/database/del", &admin.DatabaseController{}, "post:Del")
	beego.Router("/database/restore", &admin.DatabaseController{}, "post:Restore")

	//导表任务
	beego.Router("/ExportTask/all", &admin.ExportTaskController{}, "post:All")

	//ueditor
	beego.Router("/upload", &admin.UploadController{})
	//图片上传
	beego.Router("/picupload", &admin.UploadController{}, "post:PicUpload")
	//刷新缓存
	beego.Router("/system/refresh", &admin.SystemController{}, "post:Refresh")

	//微信相关
	beego.Any("/wchat", wechat.Resolve)
	beego.Router("/login/LoginWithWchat", &admin.LoginController{}, "post:LoginWithWchat")
	beego.Router("/wchatauthcallback", &admin.LoginController{}, "*:WchatLoginCallback")
	// beego.Router("/CheckWchatLogin", &admin.LoginController{}, "get:CheckWchatLogin")

	initCompany()
}
