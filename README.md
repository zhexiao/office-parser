# office-parser
Word、Excel数据解析。

# 配置文件
修改utils/settings.go对系统进行配置，如果不存在此文件，则复制settings.go.example并对应修改此文件。
```
$ cp utils/settings.go.example utils/settings.go 
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