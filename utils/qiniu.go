package utils

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type Qiniu struct {
	accessKey string
	secretKey string
	bucket    string
	zone      string
}

func (q Qiniu) uploadCloud(key string, localFile string) storage.PutRet {
	//上传
	putPolicy := storage.PutPolicy{
		Scope:   q.bucket,
		Expires: 3600,
	}

	mac := qbox.NewMac(q.accessKey, q.secretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	switch q.zone {
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

func UploadFileToQiniu(key string, localFile string) string {
	cfg := QiniuCfg()

	q := Qiniu{
		accessKey: cfg.accessKey,
		secretKey: cfg.secretKey,
		bucket:    cfg.bucket,
		zone:      cfg.zone,
	}
	ret := q.uploadCloud(key, localFile)

	//返回地址
	return fmt.Sprintf("%s/%s", cfg.domain, ret.Key)
}
