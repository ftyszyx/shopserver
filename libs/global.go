package libs

//返回的码
const SuccessCode = 1
const ErrorCode = 2
const AuthFail = 21      //账号信息不存在，需要重新登录
const NoAccessRight = 22 //没有权限
const WchatLoginOk = 23  //登录成功

//过期时间
//const ExpireTimeSpace = 86400 //过期时间

//订单状态
const OrderStatusWaitPay = 0   //待付款
const OrderStatusWaitcheck = 1 //买家确认付款
const OrderStatusWaitSend = 2  //待发货（卖家家确认付款）
const OrderStatusSend = 3      //已发货
const OrderStatusRefund = 4    //退款中
const OrderStatusOver = 5      //已完成
const OrderStatusClose = 6     //已关闭
const OrderStatusDelete = 7    //已删除
var OrderStatusArr = []string{"待付款", "付款待确认", "待发货", "已发货", "退款中", "已完成", "已关闭", "已删除"}

//关闭类型
const OrderCloseByClient = 0
const OrderCloseByAdmin = 1
const OrderCloseRefund = 2

var OrderCloseTypeArr = []string{"客户关闭", "管理员关闭", "退款关闭"}

//用户类型
const UserSystem = 3
const UserAdmin = 2
const UserMember = 1

//订单类型
const Order_type_min = 0    //
const Order_type_normal = 0 //普通
const Order_type_photo = 1  //图片
const Order_type_video = 2  //视频
const Order_type_date = 3   //日期
const Order_type_max = 3    //
var OrderVipTypeArr = []string{"普通", "图片", "视频", "日期"}
var OrderVipTypemoney = []int{0, 5, 10, 0}

//获取所有数据类型
const GetAll_type = 0     //所有
const GetAll_type_num = 1 //只获取总数

//物流状态

const ShipOnWay = 0    //在途
const ShiSent = 1      //揽件
const ShipProblem = 2  //出问题
const ShipSign = 3     //签收
const ShipRefund = 4   //退签发件人已
const ShipCitySend = 5 //同城派件
const ShipBack = 6     //货物正处于退回发件人的途中
const ShipNotExit = 7  //不存在

//发货方式
const Supply_source_zhiyou = "1"  //直邮
const Supply_source_baoshui = "2" //保锐

var ShipStatusArr = []string{"在途", "揽件", "出问题", "签收", "退签发件人已", "同城派件", "退回"}

//总体状态
const ShipNotBeginValue = 0
const ShiOverseaValue = 1
const ShipOverseaOverValue = 2

var ShipOverseaStatusArr = []string{"未开始", "进行中", "完成"}

// const ShipindoorValue = 3

type LogistcisCodeType struct {
	Title string
	Code  string
}

//物流公司
var LogisticsCodeArr = []LogistcisCodeType{
	{"百世汇通", "huitongkuaidi"},
	{"EMS", "ems"},
	{"顺丰", "shunfeng"},
	{"天天", "tiantian"},
	{"圆通速递", "yuantong"},
	{"韵达快运", "yunda"},
	{"中通速递", "zhongtong"}}
