package libs

import (
	"context"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

func GetManger() *storage.BucketManager {
	accesskey := beego.AppConfig.String("qiniu.accesskey")
	secretkey := beego.AppConfig.String("qiniu.secretkey")
	mac := qbox.NewMac(accesskey, secretkey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	return bucketManager
}

func GetUploadToken(key string, bucket string) string {
	accesskey := beego.AppConfig.String("qiniu.accesskey")
	secretkey := beego.AppConfig.String("qiniu.secretkey")

	Scopestr := bucket
	if key != "" {
		Scopestr = fmt.Sprintf("%s:%s", bucket, key)
	}
	putPolicy := storage.PutPolicy{
		Scope:      Scopestr,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	putPolicy.Expires = 7200 //示例2小时有效期
	mac := qbox.NewMac(accesskey, secretkey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

func UploadFile(key string, localfile string, bucket string) (storage.PutRet, error) {
	upToken := GetUploadToken(key, bucket)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localfile, &putExtra)
	return ret, errors.WithStack(err)
}

func getPrivateRrl(domain string, key string) string {
	accesskey := beego.AppConfig.String("qiniu.accesskey")
	secretkey := beego.AppConfig.String("qiniu.secretkey")
	mac := qbox.NewMac(accesskey, secretkey)
	// domain := "https://image.example.com"
	// key := "这是一个测试文件.jpg"
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)
	return privateAccessURL
}
