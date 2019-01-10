package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/control/logistics"
)

//物流网站

func initShipHome() {
	beego.Router("/", &logistics.ShipHomeControl{}, "get:Home")
	beego.Router("/home", &logistics.ShipHomeControl{}, "get:Home")
	beego.Router("/home/index", &logistics.ShipHomeControl{}, "get:Home")
	beego.Router("/about", &logistics.ShipHomeControl{}, "get:About")
	beego.Router("/track/?:id", &logistics.ShipHomeControl{}, "get:Track")
	beego.Router("/service", &logistics.ShipHomeControl{}, "get:Service")
	beego.Router("/new/:id", &logistics.ShipHomeControl{}, "get:New")
	beego.Router("/uploadinfo", &logistics.ShipHomeControl{}, "get:Upload")
	beego.Router("/newslist/:page", &logistics.ShipHomeControl{}, "get:NewsList")
}
