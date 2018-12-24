#shop_server

#商城服务端

编译后生成shop_server.exe

启动：./shop_server

参数：

log  指定日志文件名

backupsql 支行数据库备份命令(只是执行一下命令，服务器将不启动)

public  设置静态资源目录

port  服务器启动的商品号



##go语言的一些笔记
1、字符串打印格式

// 使用动词 v 格式化 arg 列表，非字符串元素之间添加空格
Print(arg列表)

// 使用动词 v 格式化 arg 列表，所有元素之间添加空格，结尾添加换行符

Println(arg列表)
// 使用格式字符串格式化 arg 列表
Printf(格式字符串, arg列表)
2、
https://www.cnblogs.com/golove/p/3284304.html


3、string byte[]互转
s := “abc”
b := []byte(s)
s2 := string(b)

4、string int 互转
int,err:=strconv.Atoi(string)
#string到int64  
int64, err := strconv.ParseInt(string, 10, 64)  
#int到string  
string:=strconv.Itoa(int)  


5、interface其实就是类型的指针，

6、govendor使用
安装 go get -u -v github.com/kardianos/govendor
#初始化vendor目录
govendor init
#拉取本地vendor.json中的依赖包
govendor sync
#增加
govendor add github.com/astaxie/beego/cache


7、
用nohup 运行程序
nohup ./beepkg &
nohup command > myout.file 2>&1 &

通过ps -aux|grep 来查看程序运行状态
也可以通过 jobs -l来查看（有时候看不到）

强制类型转换：
int64(a) 


netapp -authtoken    d2ca7c6ba422f9dd 
web:http://k9xmvz.natappfree.cc


wget获取网站：
wget -r --no-parent -e robots=off http://www.example.com 
-r 递归
--no-parent 不访问父节点
-e robots=off  让wget耍流氓无视robots.txt协议
-U “Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9.1.6) Gecko/20091201 Firefox/3.5.6” 伪造agent信息

定时任务：
crontab -e
第1列分钟0～59
第2列小时0～23（0表示子夜）
第3列日1～31
第4列月1～12
第5列星期0～7（0和7表示星期天）
第6列要运行的命令


