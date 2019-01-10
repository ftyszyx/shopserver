package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/control/home"
)

func initCompany() {
	beego.Router("/", &home.HomeControl{}, "get:Home")
	beego.Router("/home", &home.HomeControl{}, "get:Home")
	beego.Router("/home/index", &home.HomeControl{}, "get:Home")
	beego.Router("/home/news", &home.HomeControl{}, "get:News")
	beego.Router("/home/products", &home.HomeControl{}, "get:Products")
	beego.Router("/home/post/:id", &home.HomeControl{}, "get:Post")
}
