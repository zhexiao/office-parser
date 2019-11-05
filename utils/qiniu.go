package utils

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"io"
	"log"
)

type Qiniu struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Zone      string

	Domain string
}

var OfficeParserQiniuCfg *Qiniu
var uploadPrefix = "office_parser"

//文件传输
func UploadFileToQiniu(key string, localFile string) string {
	checkCfg()

	//图片地址新增前缀office_parser
	pathKey := fmt.Sprintf("%s/%s", uploadPrefix, key)
	ret := OfficeParserQiniuCfg.uploadCloud(pathKey, localFile)

	//返回地址
	return fmt.Sprintf("%s/%s", OfficeParserQiniuCfg.Domain, ret.Key)
}

//数据传输
func UploadDataToQiniu(key string, data io.Reader, size int64) string {
	checkCfg()

	//图片地址新增前缀office_parser
	pathKey := fmt.Sprintf("%s/%s", uploadPrefix, key)
	ret := OfficeParserQiniuCfg.uploadDataToCloud(pathKey, data, size)

	//返回地址
	return fmt.Sprintf("%s/%s", OfficeParserQiniuCfg.Domain, ret.Key)
}

func checkCfg() {
	if OfficeParserQiniuCfg == nil {
		log.Panic("没有实例化office-parser的七牛配置，请检查`OfficeParserQiniuCfg`变量")
	}
}

func (q Qiniu) readCfg() storage.Config {
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

	return cfg
}

func (q Qiniu) readToken() string {
	//上传
	putPolicy := storage.PutPolicy{
		Scope:   q.Bucket,
		Expires: 3600 * 24 * 365 * 10,
	}

	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

func (q Qiniu) uploadCloud(key string, localFile string) storage.PutRet {
	upToken := q.readToken()
	cfg := q.readCfg()

	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
	}

	return ret
}

func (q Qiniu) uploadDataToCloud(key string, data io.Reader, size int64) storage.PutRet {
	upToken := q.readToken()
	cfg := q.readCfg()

	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	err := formUploader.Put(context.Background(), &ret, upToken, key, data, size, &putExtra)
	if err != nil {
		fmt.Println(err)
	}

	return ret
}
