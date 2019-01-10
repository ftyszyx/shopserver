package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/control/base"
	"github.com/zyx/shop_server/control/corecontrol"
	"github.com/zyx/shop_server/wechat"
)

func initCommon() {
	//登录
	beego.Router("/login/login", &corecontrol.LoginController{}, "*:Login")
	beego.Router("/login/logout", &corecontrol.LoginController{}, "*:LoginOut")
	beego.Router("/login/getcaptcha", &corecontrol.LoginController{}, "get:GetCaptchaCode")
	beego.Router("/login/getphonecode", &corecontrol.LoginController{}, "post:GetPhoneCode")
	beego.Router("/login/loginwithphone", &corecontrol.LoginController{}, "post:LoginWithPhone")

	//用户管理
	beego.Router("/user/all", &corecontrol.UserController{}, "post:All")
	beego.Router("/user/edit", &corecontrol.UserController{}, "post:Edit")
	beego.Router("/user/add", &corecontrol.UserController{}, "post:Add")
	beego.Router("/user/getUserInfo", &corecontrol.UserController{}, "post:GetUserInfo")
	beego.Router("/user/getshopuserinfo", &corecontrol.UserController{}, "get:GetShopUserInfo")
	beego.Router("/user/del", &corecontrol.UserController{}, "post:Del")
	beego.Router("/user/changePassword", &corecontrol.UserController{}, "post:ChangePassword")
	beego.Router("/user/ChangeValid", &corecontrol.UserController{}, "post:ChangeValid")
	beego.Router("/user/UpdateCart", &corecontrol.UserController{}, "post:UpdateCart")

	beego.Router("/user/UpdateName", &corecontrol.UserController{}, "post:UpdateName")
	beego.Router("/user/UpdateHead", &corecontrol.UserController{}, "post:UpdateHead")
	beego.Router("/user/UpdatePhone", &corecontrol.UserController{}, "post:UpdatePhone")
	beego.Router("/user/UpdateAddress", &corecontrol.UserController{}, "post:UpdateAddress")
	beego.Router("/user/RefreshToken", &corecontrol.UserController{}, "post:RefreshToken")
	beego.Router("/user/UpdateAccount", &corecontrol.UserController{}, "post:UpdateAccount")

	beego.Router("/user/rsetpass", &corecontrol.UserController{}, "post:ResetPassword")

	//用户组
	beego.Router("/user_group/all", &corecontrol.UserGroupController{}, "post:All")
	beego.Router("/user_group/edit", &corecontrol.UserGroupController{}, "post:Edit")
	beego.Router("/user_group/add", &corecontrol.UserGroupController{}, "post:Add")
	beego.Router("/user_group/del", &corecontrol.UserGroupController{}, "post:Del")

	//文章
	beego.Router("/post/all", &corecontrol.PostController{}, "post:All")
	beego.Router("/post/edit", &corecontrol.PostController{}, "post:Edit")
	beego.Router("/post/add", &corecontrol.PostController{}, "post:Add")
	beego.Router("/post/del", &corecontrol.PostController{}, "post:Del")
	beego.Router("/post/abandon", &corecontrol.PostController{}, "post:Abandon")

	//文章类型
	beego.Router("/posttype/all", &corecontrol.PostTypeController{}, "post:All")
	beego.Router("/posttype/edit", &corecontrol.PostTypeController{}, "post:Edit")
	beego.Router("/posttype/add", &corecontrol.PostTypeController{}, "post:Add")
	beego.Router("/posttype/del", &corecontrol.PostTypeController{}, "post:Del")
	beego.Router("/posttype/abandon", &corecontrol.PostTypeController{}, "post:Abandon")

	//图片
	beego.Router("/photo/all", &corecontrol.PhotoController{}, "post:All")
	beego.Router("/photo/edit", &corecontrol.PhotoController{}, "post:Edit")
	beego.Router("/photo/add", &corecontrol.PhotoController{}, "post:Add")
	beego.Router("/photo/del", &corecontrol.PhotoController{}, "post:Del")
	beego.Router("/photo/movemulti", &corecontrol.PhotoController{}, "post:MoveMulti")
	beego.Router("/photo/getuploadtoken", &corecontrol.PhotoController{}, "post:GetUploadtoken")

	//相册
	beego.Router("/album/all", &corecontrol.AlbumController{}, "post:All")
	beego.Router("/album/edit", &corecontrol.AlbumController{}, "post:Edit")
	beego.Router("/album/add", &corecontrol.AlbumController{}, "post:Add")
	beego.Router("/album/del", &corecontrol.AlbumController{}, "post:Del")
	beego.Router("/album/changecover", &corecontrol.AlbumController{}, "post:ChangeCover")
	beego.Router("/album/changedefault", &corecontrol.AlbumController{}, "post:ChangeDefault")

	//模块
	beego.Router("/module/all", &corecontrol.ModuleController{}, "post:All")
	//日志
	beego.Router("/log/all", &corecontrol.LogController{}, "post:All")
	//设置
	beego.Router("/config/all", &corecontrol.ConfigController{}, "post:All")
	beego.Router("/config/edit", &corecontrol.ConfigController{}, "post:Edit")

	//导表模板
	beego.Router("/Export/all", &corecontrol.ExportController{}, "post:All")
	beego.Router("/Export/edit", &corecontrol.ExportController{}, "post:Edit")
	beego.Router("/Export/add", &corecontrol.ExportController{}, "post:Add")
	beego.Router("/Export/del", &corecontrol.ExportController{}, "post:Del")

	//数据库备份
	beego.Router("/database/all", &corecontrol.DatabaseController{}, "post:All")
	beego.Router("/database/edit", &corecontrol.DatabaseController{}, "post:Edit")
	beego.Router("/database/add", &corecontrol.DatabaseController{}, "post:Add")
	beego.Router("/database/del", &corecontrol.DatabaseController{}, "post:Del")
	beego.Router("/database/restore", &corecontrol.DatabaseController{}, "post:Restore")

	//导表任务
	beego.Router("/ExportTask/all", &corecontrol.ExportTaskController{}, "post:All")

	//ueditor
	beego.Router("/upload", &base.UploadController{})
	//图片上传
	beego.Router("/picupload", &base.UploadController{}, "post:PicUpload")
	beego.Router("/uploadidnum", &base.UploadController{}, "post:UploadIDNum")

	//刷新缓存
	beego.Router("/system/refresh", &corecontrol.SystemController{}, "post:Refresh")

	//微信相关
	beego.Any("/wchat", wechat.Resolve)
	beego.Router("/login/LoginWithWchat", &corecontrol.LoginController{}, "post:LoginWithWchat")
	beego.Router("/wchatauthcallback", &corecontrol.LoginController{}, "*:WchatLoginCallback")
	beego.Router("/user/GetWchatJsConf", &corecontrol.UserController{}, "*:GetWchatJsConf")

}
