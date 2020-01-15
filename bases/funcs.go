package bases

import (
	"regexp"
)

//返回需要读取的字符串列表
func ReadNum(text string) []string {
	//这里必须使用非贪婪模式
	reg := regexp.MustCompile(`\[([\w\W\d-]+?)\]`)
	data := reg.FindAllStringSubmatch(text, -1)

	var res []string
	for _, dt := range data {
		res = append(res, dt[1])
	}

	return res
}

//返回需要读取的文本文字
func ReadText(text string) string {
	//这里必须使用非贪婪模式
	reg := regexp.MustCompile(`\[([\w\W\d-]+?)\](.*)`)
	data := reg.FindAllStringSubmatch(text, -1)

	var res string
	for _, dt := range data {
		res = dt[2]
	}

	return res
}
