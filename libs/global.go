package libs

//返回的码
const SuccessCode = 1
const ErrorCode = 2
const AuthFail = 21      //账号信息不存在，需要重新登录
const NoAccessRight = 22 //没有权限
const WchatLoginOk = 23  //登录成功

//过期时间
//const ExpireTimeSpace = 86400 //过期时间

//获取行数列名
const SQL_COUNT_NAME = "tbcount"

//订单状态
const OrderStatusWaitPay = 0   //待付款
const OrderStatusWaitcheck = 1 //买家确认付款
const OrderStatusWaitSend = 2  //待发货（卖家家确认付款）
const OrderStatusSend = 3      //已发货
const OrderStatusRefund = 4    //退款中
const OrderStatusOver = 5      //已完成
const OrderStatusClose = 6     //已关闭
const OrderStatusDelete = 7    //已删除

//关闭类型
const OrderCloseByClient = 0
const OrderCloseByAdmin = 1
const OrderCloseRefund = 2

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

//获取所有数据类型
const GetAll_type = 0     //所有
const GetAll_type_num = 1 //只获取总数
