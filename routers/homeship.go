package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/admin"
)

func initShipHome() {
	beego.Router("/", &admin.ShipHomeControl{}, "get:Home")
	beego.Router("/home", &admin.ShipHomeControl{}, "get:Home")
	beego.Router("/home/index", &admin.ShipHomeControl{}, "get:Home")
	beego.Router("/about", &admin.ShipHomeControl{}, "get:About")
	beego.Router("/track/?:id", &admin.ShipHomeControl{}, "get:Track")
	beego.Router("/service", &admin.ShipHomeControl{}, "get:Service")
	beego.Router("/new/:id", &admin.ShipHomeControl{}, "get:New")
	beego.Router("/uploadinfo", &admin.ShipHomeControl{}, "get:Upload")
	beego.Router("/newslist/:page", &admin.ShipHomeControl{}, "get:NewsList")
}
