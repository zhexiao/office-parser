package utils

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"log"
)

type Qiniu struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Zone      string

	Domain string
}

func (q Qiniu) uploadCloud(key string, localFile string) storage.PutRet {
	//上传
	putPolicy := storage.PutPolicy{
		Scope:   q.Bucket,
		Expires: 3600,
	}

	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	switch q.Zone {
	case "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	default:
		cfg.Zone = &storage.ZoneHuadong
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
	}

	return ret
}

var OfficeParserQiniuCfg *Qiniu

func UploadFileToQiniu(key string, localFile string) string {
	if OfficeParserQiniuCfg == nil {
		log.Fatal("没有实例化office-parser的七牛配置，请检查`OfficeParserQiniuCfg`变量")
	}

	//图片地址新增前缀office_parser
	pathKey := fmt.Sprintf("office_parser/%s", key)
	ret := OfficeParserQiniuCfg.uploadCloud(pathKey, localFile)

	//返回地址
	return fmt.Sprintf("%s/%s", OfficeParserQiniuCfg.Domain, ret.Key)
}
