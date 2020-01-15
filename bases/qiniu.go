package bases

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"io"
)

type CT_Qiniu struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Zone      string

	Domain       string
	UploadPrefix string
}

var OpQiniu *CT_Qiniu

func checkCfg() error {
	if OpQiniu == nil {
		return NewOpError(NormalError, "没有实例化office-parser的七牛配置，请检查`OpQiniu`变量")
	}

	return nil
}

//文件传输
func UploadFileToQiniu(key string, localFile string) (string, error) {
	if err := checkCfg(); err != nil {
		return "", err
	}

	//图片地址新增前缀office_parser
	pathKey := fmt.Sprintf("%s/%s", OpQiniu.UploadPrefix, key)
	ret, err := OpQiniu.uploadCloud(pathKey, localFile)
	if err != nil {
		return "", err
	}

	//返回地址
	return fmt.Sprintf("%s/%s", OpQiniu.Domain, ret.Key), nil
}

//数据传输
func UploadDataToQiniu(key string, data io.Reader, size int64) (string, error) {
	if err := checkCfg(); err != nil {
		return "", err
	}

	//图片地址新增前缀office_parser
	pathKey := fmt.Sprintf("%s/%s", OpQiniu.UploadPrefix, key)
	ret, err := OpQiniu.uploadDataToCloud(pathKey, data, size)
	if err != nil {
		return "", err
	}

	//返回地址
	return fmt.Sprintf("%s/%s", OpQiniu.Domain, ret.Key), nil
}

//读取配置文件
func (q *CT_Qiniu) readCfg() storage.Config {
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

//读取token
func (q *CT_Qiniu) readToken() string {
	//上传
	putPolicy := storage.PutPolicy{
		Scope:   q.Bucket,
		Expires: 3600 * 24 * 365 * 10,
	}

	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

//上传文件到七牛
func (q *CT_Qiniu) uploadCloud(key string, localFile string) (storage.PutRet, error) {
	upToken := q.readToken()
	cfg := q.readCfg()

	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		return storage.PutRet{}, err
	}

	return ret, nil
}

//上传数据到七牛
func (q *CT_Qiniu) uploadDataToCloud(key string, data io.Reader, size int64) (storage.PutRet, error) {
	upToken := q.readToken()
	cfg := q.readCfg()

	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	err := formUploader.Put(context.Background(), &ret, upToken, key, data, size, &putExtra)
	if err != nil {
		return storage.PutRet{}, err
	}

	return ret, nil
}
