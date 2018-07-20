package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/admin"
)

func initCompany() {
	beego.Router("/", &admin.HomeControl{}, "get:Home")
	beego.Router("/home", &admin.HomeControl{}, "get:Home")
	beego.Router("/home/index", &admin.HomeControl{}, "get:Home")
	beego.Router("/home/news", &admin.HomeControl{}, "get:News")
	beego.Router("/home/products", &admin.HomeControl{}, "get:Products")
	beego.Router("/home/post/:id", &admin.HomeControl{}, "get:Post")
}
