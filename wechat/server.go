package wechat

import (
	"github.com/astaxie/beego/logs"
	"github.com/zyx/shop_server/wechat/oauth"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/zyx/shop_server/wechat/message"
)

var Instance *Wechat
var OauthInstance *oauth.Oauth

func init() {
	logs.Info("init wechat")
	config := &Config{
		AppID:          beego.AppConfig.String("wechat.appid"),
		AppSecret:      beego.AppConfig.String("wechat.appsecret"),
		Token:          beego.AppConfig.String("wechat.token"),
		EncodingAESKey: beego.AppConfig.String("wechat.encodingAESKey"),
	}
	Instance = NewWechat(config)
	OauthInstance = Instance.GetOauth()
}

func Resolve(ctx *context.Context) {
	// 传入request和responseWriter
	logs.Info("resolve")
	server := Instance.GetServer(ctx.Request, ctx.ResponseWriter)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {

		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	//发送回复的消息
	server.Send()
}