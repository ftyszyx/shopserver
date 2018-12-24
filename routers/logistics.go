package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/admin"
)

//物流相关
func initLogistics() {

	//物流
	beego.Router("/Logistics/all", &admin.LogisticsController{}, "post:All")
	beego.Router("/Logistics/add", &admin.LogisticsController{}, "post:Add")
	beego.Router("/Logistics/edit", &admin.LogisticsController{}, "post:Edit")
	beego.Router("/Logistics/UpdateTask", &admin.LogisticsController{}, "post:UpdateTask")
	beego.Router("/Logistics/UpdateAllTask", &admin.LogisticsController{}, "post:UpdateAllTask")
	beego.Router("/Logistics/ExportCsv", &admin.LogisticsController{}, "post:ExportCsv")
	beego.Router("/Logistics/GetLogicsInfo", &admin.LogisticsController{}, "post:GetLogicsInfo")
	beego.Router("/Logistics/AddLogicAPI", &admin.LogisticsController{}, "post:AddLogicAPI")
	beego.Router("/Logistics/UploadeLogistics", &admin.LogisticsController{}, "post:UploadeLogistics")
	beego.Router("/Logistics/SyncErpData", &admin.LogisticsController{}, "post:SyncErpData")
	beego.Router("/Logistics/ClientChangeInfo", &admin.LogisticsController{}, "post:ClientChangeInfo")

	//物流任务
	beego.Router("/LogisticsTask/all", &admin.LogisticsTaskController{}, "post:All")
	beego.Router("/LogisticsTask/edit", &admin.LogisticsTaskController{}, "post:Edit")
	beego.Router("/LogisticsTask/add", &admin.LogisticsTaskController{}, "post:Add")
	beego.Router("/LogisticsTask/del", &admin.LogisticsTaskController{}, "post:Del")
}
