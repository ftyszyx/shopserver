package routers

import (
	"github.com/astaxie/beego"
	"github.com/zyx/shop_server/control/logistics"
)

//物流相关
func initLogistics() {

	//物流
	beego.Router("/Logistics/all", &logistics.LogisticsController{}, "post:All")
	beego.Router("/Logistics/add", &logistics.LogisticsController{}, "post:Add")
	beego.Router("/Logistics/edit", &logistics.LogisticsController{}, "post:Edit")
	beego.Router("/Logistics/UpdateTask", &logistics.LogisticsController{}, "post:UpdateTask")
	beego.Router("/Logistics/UpdateAllTask", &logistics.LogisticsController{}, "post:UpdateAllTask")
	beego.Router("/Logistics/ExportCsv", &logistics.LogisticsController{}, "post:ExportCsv")
	beego.Router("/Logistics/GetLogicsInfo", &logistics.LogisticsController{}, "post:GetLogicsInfo")
	beego.Router("/Logistics/AddLogicAPI", &logistics.LogisticsController{}, "post:AddLogicAPI")
	beego.Router("/Logistics/UploadeLogistics", &logistics.LogisticsController{}, "post:UploadeLogistics")
	beego.Router("/Logistics/SyncErpData", &logistics.LogisticsController{}, "post:SyncErpData")
	beego.Router("/Logistics/ClientChangeInfo", &logistics.LogisticsController{}, "post:ClientChangeInfo")

	//物流任务
	beego.Router("/LogisticsTask/all", &logistics.LogisticsTaskController{}, "post:All")
	beego.Router("/LogisticsTask/edit", &logistics.LogisticsTaskController{}, "post:Edit")
	beego.Router("/LogisticsTask/add", &logistics.LogisticsTaskController{}, "post:Add")
	beego.Router("/LogisticsTask/del", &logistics.LogisticsTaskController{}, "post:Del")
}
