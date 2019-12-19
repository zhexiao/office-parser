# office-parser
把教育信息化体系中的Word试题，Excel试卷、知识点等数据解析成json内容。

# 配置文件
1. 如果独立使用，则可以创建settings.yaml自定义配置
```
$ cp settings.yaml.example settings.yaml
```

2. 如果作为第三方库使用，则可以由用户通过以下方法配置
```
func init() {
	utils.OfficeParserQiniuCfg = &utils.Qiniu{
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

# 测试
支持的类型：
1. book
2. paper
3. outline
4. cognition_map
5. cognition_sp
6. question

运行
```
$ go run main.go -h

$ go run main.go -f ./test/book.xlsx -t book
```