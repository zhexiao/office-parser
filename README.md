# office-parser
Word、Excel数据解析。

# 配置文件
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