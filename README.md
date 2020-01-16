# office-parser
把教育信息化体系中的Word试题，Excel试卷、知识点等数据解析成json内容。

[![Build Status](https://travis-ci.org/zhexiao/office-parser.svg?branch=master)](https://travis-ci.org/zhexiao/office-parser)
[![codecov](https://codecov.io/gh/zhexiao/office-parser/branch/master/graph/badge.svg)](https://codecov.io/gh/zhexiao/office-parser)
![go](https://img.shields.io/badge/go->%3D1.13-blue)

# 运行
1. Go >= 1.13
2. 开启go mod
```
# 根据实际情况填写路径
$ go env -w GOPATH=/goproj
$ go env -w GOBIN=/goproj/bin
$ go env -w GO111MODULE=on
$ go env -w GOPROXY=https://goproxy.cn,direct 
```

#### 拉库
```
$ go mod download
```

#### 运行测试
支持的类型：word、paper、question、book、outline、cognition_map、cognition_sp
```
$ go run main.go -h

$ go run main.go -t word -f ./test/ywgs.docx 
$ go run main.go -t book -f ./test/book.xlsx 
```

# 配置文件
1. 如果独立使用，则可以创建settings.yaml自定义配置
```
$ cp settings.yaml.example settings.yaml
```

2. 如果作为第三方库使用，则可以由用户通过以下方法配置
```bash
func init() {
	utils.OpQiniu = &bases.CT_Qiniu{
		//七牛key
		AccessKey: "key",
		//七牛secret
		SecretKey: "secret",
		//七牛存储的bucket
		Bucket: "ups",
		//所属区域
		Zone: "ZoneHuanan",
		//访问的域名地址
		Domain: "https://test.com",
		//路径
		UploadPrefix : "office_parser",
	}

	//初始化WMF的配置
	utils.WmfConfiguration = &utils.CT_WmfCfg{
		Uri: "http://127.0.0.1:10002/convert",
	}
}
```


