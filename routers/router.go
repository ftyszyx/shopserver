package routers

import (
	"fmt"

	"github.com/astaxie/beego"
)

func init() {

}

func InitAllRoute() {
	appname := beego.AppConfig.String("appname")

	fmt.Println(fmt.Sprintf("init router:appname: %s", appname))
	initCommon()

	//商城
	if appname == "shop" {
		initLogistics()
		initShop()
	}

	//公司物流
	if appname == "ship" {
		initLogistics()
		initShipHome() //物流前端一些接口
	}

	//公司首页
	if appname == "home" {
		initCompany() //公司网站
	}
}
